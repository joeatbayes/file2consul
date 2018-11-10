package jutil

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type PerfMeasure struct {
	NumReq            int
	NumSuc            int
	NumFail           int
	StartTime         time.Time
	NumSinceStatPrint int
	PrintEvery        int
	BytesRead         int64
	ByteWritten       int64
}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func (u *PerfMeasure) PrintStat(fi *os.File) {
	elapSec := time.Since(u.StartTime).Seconds()
	reqPerSec := float64(u.NumReq) / float64(elapSec)
	failRate := u.NumFail / u.NumReq
	fmt.Fprintln(fi, "numReq=", u.NumReq, "elapSec=", elapSec, "numSuc=", u.NumSuc, "numFail=", u.NumFail, "failRate=", failRate, "reqPerSec=", reqPerSec)
	fi.Sync()
	u.NumSinceStatPrint = 0
}

func (u *PerfMeasure) CheckAndPrintStat(fi *os.File) {
	if u.NumSinceStatPrint >= u.PrintEvery {
		u.PrintStat(fi)
	}
}

func MakePerfMeasure(printEvery int) *PerfMeasure {
	return &PerfMeasure{StartTime: time.Now().UTC(), PrintEvery: printEvery}
}
