package jutil

// tested with go 1.11.2 some functions not availble in older versions.

import (
	"bytes"
	"fmt"

	//"net/http"
	"strings"
	//"io"
	//"io/ioutil"
	//"log"
	//"net/url"
	"os"
	//"os/exec"
)

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func MakeDirName(baseDir string, id string) string {
	// newpath := filepath.Join(".", "public")
	//return dirName
	return "data/xyz"
}

func EnsurDir(dpath string) {

	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		os.MkdirAll(dpath, os.ModePerm)
	}

}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// Map Bad characters we do not want in
// our file names into their ASII number
// equivelants.
func MapBadIDChar(strIn string) string {
	var buff bytes.Buffer
	slower := strings.ToLower(strIn)
	for _, d := range slower {
		// if character outside range of 0..9, a..z, -_ then force conversion
		// to ascii number value.
		if (d >= '0' && d <= '9') || (d >= 'a' && d <= 'z') || d == '.' || d == '-' || d == '_' {
			buff.WriteByte(byte(d))
		} else {
			buff.WriteString(fmt.Sprintf("%%%X", d))
		}
	}
	return buff.String()
}
