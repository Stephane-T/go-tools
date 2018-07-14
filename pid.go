

// Copyright (c) <year>, <copyright holder>
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 
// The views and conclusions contained in the software and documentation are those
// of the authors and should not be interpreted as representing official policies,
// either expressed or implied, of the <project name> project.



package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	var p *os.Process
	var pid int
	var f *os.File
	var err error
	var data []byte
	var system *exec.Cmd

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage : %s pidfile command [command args]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Source of this file can be found here: https://github.com/Stephane-T/go-tools\n")
		os.Exit(3)
	}

	data, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		goto OK
	}

	pid, err = strconv.Atoi(strings.TrimRight(string(data), "\n"))
	if err != nil {
		goto OK
	}

	fmt.Fprintf(os.Stderr, "Found current PID: %d\n", pid)

	p, err = os.FindProcess(pid)
	if err != nil {
		goto OK
	}

	err = p.Signal(syscall.Signal(0))
	if err == nil {
		os.Exit(3)
	}

OK:
	f, err = os.OpenFile(os.Args[1], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not open %s : %v\n", os.Args[1], err)
		os.Exit(3)
	}
	fmt.Fprintf(f, "%d\n", os.Getpid())
	f.Close()

	system = exec.Command(os.Args[2], os.Args[3:]...)
	data, err = system.Output()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))

}
