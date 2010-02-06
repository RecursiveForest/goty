package main

import (
	"./goty"
	"fmt"
	"os"
	"bufio"
)

func main () {
	if con, err := goty.Dial("irc.freenode.org:6667", "goty-bot"); err != nil {
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
