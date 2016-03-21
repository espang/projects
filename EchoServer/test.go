package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "hello")

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Fprintf(conn, "\n")
	}()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("Message: ", message)

}
