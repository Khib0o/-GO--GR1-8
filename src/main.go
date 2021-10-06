package main

import (
	"fmt"
	"os"
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

func (item graph) toString() string{
	return item.nom
}

func readFile() {
	dat, err := os.ReadFile("../entree.txt")
	check(err)
	str := strings.Split(string(dat), "!")
	var graphs [10]graph
	for elm := range str{
		if elm == 0 {
			continue
		}
		stri := []string{"A","B","C"}
		numb := [][]int{{1,1},{1,1}}
		graphs[elm-1] = graph{"test", stri,  numb}
		fmt.Println(graphs[elm-1].toString())
	}
	re := regexp.MustCompile("\\{(.*)\\}")
	sArretes := re.FindAllString(str[1],-1)
	traits := make([][]int, len(sArretes))
	for i := range traits {
    	traits[i] = make([]int, len(sArretes))
	}
	for elm := range sArretes {
		test1 := strings.Replace(sArretes[elm], "{", "", -1)
		test2 := strings.Replace(test1, "}", "", -1)
		test3 := strings.Split(test2, ",")
		fmt.Println("------------------")
		fmt.Println(test3)
		for elm1 := range test3 {
			intVar, _ := strconv.Atoi(test3[elm1])
			traits[elm][elm1] = intVar
		}
	}
	fmt.Println("------------------")
	fmt.Println(traits)
}

func main()  {
	readFile()
	
}