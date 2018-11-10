package main

import (
	"bufio"
	"fmt"
	//"io/ioutil"
	//"jutil"
	//"net/http"
	"os"
	//"strings"
	//"time"
)

func main() {

	inFiName := "data/simple-config/basic-prop.txt"
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("error opening ID file ", inFiName, " err=", err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		aline := scanner.Text()
		if len(aline) > 0 {
			fmt.Println(aline)
		}
	}
} //main
