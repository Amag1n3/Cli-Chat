package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net"
	"strings"
	"sync"
  _ "github.com/mattn/go-sqlite3"
)

type connection struct {
	conn     net.Conn
	username string
}
var (
	clients      = make(map[net.Conn]*connection)
	clientsMutex sync.Mutex
	db           *sql.DB
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "/home/amogh/Desktop/cli_chat/users.db")
	if err != nil {
		fmt.Println("Error opening DB:", err)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println("Error creating table:", err)
	}
}


func userExists(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		fmt.Println("Error checking user: ", err)
	}
	return count > 0
}

func createNewUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %v", err)
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?,?)", username, hash)
	if err != nil {
		return fmt.Errorf("Error adding user to DB: %v", err)
	}
	return nil
}

func authenticateUser(username, password string) bool {
	var storedpswd string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedpswd)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
    fmt.Println("Error quering user: ", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedpswd), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username: ")
		return
	}
	username = strings.TrimSpace(username)

	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password: ")
		return
	}
	password = strings.TrimSpace(password)

	if userExists(username) {
		if authenticateUser(username, password){
			client := &connection{conn: conn, username: username}
			clientsMutex.Lock()
			clients[conn] = client
			clientsMutex.Unlock()

			handleBroadcast(conn, fmt.Sprintf("User %s Joined!", client.username))

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				update := scanner.Text()
				if update == "" {
					continue
				}
				handleBroadcast(conn, fmt.Sprintf("%s: %s", client.username, update))
			}
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()

			handleBroadcast(nil, fmt.Sprintf("%s has left the chat.\n", client.username))
		}else{
    fmt.Fprintln(conn, "Invalid Credentials")
    }
  }else{
    fmt.Fprintln(conn, "User doesn't exist, would you like to create a new ID? [yes/y/no/n]")
    response, err := reader.ReadString('\n')
    if err != nil{
      fmt.Println("Error reading response: ", err)
      return
    }
    response = strings.TrimSpace(response)
    if strings.ToLower(response) == "yes" || strings.ToLower(response) == "y"{
      err := createNewUser(username, password)
      if err != nil{
        fmt.Println("Error regsitering user: ", err)
        return
      }
      fmt.Fprintln(conn, "User Registered Successfully, Login Again")
    }else{
      fmt.Fprintln(conn, "invalid choice, exiting!")
    }
  }


}

func handleBroadcast(sender net.Conn, message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for conn, client := range clients {
		if conn != sender {
			_, err := fmt.Fprintf(client.conn, message+"\n")
			if err != nil {
				fmt.Println("Error sending message back to users: ", err)
			}
		}
	}
}

func main() {
  initDB()
  defer db.Close()
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening to server: ", err)
	}

	defer listener.Close()
	fmt.Println("Server started!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
