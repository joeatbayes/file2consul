package jutil

// tested with go 1.11.2 some functions not availble in older versions.

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"

	//"net/http"
	"strings"
	//"io"
	//"io/ioutil"
	//"log"
	"net/url"
	"os"

	//"os/exec"
	"encoding/base64"
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

func SaveStrsToFile(inArr []string, fiName string, b64Val bool, pargs *ParsedCommandArgs) {
	start := Nowms()
	verboseFlg := pargs.Exists("verbose")
	if verboseFlg {
		fmt.Println("SaveDictToFile ", fiName, " b64Val=", b64Val)
	}
	f, err := os.Create(fiName)
	if err != nil {
		fmt.Println("ERROR: opening file for write fiName=", fiName, " err=", err)
		return
	}

	for ndx, val := range inArr {
		key := strconv.FormatInt(int64(ndx), 10)
		var saveStr string
		if b64Val {
			saveStr = key + "=" + val + "\n"
		} else {
			val = strings.Replace(val, "\n", "\t", -1) // must replace embeded CR or will mess up readability.
			val = strings.Replace(val, "\r", " ", -1)  // in saved dictionary
			saveStr = key + "=" + val + "\n"
		}
		if verboseFlg {
			fmt.Println("saveStr=", saveStr)
		}
		_, err := f.WriteString(saveStr)
		if err != nil {
			fmt.Println("ERROR: writing to ", fiName, " err=", err)
			defer f.Close()
			return
		}
	}
	f.Sync()
	f.Close()
	Elap("saveDictFile "+fiName, start, Nowms())
}

/* save dictionary to file.  Use base64 encoding for values
because some values may contain vertical whitespace When b64Val
is true then encode the values otherwise write them native.*/
func SaveDictToFile(sdict map[string]string, fiName string, b64Flg bool, pargs *ParsedCommandArgs) {
	// TODO: May be faster to endode this in a in memory
	// butter and write as a block.
	//b64Flg := pargs.Exists("b64")
	URLEncode := pargs.Exists("urlencode")

	start := Nowms()
	verboseFlg := pargs.Exists("verbose")
	if verboseFlg {
		fmt.Println("SaveDictToFile ", fiName, " b64Flg=", b64Flg)
	}
	//fmt.Println("SaveDictToFile ", fiName, " b64Flg=", b64Flg, " urlencode=", URLEncode)
	f, err := os.Create(fiName)
	if err != nil {
		fmt.Println("ERROR: opening dict file for write fiName=", fiName, " err=", err)
		return
	}

	for key, val := range sdict {
		var saveStr string

		if URLEncode {
			saveStr = key + "=" + url.QueryEscape(val)
		} else if b64Flg {
			saveStr = key + "=" + base64.StdEncoding.EncodeToString([]byte(val))
		} else {
			val = strings.Replace(val, "\n", "\t", -1) // must replace embeded CR or will mess up readability.
			val = strings.Replace(val, "\r", " ", -1)  // in saved dictionary
			saveStr = key + "=" + val
		}
		if verboseFlg {
			fmt.Println("saveStr=", saveStr)
		}
		_, err := f.WriteString(saveStr)
		if err != nil {
			fmt.Println("ERROR: writing to ", fiName, " err=", err)
			defer f.Close()
			return
		}
		_, err = f.WriteString("\n")
	}
	f.Sync()
	f.Close()
	Elap("saveDictFile "+fiName, start, Nowms())
}

func LoadDictFile(inFiName string, b64Decode bool, pargs *ParsedCommandArgs) map[string]string {
	start := Nowms()
	verboseFlg := pargs.Exists("verbose")
	if verboseFlg {
		fmt.Println("LoadDictFile ", inFiName, "b64Decode=", b64Decode)
	}
	tout := make(map[string]string)
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("ERROR: loadDictFile: error opening in file ", inFiName, " err=", err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	lineCnt := 0
	for scanner.Scan() {
		aline := scanner.Text()
		lineCnt++
		if len(aline) <= 0 {
			continue
		}
		aline = strings.TrimSpace(aline)
		if strings.HasPrefix(aline, "#") {
			continue
		}

		arr := strings.SplitN(aline, "=", 2)
		if len(arr) != 2 {
			fmt.Println("line#", lineCnt, "fails split on = test", " line=", aline)
			continue
		}
		aKey := strings.TrimSpace(arr[0])
		aVal := arr[1]
		if b64Decode {
			uDec, _ := base64.StdEncoding.DecodeString(aVal)
			aVal = string(uDec)
		}
		tout[aKey] = aVal
		if verboseFlg {
			fmt.Println("LoadDictFile: key=", aKey, " val=", aVal)
		}
	}
	Elap("loadDictFile "+inFiName, start, Nowms())
	return tout
}
