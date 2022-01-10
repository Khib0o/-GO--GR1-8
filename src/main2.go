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

type infoCon struct {
	connexion net.Conn
	numero int
}

type problem struct {
	connexion net.Conn
	graphique graph
	sommetDepart string
	nbSommets int
	numero int
}

type solution struct {
	connexion net.Conn
	reponse string
	nbSommets int
	numero int
}

const WORKER int = 10

func main()  {

	port := getPort()

	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	ln, err := net.Listen("tcp", portString)
	check(err)

	connum := 0

	var channelConnexion chan infoCon
	channelConnexion = make(chan infoCon, 10)

	var channelProblem chan problem
	channelProblem = make(chan problem, 10)

	var channelSolution chan solution
	channelSolution = make(chan solution, 10)

	go handleConnection(channelConnexion, channelProblem)
	go jeSuisLeDernierMaillon(channelSolution)

	for compteur := 0 ; compteur < WORKER ; compteur ++{
		go jeSuisUnWorker(channelProblem, channelSolution, compteur)
	}

	for {
		conn, errconn := ln.Accept()
		check(errconn)
		connum += 1

		var conne infoCon
		conne.numero = connum
		conne.connexion = conn

		fmt.Println("\n --Client connecté --\n")
		channelConnexion <- conne
	}

}

func handleConnection(inp chan infoCon, out chan problem) {
	
	for {

		conn := <- inp
		
		connReader := bufio.NewReader(conn.connexion)

		for {
			inputLine, err := connReader.ReadString('$')
			if err != nil {
	            fmt.Printf("Error :|%s|\n", err.Error())
				fmt.Printf("--Client disconnected--")
				break
			}

			inputLine = strings.TrimSuffix(inputLine, "$")
			graphReceived := readString(inputLine)

			for i := range graphReceived.points{

				var toPush problem
				toPush.connexion = conn.connexion
				toPush.sommetDepart = graphReceived.points[i]
				toPush.graphique = graphReceived
				toPush.nbSommets = len(graphReceived.points)
				toPush.numero = conn.numero
				out <- toPush


			}


		}

	}

}

func jeSuisUnWorker(inp chan problem, out chan solution, num int){
	
	for{

		prob := <- inp
		reponse := ""
		reponse = formulateAnswer(reponse, solveGraph(prob.graphique, prob.sommetDepart),prob.graphique.points)

		var soluce solution
		soluce.reponse = reponse
		soluce.connexion = prob.connexion
		soluce.nbSommets = prob.nbSommets
		soluce.numero = prob.numero

		out <- soluce
	}

}

func jeSuisLeDernierMaillon(inp chan solution){

	tableauCompte := make([]int, 10)
	tableauRep := make([]string, 10)

	for i := 0; i<10; i++{
		tableauCompte[i] = 0
		tableauRep[i] = ""
	}

	for{
		
		soluce := <- inp
		if (tableauCompte[soluce.numero%10]==0){

			tableauCompte[soluce.numero%10] = 1
			tableauRep[soluce.numero%10] += soluce.reponse

		}else{
			if (tableauCompte[soluce.numero%10]==soluce.nbSommets-1){

				tableauRep[soluce.numero%10] += soluce.reponse
				io.WriteString(soluce.connexion, fmt.Sprintf("%s$", tableauRep[soluce.numero%10]))
				soluce.connexion.Close()
				tableauCompte[soluce.numero%10] = 0
				tableauRep[soluce.numero%10] = ""

			}else{
				tableauCompte[soluce.numero%10]+= 1
				tableauRep[soluce.numero%10] += soluce.reponse
			}
		}

	}

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