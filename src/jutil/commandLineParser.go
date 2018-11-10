/* CommandLineParser.go

Parse command line parameters into easy to use format.

*/
package jutil

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	s "strings"
)

type ParsedCommandArgs struct {
	ExeName   string
	PositArgs []string
	NamedStr  map[string]string
	NamedInt  map[string]int
}

// Parse command line parameters
// those as simple parms will show
// up in the PositArgs array in the
// they appear in the command line
// -name
//   pattern of name without value
//   will be placed in NamedInt[key]
//   with values of 0.  Where key is
//   the value supplied in name with
//   spaces removed and converted to
//   lower case.
//   eg:  -Test  will appear in NamedInt
//   as  NamedInt["test"] = 0
//
// -name
//    will show up in NamedInt[key] with
//    value of 0
//    eg: -TEST will shoup up as
//    NamedInt["test"] = 0
//
// -name=false
//   will show up in NamedInt[key] with
//   value of 0.  False is case insensitive.
//   eg: -Test=FALSE will show up as
//   NamedInt["test"] = 0.  No spaces are
//   allowed around the equal
//
// -name=True
//   will show up in NamedInt[key] with value
//   of 1.   True is case insensitive.
//   eg: -TEST=tRue  will show up as
//   NamedInt["test"]= 1
//
// -name=jimbo
//   where jimbo is a string that can not be
//   safely parsed to a integer value will show
//   up in  NameeStr[key] with value on right side
//   of equal as value.
//   eg:  -Name=C:\JoeApple\93.csv" will show up
//   in NamedStr["name"] = "C:\JoeApple\93.csv"
//   No spaces are allowed surrounding name or
//   in the value parameter.
//
// Sample:
//   ex-command-args inputFile -class=classify.out.txt -numBuck=10 -optimize=03 -X5 Apple")
//
func ParseCommandLine(args []string) *ParsedCommandArgs {
	tout := new(ParsedCommandArgs)
	tout.PositArgs = make([]string, 1)
	tout.NamedStr = make(map[string]string)
	tout.NamedInt = make(map[string]int)
	tout.ExeName = args[0]
	tout.PositArgs[0] = tout.ExeName
	numArgs := len(args)
	for ndx := 1; ndx < numArgs; ndx++ {
		arg := s.TrimSpace(args[ndx])
		if s.HasPrefix(arg, "-") {
			a := s.SplitN(arg, "=", 2)
			key := s.TrimPrefix(s.ToLower(a[0]), "-")
			if len(a) == 1 {
				tout.NamedInt[key] = 0
			} else {
				ptxt := s.TrimSpace(a[1])
				lctxt := s.ToLower(ptxt)
				//fmt.Printf("key=%v ptxt=%v\n", key, ptxt)
				if lctxt == "true" {
					tout.NamedInt[key] = 1
				} else if lctxt == "false" {
					tout.NamedInt[key] = 0
				} else if len(ptxt) < 1 {
					tout.NamedInt[key] = 0
				} else {
					i64, err := strconv.ParseInt(ptxt, 10, 32)
					//fmt.Printf("key=%v ptxt=%v i64=%v err=%v\n", key, ptxt, i64, err)
					if err == nil {
						ival := int(i64)
						tout.NamedInt[key] = ival
						tout.NamedStr[key] = ptxt
					} else {
						tout.NamedStr[key] = ptxt
					}
				}
			}
		} else {
			tout.PositArgs = append(tout.PositArgs, arg)
		}
	} // for
	return tout
}

func (parg *ParsedCommandArgs) String() string {
	var sbb bytes.Buffer
	sb := &sbb
	fmt.Fprintf(sb, "  EXEName=%v\n  Positional\n", parg.ExeName)

	for ndx, val := range parg.PositArgs {
		fmt.Fprintf(sb, "    ndx=%v\tval=%v\n", ndx, val)
	}

	fmt.Fprintf(sb, "  NamedStrArg\n")
	for key, val := range parg.NamedStr {
		fmt.Fprintf(sb, "    key=%v\t val=%v\n", key, val)
	}

	fmt.Fprintf(sb, "  NamedIntArg\n")
	for key, val := range parg.NamedInt {
		fmt.Fprintf(sb, "    key=%v\t val=%v\n", key, val)
	}

	return sb.String()

}

// Return true is the name was defined in the command line
// otherwise false.  It can be either a int or string parm
// and it will return true.
func (parg *ParsedCommandArgs) Exists(name string) bool {
	_, found := parg.NamedInt[name]
	if found {
		return true
	}

	_, found2 := parg.NamedStr[name]
	if found2 {
		return true
	}
	return false
}

// Return the named Integer value or the specified default
// if not found.
func (parg *ParsedCommandArgs) Ival(name string, defaultVal int) int {
	val, found := parg.NamedInt[name]
	if found {
		return val
	} else {
		return defaultVal
	}
}

// Return the named float 32 value or the specified default
// if not found.
func (parg *ParsedCommandArgs) Fval(name string, defaultVal float32) float32 {
	val, found := parg.NamedStr[name]
	if found {
		val = s.TrimSpace(val)
		f64, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return defaultVal
		} else {
			return float32(f64)
		}
	} else {
		return defaultVal
	}
}

// Return the named float 54 value or the specified default
// if not found.
func (parg *ParsedCommandArgs) F64val(name string, defaultVal float64) float64 {
	val, found := parg.NamedStr[name]
	if found {
		val = s.TrimSpace(val)
		f64, err := strconv.ParseFloat(val, 64)
		//fmt.Printf("\n\n\nF64val ptxt=%s  f64=%v err=%v\n", f64, val, err)
		if err != nil {
			return defaultVal
		} else {
			return f64
		}
	} else {
		return defaultVal
	}
}

// return string parameter value if it was defined or
// the default value if not.
func (parg *ParsedCommandArgs) Sval(name string, defaultVal string) string {
	val, found := parg.NamedStr[name]
	if found {
		return val
	} else {
		return defaultVal
	}
}

// return boolean equivelant of parameter value specified or
// default if not specified.
func (parg *ParsedCommandArgs) Bval(name string, defaultVal bool) bool {
	val, found := parg.NamedInt[name]
	//fmt.Printf("Bval name=%s val=%v  found=%v", name, val, found)
	if found {
		return val == 1
	} else {
		return defaultVal
	}
}

// Parse a parameter as a list of strings where
// the integer payload is the position within the
// specified list where the string was found
// used to parse a comma delimited set of strings
// such as a list of column names. If name is not
// found then returns a empty dictionary
func (parg *ParsedCommandArgs) SListDict(aname string) map[string]int {
	tout := make(map[string]int)
	val, found := parg.NamedStr[aname]
	if found {
		sarr := s.Split(val, ",")
		for ndx, astr := range sarr {
			astr = s.TrimSpace(astr)
			tout[astr] = ndx
		}
	}
	return tout
}

var str = "Old Value!"
var ParmMatch, ParmErr = regexp.Compile("\\{.*?\\}")

func (parg *ParsedCommandArgs) Interpolate(str string) string {
	ms := ParmMatch.FindAllIndex([]byte(str), -1)
	if len(ms) < 1 {
		return str // no match found
	}
	//sb := strings.Builder
	var sb []string
	last := 0
	slen := len(str)
	for _, m := range ms {
		start, end := m[0]+1, m[1]-1
		//fmt.Printf("m[0]=%d m[1]=%d match = %q\n", m[0], m[1], str[start:end])
		if start > last-1 {
			// add the string before the match to the buffer
			sb = append(sb, str[last:start-1])
		}
		aMatchStr := strings.ToLower(str[start:end])
		// substitute match string with parms value
		// or add it back in with the {} protecting it
		// TODO: Add lookup from enviornment variable
		//  if do not find it in the command line parms
		lookVal := parg.Sval(strings.ToLower(aMatchStr), "{"+aMatchStr+"}")
		//fmt.Printf("matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
		sb = append(sb, lookVal)
		last = end + 1
	}
	if last < slen-1 {
		// append any remaining characters after
		// end of the last match
		sb = append(sb, str[last:slen])
	}
	return strings.Join(sb, "")
}

func CommandLineParserTest() {
	args := os.Args
	fmt.Println("CommandLineParserTest()\n")
	fmt.Println("\texample:  executableName inputFile -class=classify.out.txt -numBuck=10 -optimize=03 -X5 Apple")
	fmt.Println("\t  Note: pattern name=value with no spaces presents as a single parameter")
	fmt.Println("\t  Note: pattern name space value presents as two separate paramters")
	fmt.Println("\t  order of named parameters is not important.  Order of postional may be important")
	fmt.Println("\t  Most common error is forgetting - as prefix for named parms")
	fmt.Printf("\t   Number of Args: %v  args: %s\n", len(args), args)

	pargs := ParseCommandLine(args)

	// see of a simple parm like -X5 was
	// defined.
	if pargs.Exists("x5") {
		fmt.Println("Found X5 Setting")
	}

	// normal way to access value in a dictionary
	// if paramater such as -NumBuck=10 was specified.
	numBuck, nbFound := pargs.NamedInt["numbuck"]
	if nbFound {
		fmt.Printf("numBuck=%v\n", numBuck)
	}

	// Shorthand way of accesing Integer parms
	// such as -NumBuck=99
	numBuck2 := pargs.Ival("numbuck", 11)
	fmt.Printf("easy way numBuck2=%v\n", numBuck2)

	// Shorthand way of accesing Integer parms
	// with a default for paramter most likely
	// not supplied.
	fnum := pargs.Ival("favorite_num", -9181)
	fmt.Printf("favorite_num=%v\n", fnum)

	// Shorthand way of accessing str Parms with default
	// such as  -class=./data/2015-01-01.csv
	classFile := pargs.Sval("class", "")
	fmt.Printf("class =%s\n", classFile)

}
