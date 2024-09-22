package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error getting username: ", err)
	}
	username = strings.TrimSpace(username)

	fmt.Println("Enter your password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error getting password: ", err)
	}
	password = strings.TrimSpace(password)

	_, err = fmt.Fprintf(conn, "%s\n", username)
	if err != nil {
		fmt.Println("Error sending name to server: ", err)
	}
	_, err = fmt.Fprintf(conn, "%s\n", password)
	if err != nil {
		fmt.Println("Error sending password to server: ", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println(text)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from server: ", err)
		}
	}()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message: ", err)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error sending message to server: ", err)
		}
	}

}
