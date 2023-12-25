package main

import (
    "fmt"
    "net"
    "strings"
    "strconv"
    "encoding/json"
)

type Connection struct {
    path string
    request_type string
    params map[string]string
}

var connections []Connection = make([]Connection, 0)

func get_path(request string) string {
    return strings.Split(request, " ")[1]
}

func get_request_type(request string) string {
    return strings.Split(request, " ")[0]
}

func get_params(request string) map[string]string {
    params := make(map[string]string)

    values := strings.Split(request, "\n")
    values = append(values[:0], values[1:]...)

    for _,value := range values {
        split := strings.Split(value, ":")

        if len(split) == 0 || len(split) < 2 {
            continue
        }

        params[split[0]] = strings.Replace(split[1], "\r", "",  -1)
    }

    return params
}

func answer(msg string, code int,  connection net.Conn) {
    new_msg := "HTTP/1.1 " + strconv.Itoa(code) + "\nContent-Type: text/plain\nContent-Length: " + strconv.Itoa(len(msg)) + "\n\n" + msg
    connection.Write([]byte(new_msg))
}

func handle_connection(connection net.Conn) {
    buffer := make([]byte, 10240)
    connection.Read(buffer)

    path := get_path(string(buffer))
    request_type := get_request_type(string(buffer))
    params := get_params(string(buffer))

    json_bytes, err := json.Marshal(params)

    if err != nil {
        return
    }

    answer(path + "\n" + request_type + "\n" + string(json_bytes), 200, connection)

    con := Connection {path, request_type, params}
    connections = append(connections, con)

    connection.Close()
}

func main() {
    listener, err := net.Listen("tcp", "127.0.0.1:8080")

    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Listening to 127.0.0.1:8080 for requests")

    for {
        connection, err := listener.Accept()

        if err != nil {
            continue
        }

        go handle_connection(connection)
    }
}
