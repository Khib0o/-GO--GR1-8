package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"bufio"
    "strings"
    "time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	graph_string := ""
	port := 10000

	if len(os.Args) != 3 {
		fmt.Printf("Usage: go run client.go <portnumber> <graph_string>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		check(err)
		port = portNumber
		graph_string = os.Args[2]
	}

	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {

        defer conn.Close()
        reader := bufio.NewReader(conn)
		fmt.Printf("#DEBUG MAIN connected\n")

		fmt.Sprintf(graph_string)

        io.WriteString(conn, fmt.Sprintf("salut gros fdp $ vas niquer tes morts"))

		fmt.Printf("#DEBUG MAIN message envoyé\n")
            
        resultString, err := reader.ReadString('\n')
        if (err != nil){
            fmt.Printf("DEBUG MAIN could not read from server")
            os.Exit(1)
        }
        resultString = strings.TrimSuffix(resultString, "\n")
        fmt.Printf("#DEBUG server replied : |%s|\n", resultString)
        time.Sleep(1000 * time.Millisecond)

	}

}
