/* EXReadSingleMultiThreaded.go - keep here showing multi-threaded
reader sending files from the input file to a remote server */
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"mdsClient"
	"net/http"
	"os"
	"strings"
	"time"
)

type MTReader struct {
	perf *mdsClient.PerfMeasure
}

func makeMTReader() *MTReader {
	r := MTReader{}
	r.perf = mdsClient.MakePerfMeasure(100000)
	return &r
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func (u *MTReader) procLine(baseURI string, line string) {

	id := strings.TrimSpace(line)
	if len(id) < 1 {
		fmt.Println("L42: Empty Line")
		return
	}
	u.perf.NumReq += 1
	//fmt.Println("aline=", aline)

	uri := "http://127.0.0.1:9601/mds/" + id

	hc := http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("error opening uri=", uri, " err=", err)
	}

	req.Header.Set("Connection", "close")
	req.Close = true
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("err=", err, "resp=", resp)
		return
	}
	defer resp.Body.Close()

	//fmt.Println(" reps=", resp, "err=", err)
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		u.perf.NumSuc += 1
	} else {
		u.perf.NumFail += 1
	}
	u.perf.NumSinceStatPrint += 1
	u.perf.CheckAndPrintStat(os.Stdout)
	u.perf.BytesRead += int64(len(body))
	//fmt.Println(" body=", string(body))
	//len(body)
	//defer mdsClient.TimeTrack(now, "finished id=" + id)
	time.Sleep(1500)

}

var (
	server *http.Server
	client *http.Client
)

func main() {
	const procs = 2
	const MaxWorkerThread = 12 //5 //150 //5 //15 // 50 // 350
	const MaxQueueSize = 3000
	const BaseURI = "http://127.0.0.1:9601/mds/"

	fmt.Println("EXReadSingleMultiThreaded.go")
	idFiName := "data/sample/physicians.id.txt"

	fmt.Println("idFiName=", idFiName, " baseURI=", BaseURI)
	fmt.Println("MaxWorkerThread=", MaxWorkerThread)
	fmt.Println("MaxQueueSize=", MaxQueueSize)
	u := makeMTReader()

	start := time.Now().UTC()
	idFile, err := os.Open(idFiName)
	if err != nil {
		fmt.Println("error opening ID file ", idFiName, " err=", err)
	}
	defer idFile.Close()

	scanner := bufio.NewScanner(idFile)
	linesChan := make(chan string, MaxQueueSize)
	done := make(chan bool)

	// Spin up 100 worker threads to post
	// content to the server.
	for i := 0; i < MaxWorkerThread; i++ {
		go func() {
			for {
				aline, more := <-linesChan
				if more {
					u.procLine(BaseURI, aline)

				} else {
					done <- true
					return
				}
			}
		}()
	}

	// Add the rows to the Queue
	// so we can process them in parralell
	// It is blocked at MaxQueueSize by the
	// channel size.
	for scanner.Scan() {
		aline := scanner.Text()
		if len(aline) > 0 {
			// Add each line multiple times
			// to provide more work and exercise
			// cache
			linesChan <- aline
			linesChan <- aline
			linesChan <- aline
		}
	}
	close(linesChan)
	fmt.Println("Queued all jobs")
	mdsClient.TimeTrackMin(os.Stdout, start, "finished Queing")
	<-done // wait until queue has been marked as finished.
	fmt.Println("Sent All Jobs")
	mdsClient.TimeTrackMin(os.Stdout, start, "inserted all records")
	u.perf.PrintStat(os.Stdout)
}
