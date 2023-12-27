package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type User struct {
  ip net.Addr
  name string
  connected int64
}

var users []User
var connections []net.Conn

func containsUser(ip net.Addr) bool {
  found := false

  for _, us := range users {
    found = us.ip == ip
  }

  return found
}

func containsUserName(name string) bool {
  found := false
  
  for _, us := range users {
    found = us.name == name
  }

  return found
}

func getHeaders(req string) map[string]string {
  headers := make(map[string]string, 0)

  vals := strings.Split(req, "\n")
  vals = append(vals[:0], vals[1:]...)

  for _, val := range vals {
    split := strings.Split(val, ":")

    if len(split) == 0 || len(split) < 2 {
      continue
    }

    headers[split[0]] = split[1]
  }

  return headers
}

func containsHeader(headers map[string]string, header string) bool {
  found := false

  for key, _ := range headers {
    found = header == key
  }

  return found
}

func getHeader(headers map[string]string, header string) string {
  head := ""

  for key, _ := range headers {
    if key == header {
      head = key
    }
  }

  return head
}

func handleConnection(conn net.Conn) {
  if containsUser(conn.LocalAddr()) {
    conn.Write([]byte("You are already connected to the server"))
    return 
  }

  var user User
  user.ip = conn.LocalAddr()
  user.connected = time.Now().UnixNano() / int64(time.Millisecond)

  buffer := make([]byte, 1024)
  conn.Read(buffer)

  if !containsHeader(getHeaders(string(buffer)), "NAME") {
    conn.Write([]byte("NAME header is missing"))
    return
  }

  user.name = getHeader(getHeaders(string(buffer)), "NAME")

  if !containsHeader(getHeaders(string(buffer)), "MSG") {
    conn.Write([]byte("MSG header is missing"))
    return
  }

  if containsUserName(user.name) {
    conn.Write([]byte("Username is already taken"))
    return
  }

  users = append(users, user)

  log.Printf("New connection (%s)", conn.LocalAddr().String())
  
  for _, con := range connections {
    con.Write([]byte(user.name + ": " + getHeader(getHeaders(string(buffer)), "MSG")))
  }

  log.Printf("Received a message from %s %s (MSG: %s)", user.ip.String(), user.name, getHeader(getHeaders(string(buffer)), "MSG"))
}

func getPort(config map[string]interface{}) string {
  port := "5454"

  for k, v := range config {
    if k == "port" {
      port = string(v.(string))
    }
  }

  return port
}

func main() {
  configData, err := os.ReadFile("config.json")

  if err != nil {
    log.Printf("Error occured while trying to read config.json (%s)", err)
    return
  }

  config := make(map[string]interface{}, 0)

  json.Unmarshal(configData, &config)

  port := getPort(config)

  server, err := net.Listen("tcp", ":" + port)
  users = make([]User, 0)
  connections = make([]net.Conn, 0)

  if err != nil {
    log.Fatalf("Unable to bind a TCP server to port %s (%s)", port, err)
    return 
  }

  for {
    conn, err := server.Accept()

    if err != nil {
      log.Printf("Unable to accept connection (%s)", err)
      return
    }

    connections = append(connections, conn)

    go handleConnection(conn)
  }
}
