package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
	"io"
	"strings"
	"regexp"
	"strconv"
)

type graph struct {
	nom string
	points []string
	arretes [][]int
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getPort() int{

	if len(os.Args)!=2 {
		fmt.Printf("utilisation : go run main.go <Port Number>")
		os.Exit(1)
	}else{
		portNumber, err := strconv.Atoi(os.Args[1])
		check(err)
		fmt.Printf("Port séléctionner : ",portNumber)
		return portNumber
	}

	return -1

}

func (item graph) toString() string{
	return item.nom
}

func readString(maString string) graph{

	re := regexp.MustCompile("\\{(.*)\\}")
	nbArretes := re.FindAllString(maString,-1)

	traits := make([][]int, len(nbArretes))
	for i := range traits {
    	traits[i] = make([]int, len(nbArretes))
	}

	for elm := range nbArretes {
		test1 := strings.Replace(nbArretes[elm], "{", "", -1)
		test2 := strings.Replace(test1, "}", "", -1)
		test3 := strings.Split(test2, ",")
		
		for elm1 := range test3 {
			intVar, err := strconv.Atoi(test3[elm1])
			check(err)
			traits[elm][elm1] = intVar
		}
	}

	lignes := strings.Split(maString, "\n")
	nomGraph := lignes[0]
	nomSommet := strings.Split(lignes[1],",")

	var ret graph
	ret.nom = nomGraph
	ret.points = nomSommet
	ret.arretes = traits

	return ret
}

func handleConnection(connection net.Conn, connum int) {

	defer connection.Close()
	
	connReader := bufio.NewReader(connection)

	for {
		inputLine, err := connReader.ReadString('$')
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
            fmt.Printf("Error :|%s|\n", err.Error())
			break
		}

		inputLine = strings.TrimSuffix(inputLine, "$")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)
        splitLine := strings.Split(inputLine, " ")
        returnedString := splitLine[len(splitLine)-1]
        fmt.Printf("#DEBUG %d RCV Returned value |%s|\n", connum, returnedString)
        io.WriteString(connection, fmt.Sprintf("%s\n", returnedString))
	}

}


func main()  {

	port := getPort()

	//readString("Graph1\nA,B,C,D\n{0,40,5,6}\n{10,0,5,9}\n{5,5,0,2}\n{6,9,2,0}")

	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	ln, err := net.Listen("tcp", portString)
	check(err)

	connum := 0

	for {
		conn, errconn := ln.Accept()
		check(errconn)

		connum += 1
		fmt.Println("Client connecté")
		go handleConnection(conn, connum)
	}

}