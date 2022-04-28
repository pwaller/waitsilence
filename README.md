# waitsilence

Delay until `stdin` has recieved no input for a specified time, and then exit.

## Options

	-timeout : length of time to wait, e.g. 10s
	-command : instead of reading from `stdin`, execute the specified command
	           inside a shell. The command is killed if `stdin` is silent for
	           the specified duration.

## Use case

Let's say you have process which you can signal to save everything in memory to
disk in lots of little files. You don't know in advance how long this will take,
but you want to wait as long as it's still writing. On Linux you can use
[`inotifywait`](//linux.die.net/man/1/inotifywait) which will spit out text as
long as files are still being modified on disk. You want to wait until the
output falls silent for a short while, and then you can be sure that writing has
finished. Then you can do this:

    inotifywait -m -r . | waitsilence -timeout 5s

Unfortunately, this will wait longer than necessary because `inotifywait` will
not be writing to its `stdout` when it has finished. This means that it won't
receive the `SIGPIPE` until it next writes, which could be some time. In that
case, you can run the command as a child inside `waitsilence`, where it will be
killed.

## Backing up minecraft worlds without stopping the server

If you say `save-all` to a minecraft server, it won't actually finish writing
for some time.

This is what I use to do a backup on my minecraft server:

	if ! screen -r -S minecraft -X stuff $'\nsave-all\nsave-off\n';
	then
	    echo "Screen is attached. Not taking world backup!"
	    exit 1
	fi

	CMD='inotifywait -m -r server/world --exclude "^./dynmap.*" | grep -v ACCESS'
	
	# Ensure minecraft has finished writing its files out:
	waitsilence -timeout 5s -command "$CMD"

	# < do whatever you want to backup your files >

	screen -S minecraft -X stuff $'save-on\nsay World backup completed.\n'

I can recommend [bup](http://github.com/bup/bup) for taking minecraft backups.
Hourly backups for 2 months currently occupy ~20GB, whereas vanilla copies would
now occupy more than 1.5TB.

## Installation

Use the [`go get`](//youtu.be/XCsL89YtqCs) tool, which will install it to your
`$GOPATH`.

	go get github.com/pwaller/waitsilence
