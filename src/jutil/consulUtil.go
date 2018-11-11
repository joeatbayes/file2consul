package jutil

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func SetConsulKey(serverURI string, key string, val string) {
	fmt.Println("setConsulKey key=", key, " value=", val)
	uri := serverURI + "/v1/kv/" + key
	//fmt.Println("uri=", uri)
	start := time.Now().UTC()
	hc := http.Client{}
	req, err := http.NewRequest("PUT", uri, strings.NewReader(val))
	if err != nil {
		fmt.Println("Error: opening uri=", uri, " err=", err, " key=", key, "  val=", val)
		return
	}
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("Error: uri=", uri, " err=", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error: Expected Status Code 200, Got ", resp.StatusCode)
	} else {
		//fmt.Println("Sucess:")
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	TimeTrack(os.Stdout, start, "finished single PUT uri="+uri+"\n")
	//fmt.Println("statusCode=", resp.StatusCode)
	//fmt.Println(" body=", string(body))
}
