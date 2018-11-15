package jutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ConsulKeys struct {
	Keys []string
}

func GetConsulKeys(serverURI string, key string, sep string, dataCenter string) []string {
	uri := serverURI + "/v1/kv/" + key + "?keys&sep" + url.QueryEscape(sep)
	if dataCenter > "" {
		uri = uri + "&dc=" + url.QueryEscape(dataCenter)
	}
	//fmt.Println("uri=", uri)
	start := time.Now().UTC()
	hc := http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Error: GetConsulKeys:opening uri=", uri, " err=", err, " key=", key)
		return nil
	}
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("Error: uri=", uri, " err=", err)
		return nil
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error: GetConsulKeys: Expected Status Code 200, Got ", resp.StatusCode)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("statusCode=", resp.StatusCode)
	fmt.Println(" body=", string(body))
	bodystr := "{ \"Keys\":" + string(body) + "}"
	fmt.Println(" bodystr=", bodystr)

	// Parse Results from JSON Array
	var m ConsulKeys
	err = json.Unmarshal([]byte(bodystr), &m)
	fmt.Println(" after unmarshal m=", m)
	TimeTrack(os.Stdout, start, "finished single PUT uri="+uri+"\n")
	//fmt.Println("statusCode=", resp.StatusCode)
	//fmt.Println(" body=", string(body))
	return m.Keys

}

// Save specified key to consul.  report error if fails.
func SetConsulKey(serverURI string, key string, val string, dataCenter string) {
	fmt.Println("setConsulKey key=", key, " value=", val)
	uri := serverURI + "/v1/kv/" + key
	if dataCenter > "" {
		uri = uri + "?dc=" + url.QueryEscape(dataCenter)
	}
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
