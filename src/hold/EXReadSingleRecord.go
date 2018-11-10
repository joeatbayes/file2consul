package main

import (
	"fmt"
	"io/ioutil"
	"mdsClient"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("EXReadSingleRecord.go")
	uri := "http://127.0.0.1:9601/mds/1918171"
	fmt.Println("uri=", uri)
	start := time.Now().UTC()
	hc := http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Error: opening uri=", uri, " err=", err)
		return
	}

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("Error: uri=", uri, " err=", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error: Expected Status Code 200, Got ", resp.StatusCode)
	} else {
		fmt.Println("Sucess:")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	mdsClient.TimeTrack(os.Stdout, start, "finished single GET")
	fmt.Println("statusCode=", resp.StatusCode)
	fmt.Println(" body=", string(body))

}
