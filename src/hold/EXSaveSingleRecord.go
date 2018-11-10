package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mdsClient"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("EXSaveSingleRecord.go")
	postStr := "{'id' : 1918171, 'title', 'this is a test structure'}"
	postBody := bytes.NewBuffer([]byte(postStr)) // Convert post body String into form compatible with stream writter
	uri := "http://127.0.0.1:9601/mds/1918171"
	fmt.Println("uri=", uri)
	start := time.Now().UTC()

	hc := http.Client{}
	req, err := http.NewRequest("PUT", uri, postBody)
	if err != nil {
		fmt.Println("Error: opening uri=", uri, " err=", err)
	}

	// Add the content type header
	// when saving the document so we can get
	// it back when retrieving the document latter.
	req.Header.Set("Content-Type", "application/json")

	// By convention any headers with prefix of "meta"
	// will be saved on server and returned with the
	// document when fetched.
	req.Header.Set("meta-roles-view", "joe,jim,admin")
	req.Header.Set("meta-roles-delete", "admin,mgr1")
	req.Header.Set("meta-title", "Sample Test Structure")
	req.Header.Set("meta-purpose", "Test meta data save and retrieval")

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

	mdsClient.TimeTrack(os.Stdout, start, "finished single PUT")
	fmt.Println("statusCode=", resp.StatusCode)
	fmt.Println("body=", string(body))
}
