package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"bufio"	
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	//graph_string := "Graph1\nA,B,C,D\n{0,2,5,3}\n{2,0,0,0}\n{5,0,0,4}\n{3,0,4,0}$"

	graph_string := "TestTD2\nA,B,C,D,E,F\n{0,2,100,0,0,0}\n{2,0,10,60,0,0}\n{100,10,0,0,3,2}}\n{0,60,0,0,4,0}\n{0,0,3,4,0,2}\n{0,0,2,0,1,0}\n$"

	port := 10000

	if len(os.Args) != 4 {
		fmt.Printf("Usage: go run client.go <portnumber> <file> <graph_number>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		check(err)
		port = portNumber
		file, err := os.Open(os.Args[2])
		check(err)

		reader := bufio.NewReader(file)
		compteur := 0
		nbG, err := strconv.Atoi(os.Args[3])
		check(err)
		for {
			test, err := reader.ReadString('$')
			check(err)
			compteur= compteur + 1
    		if (compteur == nbG){
				graph_string = test
				graph_string = strings.Replace(graph_string,"\\n","\n",-1)
				fmt.Printf(graph_string)
				fmt.Printf("\n")
				//fmt.Print("TestTD2\nA,B,C,D,E,F\n{0,2,100,0,0,0}\n{2,0,10,60,0,0}\n{100,10,0,0,3,2}}\n{0,60,0,0,4,0}\n{0,0,3,4,0,2}\n{0,0,2,0,1,0}\n$")
				break
			}
		}

		graph_string = strings.TrimPrefix(graph_string, "\n")

		file.Close()

		//graph_string = "TestTD2\nA,B,C,D,E,F\n{0,2,100,0,0,0}\n{2,0,10,60,0,0}\n{100,10,0,0,3,2}}\n{0,60,0,0,4,0}\n{0,0,3,4,0,2}\n{0,0,2,0,1,0}\n$"
		//graph_string = os.Args[2]
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

        io.WriteString(conn, fmt.Sprintf(graph_string))

		fmt.Printf("#DEBUG MAIN message envoyé\n")
            
        resultString, err := reader.ReadString('$')
        if (err != nil){
            fmt.Printf("DEBUG MAIN could not read from server")
            os.Exit(1)
        }
        resultString = strings.TrimSuffix(resultString, "$")
        fmt.Printf("#DEBUG server replied : -------------\n%s\n-------------\n", resultString)

	}

}
