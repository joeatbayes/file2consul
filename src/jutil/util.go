package jutil

import (
	"fmt"
	//"net/http"
	"os"
	"strings"
	"time"
)

func TimeTrack(fi *os.File, start time.Time, name string) {
	elapsed := time.Since(start).Seconds() * 1000.0
	fmt.Fprintf(fi, "%s took %fms\n", name, elapsed)
	fi.Sync()
}

func TimeTrackMin(fi *os.File, start time.Time, name string) {
	elapsed := time.Since(start).Seconds() / 60.0
	fi.Sync()
	fmt.Fprintf(fi, "%s took %f min\n", name, elapsed)
	fi.Sync()
}

func Nowms() float64 {
	nn := time.Now()
	return float64(nn.UnixNano()) / float64(time.Millisecond)
}

func Elap(msg string, beg_time float64, end_time float64) float64 {
	elap := end_time - beg_time
	fmt.Printf("ELAP %s = %12.3f ms\n", msg, elap)
	return elap
}

func Check(msg string, e error) {
	if e != nil {
		fmt.Println("ERROR:")
		fmt.Println(e)
		panic(e)
	}
}

func KeepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}
