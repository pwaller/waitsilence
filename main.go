// waitsilence: wait until there is no input on stdin for a given amount of time.
//
// Copyright 2013 Peter Waller <peter.waller@gmail.com>.
// BSD 3-clause license
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the author nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/textproto"
	"os"
	"os/exec"
	"time"
)

var timeout = flag.Duration("timeout", 1*time.Second,
	"Amount of time for which silence required before quitting")

var cmdstr = flag.String("command", "",
	"Command to execute and wait for silence")

var verbose = flag.Bool("verbose", false, "Show stderr of process, # lines "+
	"printed")

func main() {
	flag.Parse()

	keepalive, done := make(chan bool), make(chan bool)

	var cmd *exec.Cmd
	var input io.Reader = os.Stdin

	if *cmdstr != "" {
		cmd = exec.Command("sh", "-c", *cmdstr)
		if *verbose {
			cmd.Stderr = os.Stderr
		}
		var err error
		input, err = cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		go func() {
			err := cmd.Run()
			log.Print("Command exited: ", err)
			done <- true
		}()
	}

	in := textproto.NewReader(bufio.NewReader(input))

	go func() {
		var n int64
		last := time.Now()
		for {
			_, err := in.ReadLine()
			keepalive <- true
			if err != nil {
				break
			}
			if *verbose {
				fmt.Printf("%d lines (%.2fs)\r", n, time.Since(last).Seconds())
				n++
				last = time.Now()
			}
		}
		done <- true
	}()

	start := time.Now()
mainloop:
	for {
		select {
		case <-keepalive:
		case <-done:
			break mainloop
		case <-time.After(*timeout):
			break mainloop
		}
	}
	if cmd != nil {
		cmd.Process.Kill()
	}
	log.Printf("%s silence achieved after %s", *timeout, time.Since(start))
}
