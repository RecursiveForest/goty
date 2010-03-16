package main

import (
	"bufio"
	"flag"
	"fmt"
	"goty"
	"os"
)

var server *string = flag.String("server", "irc.freenode.org:6667", "Server to connect to in format 'irc.freenode.org:6667'")
var nick *string = flag.String("nick", "goty-bot", "IRC nick to use")

func main () {
	flag.Parse()

	if con, err := goty.Dial(*server, *nick); err != nil {
		fmt.Fprintf(os.Stderr, "goty: %s\n", err.String())
	} else {
		in := bufio.NewReader(os.Stdin)

		go func() {
			for {
				str := <-con.Read
				if closed(con.Read) {
					break
				}
				fmt.Printf("<- %s\n", str)
			}
		}()

		for {
			if input, err := in.ReadString('\n'); err != nil {
				fmt.Fprintf(os.Stderr, "goty: %s\n", err.String())
				break
			} else {
				fmt.Printf("-> %s", input)
				con.Write <- input[0:len(input)-1]
			}
		}
		if err:= con.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "goty: %s\n", err.String())
		}
	}
}
