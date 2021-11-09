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

type link struct {
	sommetUn int
	sommetDeux int
	poids int
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
		fmt.Printf("Port séléctionné : ",portNumber)
		return portNumber
	}

	return -1  //should not be reached

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
		fmt.Printf("Line detected",inputLine)

		inputLine = strings.TrimSuffix(inputLine, "$")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)
        fmt.Printf(inputLine)
        io.WriteString(connection, fmt.Sprintf("J'ai fini de capter \n"))
	}

}

func solveGraph(toSolve graph, sommet string) []link {

	indexSommet := toSolve.points.findIndex(sommet)
	sommetRelies := make([]int, len(toSolve.points))
	index := 0

	ret := make([]link, len(toSolve.points)-1)

	sommetRelies[index] = indexSommet
	index = index + 1

	for index < len(toSolve.points){
		ret[index-1] = getLowestLink(graph.traits, sommetRelies)
		if sommetRelies.contain(ret[index-1].sommetUn){
			sommetRelies[index] = ret[index-1].sommetDeux
		}else{
			sommetRelies[index] = ret[index-1].sommetUn
		}
	}

	return ret
}

func getLowestLink(tab [][]int, done []int) link {

	var ret link

	indexUn := -1
	indexDeux := -1
	value := -1

	for i := range tab{
		for j := i; j < len(tab); j++{
			if !(done.contain(i))&&!(done.contain(j)){
				continue
			}
			if ((tab[i][j] <= value)||(value < 0)){
				indexUn = i
				indexDeux = j
				value = tab[i][j]
			}
		}
	}

	ret.sommetUn = indexUn
	ret.sommetDeux = indexDeux
	ret.poids = value

	return ret
}

func (self []string) findIndex(value string) int{
	for p,v := range self{
		if (v == value){
			return p
		}
	}
	return -1
}

func (self []int) contain(value int) bool{
	for p,v := range self{
		if (v == value){
			return true
		}
	}
	return false
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