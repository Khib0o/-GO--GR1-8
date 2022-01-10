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

func (g graph) toString() string{
	str := fmt.Sprintf("Le nom du graph est :\n"+g.nom+"\n Ses sommets sont :\n")
	for _,elm := range g.points{
		str = str + elm + " "
	}
	str = str + "\n"
	return str
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

	return -1

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
            fmt.Printf("Error :|%s|\n", err.Error())
			fmt.Printf("--Client %d disconnected--", connum)
			break
		}

		inputLine = strings.TrimSuffix(inputLine, "$")


		graphReceived := readString(inputLine)

		reponse := ""

		for i := range graphReceived.points{

			reponse = formulateAnswer(reponse, solveGraph(graphReceived, graphReceived.points[i]),graphReceived.points)

		}

        io.WriteString(connection, fmt.Sprintf("%s$", reponse))
	}

}

func solveGraph(toSolve graph, sommet string) []link {

	indexSommet := findIndex(toSolve.points,sommet)

	sommetRelies := make([]int, len(toSolve.points))
	coutsTot := make([]int, len(toSolve.points))

	for i := range sommetRelies{
		sommetRelies[i] = -1
		coutsTot[i] = -1
	}
	index := 0

	ret := make([]link, len(toSolve.points)-1)

	sommetRelies[index] = indexSommet
	coutsTot[index] = 0
	index = index + 1

	for index < len(toSolve.points){

		ret[index-1] = getLowestLink(toSolve.arretes, sommetRelies, coutsTot)

		sommetRelies[index] = ret[index - 1].sommetDeux
		coutsTot[index] = correspondingWeight(sommetRelies, coutsTot, ret[index-1].sommetUn) + ret[index - 1].poids

		index = index + 1


	}


	return ret
}

func getLowestLink(tab [][]int, done []int, couts []int) link {

	var ret link

	indexUn := -1
	indexDeux := -1
	value := -1

	for i := range done{

		if (done[i]>=0){


			toTest := done[i]

			for j := range tab[toTest]{

				if (value == -1){
					if ((!contain(done,j)) && (tab[toTest][j] != 0)){
						indexUn = toTest
						indexDeux = j
						value = tab[toTest][j]
					}
				}else{
					if (tab[toTest][j] != 0) && (!contain(done, j)) && (correspondingWeight(done, couts, indexUn) + value > correspondingWeight(done, couts, toTest)+tab[toTest][j]){
						indexUn = toTest
						indexDeux = j
						value = tab[toTest][j]
					}
				}

			}

		}
	}

	ret.sommetUn = indexUn
	ret.sommetDeux = indexDeux
	ret.poids = value

	return ret
}

func correspondingWeight(sommets []int, weight []int, sommet int) int {

	for i:= range sommets{
		if (sommets[i] == sommet){
			return weight[i]
		}
	}
	return -1

}

func findIndex(self []string, value string) int{
	for p,v := range self{
		if (v == value){
			return p
			p = p + 1
		}
	}
	return -1
}

func  contain(self []int, value int) bool{
	for p,v := range self{
		if (v == value){
			return true
			p = p + 1
		}
	}
	return false
}

func formulateAnswer(before string, toAdd []link, dico []string) string{

	str := before

	str += fmt.Sprintf("Point initial : %s\n", dico[toAdd[0].sommetUn])
	
	for i := range toAdd{

		str += dico[toAdd[i].sommetUn] + " -------->" + dico[toAdd[i].sommetDeux] + " (Poids : " +  strconv.Itoa(toAdd[i].poids) + ")\n"

	}

	return str
}

func main()  {

	port := getPort()

	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	ln, err := net.Listen("tcp", portString)
	check(err)

	connum := 0

	for {
		conn, errconn := ln.Accept()
		check(errconn)

		connum += 1
		fmt.Println("\n --Client connecté --\n")
		go handleConnection(conn, connum)
	}

}