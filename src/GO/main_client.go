package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

	if len(os.Args) != 6 {
		fmt.Printf("Usage: go run main_client.go <Numero de port> <Fichier> <Ligne du graph> <Afficher résultat (y/n)> <Afficher temps (y/n)>\n")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])		//Récupération du port
	check(err)

	file, err := os.Open(os.Args[2])
	check(err)

	reader := bufio.NewReader(file)
	compteur := 0
	nbG, err := strconv.Atoi(os.Args[3])
	check(err)

	for {										//Lecture de la ligne visée
		test, err := reader.ReadString('$')
		check(err)
		compteur = compteur + 1

		if compteur == nbG {
			graph_string = test
			graph_string = strings.Replace(graph_string, "\\n", "\n", -1)
			break
		}

	}

	graph_string = strings.TrimPrefix(graph_string, "\n")

	file.Close()

	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))

	conn, err := net.Dial("tcp", portString)

	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {

		defer conn.Close()
		reader := bufio.NewReader(conn)

		io.WriteString(conn, fmt.Sprintf(graph_string))		//Envoie de la chaine correspondante au graphique

		mesureTemps := time.Now()
		resultString, err := reader.ReadString('$')			//Reception de la solution
		if err != nil {
			fmt.Printf("DEBUG MAIN could not read from server")
			os.Exit(1)
		}
		elapsed := time.Since(mesureTemps)
		resultString = strings.TrimSuffix(resultString, "$")
		if os.Args[4] == "y" {
			fmt.Printf("%s", resultString)						//Affichage du résultat
		}
		if os.Args[5] == "y" {
			fmt.Printf("\nTime of execution : %s\n", elapsed)	//Lecture du temps d'execution
		}

	}

}
