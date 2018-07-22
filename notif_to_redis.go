package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var wg sync.WaitGroup

/*
// DEV ENV
const notif1_dir string = "/home/stephane/tmp/notif1/"
const trash = "/home/stephane/tmp/trash/"
*/

const notif1_dir string = "/data/qrouter/notif1/"
const trash = "/data/qrouter/notif1_trash/"
const longsms string = notif1_dir + "longsms/"

func load_folder() []string {

	var dirs []os.FileInfo
	var dir os.FileInfo
	var err error
	var ndir []string

	dirs, err = ioutil.ReadDir(notif1_dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s notif_to_redis [CRITICAL] Can not open %s: %v\n", time.Now(), notif1_dir, err)
		os.Exit(1)
	}

	for _, dir = range dirs {
		if dir.IsDir() {
			ndir = append(ndir, notif1_dir+dir.Name())
			ndir = append(ndir, longsms+dir.Name())
		}
	}
	return ndir
}

func process_file(client *redis.Client, file string, dir string) {

	var customer, bnumber, cid, p, partid, num, pkey, key string
	var bcmd *redis.BoolCmd

	tmp := strings.Split(file, "---")
	customer = tmp[0]
	bnumber = tmp[1]
	cid = tmp[2]
	p = tmp[3]

	p = strings.Split(p, ".")[0]

	if len(p) >= 10 {

		partid = p[0:8]
		num = p[8:10]

		if len(num) == 0 {
			num = "01"
		}

		pkey = bnumber + ":::" + partid

		bcmd = client.HSet(pkey, cid, num)
		if bcmd.Err() != nil {
			fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not hset %s: %v\n ", time.Now(), pkey, bcmd.Err())
		}

		bcmd = client.ExpireAt(pkey, time.Now().Add(96*time.Hour))
		if bcmd.Err() != nil {
			fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not ExpireAt %v\n", time.Now(), bcmd.Err())
		}

	}

	key = bnumber + ":::" + cid

	bcmd = client.HSet(key, "CUSTOMER", customer)
	if bcmd.Err() != nil {
		fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not hset %s: %v\n", time.Now(), key, bcmd.Err())
	}

	if len(partid) > 0 {
		bcmd = client.HSet(key, "PART", pkey)
		if bcmd.Err() != nil {
			fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not hset %v\n", time.Now(), bcmd.Err())
		}
	}

	bcmd = client.ExpireAt(key, time.Now().Add(96*time.Hour))
	if bcmd.Err() != nil {
		fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not ExpireAt %v\n", time.Now(), bcmd.Err())
	}

	err := os.Rename(dir+"/"+file, trash+"/"+file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not move file: %v\n", time.Now(), err)
	}

	fmt.Fprintf(os.Stderr, "%s notif_to_redis [INFO] Process file %s\n", time.Now(), file)
	fmt.Fprintf(os.Stderr, "%s notif_to_redis [INFO] %s / %s / %s / %s / %s / %s / %s / %s\n", time.Now(), customer, bnumber, cid, p, partid, num, pkey, key)
}

func process_dir(dir string, client *redis.Client) {
	defer wg.Done()
	var files []os.FileInfo
	var file os.FileInfo
	var err error

	files, err = ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s notif_to_redis [WARNING] Can not open %s as a dir: %v\n", time.Now(), dir, err)
	} else {
		for _, file = range files {
			if file.IsDir() {
				continue
			}
			process_file(client, file.Name(), dir)
		}
	}
}

func main() {

	runtime.GOMAXPROCS(2)

	var client *redis.Client

	client = redis.NewClient(&redis.Options{
		Addr:       "10.10.10.201:6379",
		Password:   "", // no password set
		DB:         0,  // use default DB
		MaxRetries: 100000,
	})

	var ndir []string

	for {
		ndir = load_folder()
		for _, a := range ndir {
			wg.Add(1)
			go process_dir(a, client)
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}

}
