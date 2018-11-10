package jutil

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func TimeTrack(fi *os.File, start time.Time, name string) {
	elapsed := time.Since(start).Seconds() * 1000.0
	fmt.Fprintf(fi, "%s took %fms", name, elapsed)
	fi.Sync()
}

func TimeTrackMin(fi *os.File, start time.Time, name string) {
	elapsed := time.Since(start).Seconds() / 60.0
	fi.Sync()
	fmt.Fprintf(fi, "%s took %f min\n", name, elapsed)
	fi.Sync()
}

func KeepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}
