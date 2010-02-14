package goty

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
)

type IRCConn struct {
	Sock *net.TCPConn
	Read, Write chan string
}

func Dial(server, nick string) (*IRCConn, os.Error) {
	read := make(chan string, 1000)
	write := make(chan string, 1000)
	con := &IRCConn{nil, read, write}
	err := con.Connect(server, nick)
	return con, err
}

func (con *IRCConn) Connect(server, nick string) os.Error {
	if raddr, err := net.ResolveTCPAddr(server); err != nil {
		return err
	} else {
		if c, err := net.DialTCP("tcp", nil, raddr); err != nil {
			return err
		} else {
			con.Sock = c
			r := bufio.NewReader(con.Sock)
			w := bufio.NewWriter(con.Sock)

			go func() {
				for {
					if closed(con.Read) {
						fmt.Fprintf(os.Stderr, "goty: read closed\n")
						break
					}
					if str, err := r.ReadString(byte('\n')); err != nil {
						fmt.Fprintf(os.Stderr, "goty: read: %s\n", err.String())
						break
					} else {
						if strings.HasPrefix(str, "PING") {
							con.Write <- "PONG" + str[4:len(str)-2]
						} else {
							con.Read <- str[0:len(str)-2]
						}
					}
				}
			}()

			go func() {
				for {
					str := <-con.Write
					if closed(con.Write) {
						fmt.Fprintf(os.Stderr, "goty: write closed\n")
						break
					}
					if _, err := w.WriteString(str + "\r\n"); err != nil {
						fmt.Fprintf(os.Stderr, "goty: write: %s\n", err.String())
						break
					}
					w.Flush()
				}
			}()

			con.Write <- "NICK " + nick
			con.Write <- "USER bot * * :..."
		}
	}
	return nil
}

func (con *IRCConn) Close() os.Error {
	close(con.Read)
	close(con.Write)
	return con.Sock.Close()
}
