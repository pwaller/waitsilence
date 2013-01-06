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
	"log"
	"net/textproto"
	"os"
	"time"
)

var timeout = flag.Duration("timeout", 1*time.Second,
	"Amount of time for which silence required before quitting")

func main() {
	flag.Parse()

	in := textproto.NewReader(bufio.NewReader(os.Stdin))

	keepalive := make(chan bool)

	go func() {
		for {
			line, err := in.ReadLine()
			keepalive <- true
			if err != nil {
				break
			}
		}
	}()

mainloop:
	for {
		select {
		case <-keepalive:
		case <-time.After(*timeout):
			log.Print("Timeout has expired!")
			break mainloop
		}
	}
	log.Print("Done!")
}
