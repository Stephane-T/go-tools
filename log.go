package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

func read(c chan string) {

	var line string
	var in *bufio.Scanner
	in = bufio.NewScanner(os.Stdin)

	for in.Scan() {
		line = in.Text()
		c <- string(line)
	}
	c <- "FLUSH"
	time.Sleep(2 * time.Second)

	os.Exit(0)
}

func mark(c chan string) {
	for {
		time.Sleep(30 * time.Second)
		c <- "FLUSH"
	}
}

func save_buffer(file string, buffer bytes.Buffer) {

	var fdlog *os.File
	var err error
	fdlog, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not open %s: %v\n", os.Args[1], err)
		fmt.Fprintf(os.Stdout, "%s", buffer.String())
	} else {
		fmt.Fprintf(fdlog, "%s", buffer.String())
		fdlog.Close()
	}
}

func main() {

	var buffer bytes.Buffer
	var line string

	c := make(chan string)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s file.log\n", os.Args[0])
		os.Exit(1)
	}

	go read(c)

	go mark(c)

	for {
		line = <-c
		buffer.WriteString(line)
		buffer.WriteString("\n")
		if line == "FLUSH" || buffer.Len() > 1024 {
			save_buffer(os.Args[1],buffer)
			buffer.Reset()
		}
	}

}



