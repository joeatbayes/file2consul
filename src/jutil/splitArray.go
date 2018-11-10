/* splitArray.go  Utility to split an array into two parts.  One used for
testing the other used for training.
(C) Joe Ellsworth Jan-2017 License MIT See license.txt */
package jutil

// Divide the input array into two output arrays
// main which contains most of the records and
// auxOut which receives one row for every
// oneEvery rows in main.  The first skipNumFirst
// rows always go to main before aux gets any.
// This is normally used to isolate a set of data
// between training and test but is also used to
// isolate rows for the optimizer
// skipNumFirst is used to shuffle sets to allow
// easy extraction of different records to vary
// the test set.
func SplitFloatArrOneEvery(ain [][]string, skipNumFirst int, oneEvery int) ([][]string, [][]string) {
	numRow := len(ain)
	auxNumRow := (numRow - skipNumFirst) / oneEvery
	if auxNumRow < 1 {
		auxNumRow = 1
	}
	mainNumRow := numRow - auxNumRow
	if mainNumRow < 1 {
		mainNumRow = 1
	}

	mOut := make([][]string, 0, mainNumRow+2)
	auxOut := make([][]string, 0, auxNumRow+2)
	keepCnt := 0
	for ndx, row := range ain {
		if ndx >= skipNumFirst && keepCnt >= oneEvery {
			auxOut = append(auxOut, row)
			keepCnt = 0
		} else {
			mOut = append(mOut, row)
			keepCnt += 1
		}
	} // for
	return mOut, auxOut
}

// Split the input array into two sets of rows.  Where the
// second array is sized to be len * testPortion long.
// Note: This is much faster than SplitFloatArrOneEvery because
// it only has to allocate new pointers with no underlying
// copy of the array elements
func SplitFloatArrTail(ain [][]string, testPortion float32) ([][]string, [][]string) {
	arrLen := len(ain)
	arr2Len := int(float32(arrLen) * testPortion)
	arr1EndNdx := arrLen - arr2Len
	arr1 := ain[:arr1EndNdx]
	arr2 := ain[arr1EndNdx+1:]
	return arr1, arr2
}
