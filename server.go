package main

import (
	"bufio"
	"fmt"
  "sync"
	"net"
	"strings"
)

type data struct{
  conn net.Conn
  message string
}

type connection struct{
  conn net.Conn
  username string
}
var (
	clients      = make(map[net.Conn]*connection)
	clientsMutex sync.Mutex
)

func handleConnection(conn net.Conn){
  reader := bufio.NewReader(conn)
  username, err := reader.ReadString('\n')
  if err != nil{
    fmt.Println("Error reading username: ")
    return
  }
  username = strings.TrimSpace(username)

  client := &connection{conn: conn, username: username}
  clientsMutex.Lock()
  clients[conn] = client
  clientsMutex.Unlock()

  handleBroadcast(conn, fmt.Sprintf("User %s Joined!", client.username))


  scanner := bufio.NewScanner(conn)
  for scanner.Scan(){
    update := scanner.Text()
    if update == ""{
      continue
    }
    handleBroadcast(conn, fmt.Sprintf("%s: %s", client.username, update))
  }
  clientsMutex.Lock()
  delete(clients, conn)
  clientsMutex.Unlock()

  handleBroadcast(nil, fmt.Sprintf("%s has left the chat.\n", client.username))
} 


func handleBroadcast(sender net.Conn, message string){
  clientsMutex.Lock()
  defer clientsMutex.Unlock()

  for conn, client := range clients{
    if conn != sender{
      _, err := fmt.Fprintf(client.conn, message+"\n")
      if err != nil{
        fmt.Println("Error sending message back to users: ", err)
      }
    }
  }
}



func main(){
  listener, err := net.Listen("tcp", "localhost:8080")
  if err != nil{
    fmt.Println("Error listening to server: ", err)
  }

  defer listener.Close()
  fmt.Println("Server started!")
  for{
    conn, err := listener.Accept()
    if err != nil{
      fmt.Println("Error accepting connection: ", err)
      continue
    }
    go handleConnection(conn)
  }
}
