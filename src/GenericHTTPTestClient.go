package main

// GenericHTTPClient.go
// See:  GenericHTTPTestClient.md
// See: ../data/sample/GenericTestSample.txt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type MTReader struct {
	perf       *jutil.PerfMeasure
	reqPending int
	logFile    *os.File
}

func makeMTReader(outFiName string) *MTReader {
	r := MTReader{}
	r.perf = jutil.MakePerfMeasure(25000)

	logFiName := outFiName
	var logFile, sferr = os.Create(logFiName)
	if sferr != nil {
		fmt.Println("Can not open log file ", logFiName, " sferr=", sferr)
	}
	r.logFile = logFile
	return &r
}

func (r *MTReader) done() {
	defer r.logFile.Close()
}

type TestSpec struct {
	Id        string
	Verb      string
	Uri       string
	Headers   map[string]string
	Expected  int
	Rematch   string
	ReNoMatch string
	Message   string
	Body      string
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func (u *MTReader) procLine(spec *TestSpec) {
	u.reqPending += 1
	u.perf.NumReq += 1
	//fmt.Println("L45: spec=", spec)
	//fmt.Println("L49: spec.Rematch=", spec.Rematch)
	uri := spec.Uri
	hc := http.Client{}
	reqStat := true
	errMsg := ""

	req, err := http.NewRequest(spec.Verb, uri, bytes.NewBuffer([]byte(spec.Body)))
	//fmt.Println("L50: req=", req, " err=", err)
	if err != nil {
		u.perf.NumFail += 1
		fmt.Fprintln(u.logFile, "FAIL: L74: id=", spec.Id, " message=", spec.Message, " error opening uri=", uri, " err=", err)
		u.reqPending--
		reqStat = false
		return
	}

	req.Header.Set("Connection", "close")
	req.Close = true
	resp, err := hc.Do(req)
	//fmt.Println(" L60: reps=", resp, "err=", err)
	if err != nil {
		fmt.Fprintln(u.logFile, "FAIL: L85: id=", spec.Id, " message=", spec.Message, "err=", err, "resp=", resp)
		u.perf.NumFail += 1
		u.reqPending--
		reqStat = false
		return
	}
	defer resp.Body.Close()

	//fmt.Println(" L68: reps=", resp, "err=", err)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if resp.StatusCode != spec.Expected {
		errMsg = fmt.Sprintln("\tL:97 Expected StatusCode=", spec.Expected, " got=", resp.StatusCode)
		reqStat = false
	}

	// Add Logic to check the RE Pattern
	// match for body result
	if spec.Rematch != "" && spec.Rematch > " " {
		match, _ := regexp.MatchString(spec.Rematch, bodyStr)
		//fmt.Println("L86 match=", match, "merr=", merr)
		if match != true {
			errMsg = fmt.Sprintln(" L107:failed rematch=", spec.Rematch)
			reqStat = false
		}
	}

	// Add Logic to check the RE Pattern
	// match for body result
	if spec.ReNoMatch != "" && spec.ReNoMatch > " " {
		match, _ := regexp.MatchString(spec.ReNoMatch, bodyStr)
		//fmt.Println("118 match=", match, "ReNoMatch=", spec.ReNoMatch)
		if match == true {
			errMsg = fmt.Sprintln(" L120:FAIL ReNoMatch pattern found in record reNoMatch=", spec.ReNoMatch, " match=", match)
			reqStat = false
		}
	}

	//x1, x2 := regexp.MatchString(".*xfoo.*", "seafood")
	//fmt.Println("L95:X1=", x1, "x2=", x2)

	u.perf.NumSinceStatPrint += 1
	u.perf.CheckAndPrintStat(u.logFile)

	//fmt.Println("L82: body len=", len(string(body)))
	//fmt.Println("L82: body =", string(body))
	//len(body)
	//defer mdsClient.TimeTrack(now, "finished id=" + id)

	if reqStat == true {
		u.perf.NumSuc += 1
		fmt.Fprintln(u.logFile, "SUCESS: L125: id=", spec.Id, "\tmessage=", spec.Message)
	} else {
		u.perf.NumFail += 1
		fmt.Fprintf(u.logFile, "FAIL: L128: id=%v \tmessage=%s \terr Msg=%s  \tverb=%s uri=%s\n", spec.Id, spec.Message, errMsg, spec.Verb, spec.Uri)
	}
	u.reqPending--
	time.Sleep(1500)
}

var (
	server *http.Server
	client *http.Client
)

func PrintHelp() {
	fmt.Println("GenericHttpTestClient -in=InputFileName -out=OutputFileName -MaxThread=5")
	fmt.Println("  -in defaults to data/sample/GenericTestSample.txt")
	fmt.Println("  -out defaults to GTestx1.log.txt ")
	fmt.Println(" -MaxThread defaults to 20")
}

func main() {
	const procs = 2
	const DefMaxWorkerThread = 20 //5 //150 //5 //15 // 50 // 350
	const MaxQueueSize = 3000
	const BaseURI = "http://127.0.0.1:9601/mds/"
	const DefInFiName = "data/sample/GenericTestSample.txt"
	const DefOutFiName = "GTestx1.log.txt"
	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}
	MaxWorkerThread := parms.Ival("maxthread", int(DefMaxWorkerThread))
	fmt.Println(parms.String())
	inFiName := parms.Sval("input", DefInFiName)
	outFiName := parms.Sval("output", DefOutFiName)

	u := makeMTReader(outFiName)
	fmt.Fprintln(u.logFile, "GenericHTTPTestClient.go")

	fmt.Println("TestCaseFiName=", inFiName, " baseURI=", BaseURI)
	fmt.Println("OutFileName=", outFiName)
	fmt.Println("MaxWorkerThread=", MaxWorkerThread)
	fmt.Println("MaxQueueSize=", MaxQueueSize)

	start := time.Now().UTC()
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("error opening input file ", inFiName, " err=", err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	linesChan := make(chan TestSpec, MaxQueueSize)
	done := make(chan bool)

	// Spin up 100 worker threads to post
	// content to the server.

	for i := 0; i < MaxWorkerThread; i++ {
		go func() {
			for {
				spec, more := <-linesChan
				if more {
					u.procLine(&spec)
					//fmt.Println("L128 spec=", spec)

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
	var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Text()
		aline = strings.TrimSpace(aline)
		if len(aline) < 1 {
			continue
		} else if aline[0] == '#' {
			if strings.HasPrefix(aline, "#END") {
				recStr := strings.TrimSpace(b.String())
				//fmt.Println("RecStr=", recStr)
				if len(recStr) > 0 {
					b.Reset()
					spec := TestSpec{}
					err := json.Unmarshal([]byte(recStr), &spec)
					//fmt.Println("L150: spec=", &spec, "err=", err)
					if err != nil {
						fmt.Println("L222: FAIL: to parse err=", err, "str=", recStr)
					} else {
						linesChan <- spec
					}
				}
			} else if strings.HasPrefix(aline, "#WAIT") {
				// Add a Pause
				time.Sleep(9500)
				//fmt.Fprintln(u.logFile, "L208: waiting queue=", len(linesChan), "reqPending=", u.reqPending)
				for len(linesChan) > 0 || u.reqPending > 0 {
					u.logFile.Sync()
					fmt.Fprintln(u.logFile, "L209: waiting queue=", len(linesChan), "reqPending=", u.reqPending)
					u.logFile.Sync()
					time.Sleep(5500)
					continue
				}
			} else {
				continue
			}
		} else {
			// Add current line to buffer
			b.Write([]byte(aline))
		}
	}
	close(linesChan)
	u.logFile.Sync()
	jutil.TimeTrackMin(u.logFile, start, "Finished Queing\n")
	jutil.TimeTrackMin(u.logFile, start, "Finished all test records\n")

	<-done // wait until queue has been marked as finished.
	for u.reqPending > 0 {
		time.Sleep(1500)
	}
	u.logFile.Sync()
	u.perf.PrintStat(u.logFile)
	u.logFile.Sync()
	defer u.logFile.Close()
}
