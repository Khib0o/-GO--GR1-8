package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type graph struct {
	nom     string
	points  []string
	arretes [][]int
}

type link struct {
	sommetUn   int
	sommetDeux int
	poids      int
}

type infoCon struct {
	connexion net.Conn
	numero    int
}

type problem struct {
	connexion    net.Conn
	graphique    graph
	sommetDepart string
	nbSommets    int
	numero       int
}

type solution struct {
	connexion net.Conn
	reponse   string
	nbSommets int
	numero    int
}

const WORKER int = 10		//Nombre de goroutine travaillant à la résolution de problèmes

func main() {

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

	go traitementConnexions(channelConnexion, channelProblem)
	go formulateurDeReponse(channelSolution)

	for compteur := 0; compteur < WORKER; compteur++ {
		go worker(channelProblem, channelSolution, compteur)
	}

	for {
		conn, errconn := ln.Accept()
		check(errconn)
		connum += 1

		var conne infoCon
		conne.numero = connum
		conne.connexion = conn

		channelConnexion <- conne
	}

}

//Gère les connexions entrantes
func traitementConnexions(inp chan infoCon, out chan problem) {

	for {

		conn := <-inp

		connReader := bufio.NewReader(conn.connexion)

		for {
			inputLine, err := connReader.ReadString('$')
			if err != nil {
				break
			}

			inputLine = strings.TrimSuffix(inputLine, "$")
			graphReceived := readString(inputLine)

			for i := range graphReceived.points {

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

//Résoud un sous-problème
func worker(inp chan problem, out chan solution, num int) {

	for {

		prob := <-inp
		reponse := ""
		reponse = formulateAnswer(reponse, solveGraph(prob.graphique, prob.sommetDepart), prob.graphique.points)

		var soluce solution
		soluce.reponse = reponse
		soluce.connexion = prob.connexion
		soluce.nbSommets = prob.nbSommets
		soluce.numero = prob.numero

		out <- soluce
	}

}

//Une fois toutes les solutions d'un problème arrivée, envoie la réponse au client
func formulateurDeReponse(inp chan solution) {

	tableauCompte := make([]int, 10)
	tableauRep := make([]string, 10)

	for i := 0; i < 10; i++ {
		tableauCompte[i] = 0
		tableauRep[i] = ""
	}

	for {

		soluce := <-inp
		index := soluce.numero % 10
		if tableauCompte[index] == 0 {

			tableauCompte[index] = 1
			tableauRep[index] += soluce.reponse

		} else {
			if tableauCompte[index] == soluce.nbSommets-1 {

				tableauRep[index] += soluce.reponse
				io.WriteString(soluce.connexion, fmt.Sprintf("%s$", tableauRep[index]))
				soluce.connexion.Close()
				tableauCompte[index] = 0
				tableauRep[index] = ""

			} else {
				tableauCompte[index] += 1
				tableauRep[index] += soluce.reponse
			}
		}

	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getPort() int {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : go run main_server.go <Numero de port>")
		os.Exit(1)
	} else {
		portNumber, err := strconv.Atoi(os.Args[1])
		check(err)
		fmt.Printf("Port séléctionné : ", portNumber)
		return portNumber
	}

	return -1

}

//Transforme une string en un objet de type graph
func readString(maString string) graph {

	re := regexp.MustCompile("\\{(.*)\\}")
	nbArretes := re.FindAllString(maString, -1)

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
	nomSommet := strings.Split(lignes[1], ",")

	var ret graph
	ret.nom = nomGraph
	ret.points = nomSommet
	ret.arretes = traits

	return ret
}

//Djikstra : Résoud un graphique à partir d'un sommet donné
func solveGraph(toSolve graph, sommet string) []link {

	indexSommet := findIndex(toSolve.points, sommet)

	sommetRelies := make([]int, len(toSolve.points))
	coutsTot := make([]int, len(toSolve.points))

	for i := range sommetRelies {
		sommetRelies[i] = -1
		coutsTot[i] = -1
	}
	index := 0

	ret := make([]link, len(toSolve.points)-1)

	sommetRelies[index] = indexSommet
	coutsTot[index] = 0
	index = index + 1

	for index < len(toSolve.points) {

		ret[index-1] = getLowestLink(toSolve.arretes, sommetRelies, coutsTot)

		sommetRelies[index] = ret[index-1].sommetDeux
		coutsTot[index] = correspondingWeight(sommetRelies, coutsTot, ret[index-1].sommetUn) + ret[index-1].poids

		index = index + 1

	}

	return ret
}

// Regarde le lien le plus interessant à retenir pour le prochain coup
func getLowestLink(tab [][]int, done []int, couts []int) link { 

	var ret link 

	indexUn := -1
	indexDeux := -1
	value := -1

	for i := range done {

		if done[i] >= 0 {

			toTest := done[i]

			for j := range tab[toTest] {

				if value == -1 {
					if (!contain(done, j)) && (tab[toTest][j] != 0) {
						indexUn = toTest
						indexDeux = j
						value = tab[toTest][j]
					}
				} else {
					if (tab[toTest][j] != 0) && (!contain(done, j)) && (correspondingWeight(done, couts, indexUn)+value > correspondingWeight(done, couts, toTest)+tab[toTest][j]) {
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

//Permet de récupérer le poids pour atteindre un sommet donné
func correspondingWeight(sommets []int, weight []int, sommet int) int {

	for i := range sommets {
		if sommets[i] == sommet {
			return weight[i]
		}
	}
	return -1

}

//Renvoie l'index d'un sommmet donné
func findIndex(self []string, value string) int {
	for p, v := range self {
		if v == value {
			return p
			p = p + 1
		}
	}
	return -1
}

//Vérifie si un sommet est contenu dans les sommets déjà reliés
func contain(self []int, value int) bool {
	for p, v := range self {
		if v == value {
			p = p + 1
			return true
		}
	}
	return false
}

//Permet de formuler une réponse à partir du tableau des liens choisis
func formulateAnswer(before string, toAdd []link, dico []string) string {

	str := before

	str += fmt.Sprintf("Point initial : %s\n", dico[toAdd[0].sommetUn])

	for i := range toAdd {

		str += dico[toAdd[i].sommetUn] + " -------->" + dico[toAdd[i].sommetDeux] + " (Poids : " + strconv.Itoa(toAdd[i].poids) + ")\n"

	}

	return str
}
