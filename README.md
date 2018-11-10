# file2consul - Loads config file contents into consul. 

Consul provides a fairly simple KV configuration management system.  This works but the flat model can make supporting multiple environments that may contain hundreds of discrete configuration values onerous.  
#### Environments are actually quite similar
In many companies when building complex software that runs on more than one computer we call an environment the set of computers required to run 1 copy. 

In most companies there is a PROD environment where the production software runs.  A UAT environment where the complete set is tested prior to release into production and TEST where developers can test their modules to be sure they work with other components planned to be released.   In larger companies it is not uncommon to have over a dozen of these enviornments.

 Building and maintaining the environments can be labor intensive because of all the configuration values that are slightly different between the environments.  This software seeks to make it easier to maintain multiple environments at a lower cost and with less effort. 

In many cases a new environment is actually identical to an existing environment except for a small number of changed parameters. In other instances a configuration value may only changed based on a predictable value such as a environment name such. For example in a TST environment the DB server may be at test-orcle-main.abc.com while in the PROD environment it may be at prod-orcle-main.abc.com. 

#### Minimize Cost & Opportunity for errors 
file2consul seeks to minimize the total opportunity for manually introduced errors and the amount of work in managing configuration parameters by optimizing for use of these two patterns.  

file2Consul uses interpolation and inheritance to allow a smaller set of configuration parameters to be used whenever values or keys change in predictable ways between environments. 

Interpolation allows a single setting to be changed automatically between 
environments.  A slot based inheritance model to support easier derived environments.   It will process the values in each parent in the tree sequentially building a new tree where subsequent derived environments may replace a subset of the configured environments.


## Basic Operation

    file2consul -E ENV=UAT01  -S=../prodConfig,..UAT/config,..TSTConfig,..JOESpecialConfig -sv127.0.0.0.1:8500 -cache=./cache/cache.dta
  -e Sets a local environment variable.    The system also reads environment variables currently defined for the current session.
​    
  -p A list of paths that resolve to either a directory or a file.     The system will read and process the contents of each in the order  specified.  If a value for a given key is defined multiple times  the one encountered on the last file will win.
​     
  -s Names the consul server to the values to.   If multiple servers are listed then the consul values will be copied to each of these servers.

  -cache=optional - Name of a file the system will save the full set of key, values after special processing  to.   This file is used to compare key/values from the last run so it can send only values  that have changed to consul. 

 -runPull optional if present and when the source path is a directory the the system will run a git pull in that directory to fetch most recent copy of the config settings.



### Simple Operation of Dumb version without Inheritance

```sh
go build src/file2consul-dumb.go

file2consul-dumb -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -FILE=data/simple-config/basic.prop.txt -uri=http://127.0.0.1:8500
 -file=name of input paramter file
 -uri=uri to reach console server
 -appname = variable used for interpolation
 -env =  variable used for interpolation
 -company = variable used for interpolation
 -appname = varabile used for interpolation
  other varaibles can be defined as needed
```





## Basic Interpolation
Basic interpolation allows interpolation of defied environment variables into existing values.  This can allow a single config 
string to be used across multiple environments without requiring separate files.

     $^ENV.DBServer=orcl.master.$^ENV
     
     Assuming the an environment variable ENV has been defined as UAT01
     the key becomes UAT01.DBServer while the value became orcl.master.UAT01.  Interpolation can occur on either the key
     and or the data values.

## Using Ancestor Overrides

Ancestor overrides allow for changes that are specific to an environment that can not be easily handled with interpolation changes.

Sample Usage




## Build & Setup

Install GO compiler which can be downloaded from https://golang.org/dl/ 

* Download the repository  using GIT  from command line

```sh
git clone https://github.com/joeatbayes/file2consul.git
```

* Or Download using HTTP   Open the following URI in your browser https://github.com/joeatbayes/file2consul/archive/master.zip   Once downloaded save in your desired directory and unzip.    You will need to be able to open a shell at that directory

* Add the directory where you copied the source to your [search PATH](http://www.linfo.org/path_env_var.html)  This can be done temporarily by running the [setGOEnv.bat](setGOEnv.bat) on windows or [setGOEnv.sh](setGOEnv.sh) on Linux or mac.  These are included in the downloaded repository.
* Build the software by running  [makeGO.sh][makeGO.sh] on Linux or running [makeGO.bat](makeGO.bat) on windows.   It should produce several executable files including file2Consul.exe on windows or a executable file2Consul on Linux.
* The file2Consule [executable](http://www.linfo.org/executable.html) can be copied to any location in the [search PATH](http://www.linfo.org/path_env_var.html).  It will always look for files relative to the [current working directory](http://www.informit.com/articles/article.aspx?p=2468330&seqNum=15) unless the paths specified on the command line are [absolute paths](https://www.linux.com/blog/absolute-path-vs-relative-path-linuxunix).   We generally leave our in the same directory where we downloaded the repository to make it easy access our sample input files.
* 







## Main Files



* [makeGo.sh linux](makeGO.sh)  [makeGO.bat Windows](makeGO.bat)  Builds the main executables from the GO Lang source code
* [License](LICENSE.md)

## Reference

* [git2consul](https://github.com/breser/git2consul) a similar utility but it does not support interpolation or slot based inheritance.

* [Consul Emulator in ESP32 Hardware](....) Runs the HTTP listener portion of a Consul network listener on $10 worth of hardware.    Why run Consul on VM's that will cost thousands of dollars per year simply run 3 or 4 on these 10 dollar modules put them behind a vip and you have the cheapest high availability configuration store available.

* ##### Consul

  * [Consul download Page](https://www.consul.io/downloads.html)
  * [Consul http API guide](https://www.consul.io/api/index.html) 
  * [Consul ACL security guide](https://www.consul.io/docs/guides/acl.html) 
