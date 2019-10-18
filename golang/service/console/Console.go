package console

import (
	. "github.com/saichler/utils/golang"
	"net"
	"strconv"
	"strings"
)

type Console struct {
	socket  net.Listener
	handler ConsoleCommandHandler
}

type ConsoleCommandHandler interface {
	Name() string
	HandleCommand([]string)
	CommandList() []string
}

func NewConsole(port int, handler ConsoleCommandHandler) (*Console, error) {
	socket, e := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if e != nil {
		Error("Failed to bind to console port:" + strconv.Itoa(port))
		return nil, e
	}
	console := &Console{}
	console.socket = socket
	console.handler = handler
	go console.waitForConnection()
	return console, nil
}

func (c *Console) waitForConnection() {
	for {
		conn, e := c.socket.Accept()
		if e != nil {
			Error("Failed to accept connection:", e)
			break
		}
		prompt := c.handler.Name() + ">"
		for {
			write(prompt, conn)
			line := make([]byte, 4096)
			n, e := conn.Read(line)
			if e != nil {
				Error("Failed to read line:", e)
				break
			}
			command := strings.ToLower(strings.TrimSpace(string(line[0:n])))
			if command == "exit" || command == "quit" {
				writeln("Goodby!", conn)
				break
			} else if command == "?" {
				c.printHelp(conn)
			} else if command != "" {
				args := strings.Split(command, " ")
				c.handler.HandleCommand(args)
			}
		}
		conn.Close()
	}
	c.socket.Close()
}

func write(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
}

func writeln(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func (c *Console) printHelp(conn net.Conn) {
	writeln("? - Print this help message.", conn)
}
