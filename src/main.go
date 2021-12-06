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

		inputLine = strings.TrimSuffix(inputLine, "$")

		graphReceived := readString(inputLine)

		fmt.Printf("#DEBUG Le graph est enregistré\n")

		reponse := ""

		for i := range graphReceived.points{

			reponse = formulateAnswer(reponse, solveGraph(graphReceived, graphReceived.points[i]),graphReceived.points)

		}

        io.WriteString(connection, fmt.Sprintf("%s$", reponse))
	}

}

func solveGraph(toSolve graph, sommet string) []link {

	fmt.Printf("#DEBUG Debut de la résolution en partant de %d\n", sommet)

	indexSommet := findIndex(toSolve.points,sommet)

	fmt.Printf("#DEBUG Sommet trouvé à la place numero  %d\n", indexSommet)

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

	//fmt.Printf("-----------------------------------------------------------\n", index,sommetRelies)

	for index < len(toSolve.points){

		//fmt.Printf("#DEBUG Les sommets reliés a %[1]d iteration sont %[2]d\n", index,sommetRelies)

		ret[index-1] = getLowestLink(toSolve.arretes, sommetRelies, coutsTot)

		//fmt.Printf("#DEBUG Lien qui va bien : %[1]d\n", ret[index-1])

		sommetRelies[index] = ret[index - 1].sommetDeux
		coutsTot[index] = correspondingWeight(sommetRelies, coutsTot, ret[index-1].sommetUn) + ret[index - 1].poids

		index = index + 1

		//fmt.Printf("-----------------------------------------------------------\n", index,sommetRelies)

	}

	fmt.Printf("Les sommets retenus sont : %d \n", ret)

	return ret
}

func getLowestLink(tab [][]int, done []int, couts []int) link {

	var ret link

	indexUn := -1
	indexDeux := -1
	value := -1

	for i := range done{

		if (done[i]>=0){

			//fmt.Printf("----- Je cherches parmis la colonne de %d sommet\n", done[i])

			toTest := done[i]

			for j := range tab[toTest]{


				//fmt.Printf("--------------- Je cherches parmis la ligne de %d sommet\n", j)


				if (value == -1){
					if ((!contain(done,j)) && (tab[toTest][j] != 0)){
						//fmt.Printf("--------------- Lien trouvé allant de %d  vers %d \n",toTest, j)
						indexUn = toTest
						indexDeux = j
						value = tab[toTest][j]
					}
				}else{
					if (tab[toTest][j] != 0) && (!contain(done, j)) && (correspondingWeight(done, couts, indexUn) + value > correspondingWeight(done, couts, toTest)+tab[toTest][j]){
						//fmt.Printf("--------------- Lien trouvé allant de %d  vers %d \n",toTest, j)
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

	//readString("Graph1\nA,B,C,D\n{0,40,5,6}\n{10,0,5,9}\n{5,5,0,2}\n{6,9,2,0}")

	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	ln, err := net.Listen("tcp", portString)
	check(err)

	connum := 0

	for {
		conn, errconn := ln.Accept()
		check(errconn)

		connum += 1
		fmt.Println("Client connecté\n")
		go handleConnection(conn, connum)
	}

}