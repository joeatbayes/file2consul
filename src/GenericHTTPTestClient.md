# Generic HTTP Test Client

[Generic test client](src/GenericHTTPTestClient.go) provides a data driven, multi threaded test client able to support running at many threads for a while then waiting for all concurrent threads to finish before starting the next test.  This provides basic support for read after write tests.   It also provides easily parsed output that can be used to feed test results into downstream tools.  

The input to Generic Test client is a text file containing a series of JSON strings that describe each test.   It includes a few directives such as #WAIT to indicate a desire to wait for all prior requests to finish.   Comment strings are lines starting with #WAIT

 see [src/GenericHTTPTestClient.md](src/GenericHTTPTestClient.md )  and  [src/GenericHTTPTestClient.go](src/GenericHTTPTestClient.go)



> ## Example Output


    GenericHTTPTestClient.go
    GenericHTTPTestClient.go
    SUCESS: L125: id= 0183 	message= saving JSON record
    L209: waiting queue= 0 reqPending= 1
    SUCESS: L125: id= 0184G 	message= Read after write
    L209: waiting queue= 0 reqPending= 1
    SUCESS: L125: id= 0184M 	message= Read after write existing record
    SUCESS: L125: id= 0185 	message= Check Mismatched 404 miss on known bad key
    FAIL: L128: id=0186 	message=Check unepected response code  	err Msg=	L:97 Expected StatusCode= 519  got= 404
      	verb=GET uri=http://127.0.0.1:9601/mds/test/1817127X5
    FAIL: L128: id=0184F 	message=Check re No Match functionality should fail 	err Msg= L120:FAIL ReNoMatch pattern found in record reNoMatch= .*JOHN MUIR.*  match= true
      	verb=GET uri=http://127.0.0.1:9601/mds/test/1817127X3
    L209: waiting queue= 0 reqPending= 4
    SUCESS: L125: id= 0187 	message= Delete a JSON record
    L209: waiting queue= 0 reqPending= 1
    SUCESS: L125: id= 0188 	message= Read after delete
    SUCESS: L125: id= 0189 	message= Delete a previously deleted record
    Finished Queing
     took 0.000468 min
    Finished all test records
     took 0.000501 min
    numReq= 9 elapSec= 0.0350985 numSuc= 7 numFail= 2 failRate= 0 reqPerSec= 256.42121458181975



## Assumptions

* Unless specified otherwise assumes all requests can be completed in any order which may in fact happen since they are ran in a multi threaded fashion.     



## Command Line API

    GenericHTTPTestClient  -in=data/sample/GenericTestSample.txt -out=test1.log.txt  -MaxThread=100
​    

> > * Runs the test with [GenericTestSample.txt](../data/sample/GenericTestSample.txt) as the input file.    
> > * Writing basic results to test1.log.txt 
> > * Runs with 100 client threads.
> > * **Parameters**
> > * **-in** = the name of the file containing test specifications.
> > * **-out** = the name of the file to write test results and timing to.
> > * **-MaxThread** = maximum number of concurrent requests submitted from client to servers.

## File Input Format

- A series of lines containing JSON text which represents the specifications for the test 

- Each JSON string is terminated by a #END starting a otherwise blank line. 

- **[Sample](../data/sample/GenericTestSample.txt):**

     Test ID = ID to print upon failure
     HTTP Verb = HTTP Verb to send to the server
     URI =  URI to open for this test 
     Headers = Array of Headers to send to the server
     rematch = RE pattern to match the response body against. Not match is failure.
    renomatch = RE pattern that must not be in the response data.

     expected = HTTP Response code expected from the server. 

                other response codes are treated as failure.
     body = Body string to send as Post Body the server 

- Blank rows are ignored

- Rows prefixed by # are treated as comment except when #WAIT or #END

- Rows Prefixed with "#WAIT" Cause the system to pause and wait for all previously queued requests to complete before continuing.  This can allow blocking to allow data setup calls to complete before their read equivalents to complete.

- Test Requests  can be read from file and executed in parallel threads unless blocked by #WAIT directive.

- HTTP VERBS SUPPORTED  GET,PUT,DELETE,POST

- HTTP Headers URI Encoded sent in order specified but this can not be guaranteed since it is treated internally as a map which does not guarantee ordering. 

- POST BODY IS URI Encoded in file but will be decoded prior to POSTING TO Test client.



## Build / Setup

 Download the Metadata server repository

cd  RepositoryBaseDir 

​	eg:  cd \jsoft\mdsgo

​	This is the directory where you downloaded the repository.   

```
go build src/GenericHTTPTestClient.go
```

  or  

```
makeGO.bat
```



## TODO: 

* Add "saveas" parameter to test spec so results of HTTP call are saved to local file as if fetched by Curl.  If present then treated as a relative file name relative to output file.  
* Add ability to run for a while at a given concurrency level and then increase concurrency to find the sweet spot for the server for the current set of data.
* Consider output format that uses JSON to make parsing easy.
* Modify logging output to use atomic output to avoid mixing lines in threads and to reduce sync calls.  Current version could easily mix logs on same line when heavily multi-threaded.
* Add Timing Library for performance logging example.
* Add option to run several test clients simultaneously using os spawn.  This may help ensure client is not blocked in context switching overhead.   May need a way to collate output from multiple spawned processes to get cumulative throughput.

## DONE:

* DONE: JOE:2017-10-20: Add ReNoMatch to ensure a given string is not in the result string from the service. 