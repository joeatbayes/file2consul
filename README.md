# file2consul - Loads config file contents into consul. 

Mirror key, values from property files into consul.       Reduce cost of maintaining larger configuration sets for multiple environments by reducing re-statement and manual editing.  It provides  variable expansion, interpolation, inheritance style overrides and the ability to update multiple consul servers.

> See Also: [Devops with lxd for micro services](https://bitbucket.org/joexdobs/devops_lxd_containers/src/master/documentation/devops-with-lxd.md?fileviewer=file-view-default)  file2consul is a component designed to enable my DevOps vision in.  [Guiding principals for devops](https://bitbucket.org/joexdobs/devops_lxd_containers/src/master/documentation/guiding-principals.md?fileviewer=file-view-default) 

#### Environments are actually quite similar
When building complex software that run on more than one computer we call an environment the set of computers required to run 1 full copy. 

In most companies there is a PROD environment where the production software runs.  A UAT environment where the complete set is tested prior to release into production and TEST where developers can test their modules to be sure they work with other components planned to be released.   In larger companies it is not uncommon to have over a dozen of these environments.

 Building and maintaining the environments can be labor intensive because of all the configuration values that are slightly different between the environments.  This software seeks to make it easier to maintain multiple environments at a lower cost and with less effort. 

In many cases a new environment is actually identical to an existing environment except for a small number of changed parameters. In other instances a configuration value may only changed based on a predictable value such as a environment name such. For example in a TST environment the DB server may be at test-orcle-main.abc.com while in the PROD environment it may be at prod-orcle-main.abc.com. 

#### Minimize Cost & Opportunity for errors 
file2Consul uses interpolation and inheritance to allow a smaller set of configuration parameters to be used whenever values or keys change in predictable ways between environments. 

**Interpolation** allows a single setting to be changed automatically between 
environments.  

**Inheritance**  supports easier derived environments when the differences can not be easily handled by interpolation.     Values are processed from each parent sequentially building a new tree where subsequent derived environments may replace a subset of the configured environments.

## Basic Operation

  Simple example showing building of the Prod settings using a template with mostly variable interpolation.   It uses inheritance override for a few values such as changing the number of network listeners.

> Look at the sample configuration files in [data/config/simple](data/config/simple) they are the best way to learn about how to use file2consule to reduce manual work. 

```sh
file2consul -ENV=PROD -COMPANY=ABC -APPNAME=peopleSearch -IN=data/config/simple/template::data/config/simple/prod -uri=http://127.0.0.1:8500 -CACHE=data/{env}.{appname}.CACHE.b64.txt  -savereadable=data/{env}.lastrun.txt.ini -printlines -verbose 

#Command may be shown as wrapped but it is really one longer command.

# If current directory is not in search path then you
#   may need to modify to change from file2consul to 
#   ./fileconsul on linux or .\file2consul on windows.
```

```
   -IN=name of input parameter file or directory
       If named resource is directory will process all 
	   files in that directory.    Multiple inputs
	   can be specified seprated by :: .  Each 
	   input set
	   
	   will be processed in order with any duplicate 
	   keys in subsequent files overriding those 
	   previously defined for the same key.   This 
	   provides a simple form of inheritance where
	   only the values that change between environment
	   need to be listed while the rest can be inherited
	   from a common parent.  If not specified defaults
	   to data/config/simple/basic
	  

   -URI=uri to reach consul server.   
        If separated by :: (colon) will save to each 
        server listed defaults to http://127.0.0.1:8500
        if not specified.   When set to NONE will 
        run input and compute outputs but will not 
        update consul and will not write to
        local cache file.
	
   -CACHE = name of files to use as cache file.  This
      file is read and compared to the post processing 
      key / value set to determine what values need to
      be saved to consul.  It is also re-written
	  and end of run when defined.  To clear cache
	  delete the file before running the utility.
      This value is interpolated so you can use variables 
      things like enviornment as part of file name.
      will not be written even -URI=NONE.

   -PATHDELIM =  Delimiter to use when splitting list of 
       files, paths or URI for fields like -URI, -IN.   
       Defaults to :: when not set.
       
   -PRINTLINES when this value is specified the
     system will print every input line as it is read
     to help in diagnostics.
     
   -VERBOSE When this value is specified the system 
     will print additional details about values as 
     they are set or re-set during the run. 
   
   -SAVEREADABLE=filename when specified the system will 
      save a human readable version of the keys, values
      after interpolation to the named file. 
   
   -appname = variable used for interpolation
   -env =  variable used for interpolation
   -company = variable used for interpolation
   -appname = varabile used for interpolation

   other named parameters are treated in interpolated values
   Most common error is forgetting - as prefix for named 
   parms
```

> Performance Note:  Remove the -verbose and -printlines flags and the system will run much faster. 
>

* TODO: -runPull optional if present and set to "true" and when the source path is a directory the the system will run a git pull in that directory to fetch most recent copy of the config settings.
* NOTE: Values returned from consul are base64 decoded.  You have to use a Base64 decoder to see what is actually saved in consul.  I was initially confused by this when consul looked like it was returning gibberish.

### A More complex inheritance example

 showing inheritance overrides with values derived from the higher order environments.   This example is one of the more complex where we are actually building an environment configuration for a individual developer but rather than specify everything we specify an order where we process first the base Template then Prod, Then UAT, then DEV and finally joes special properties.  This allows to ensure that we have all the basic settings identical to PROD and then change those as we work down through the other environments.   This helps prevent lower order environments from accidentally being different than prod.   While this does help with consistency we used the {env} in all the key names so even though some of the basic configuration came from prod we can be sure we do not accidentally change prod config values.    We also use {env} and other variables in several of the  data bodies to give vips that have similar but predictably different names differentiation so they do not conflict with prod assets.

```sh
file2consul -ENV=JTEST -COMPANY=ACME -APPNAME=peopleSearch -IN=data/config/simple/template::data/config/simple/prod::data/config/simple/uat::data/config/simple/dev::data/config/simple/joes-dev.prop.txt -uri=http://127.0.0.1:8500

#Command may be shown as wrapped but it is really one longer command.   

#NOTE: Absence of cache command will cause all of consul values to be udpated every time. 
```

### Simple Operation of Dumb version without Inheritance

Please note the dumb version does not attempt detect lines which have not changed.  As a result it sends every config setting to consul every time it is ran.  The full version of file2consul keeps the last values saved and only updates consul when something actually changed. 

```sh
go build src/file2consul-dumb.go

file2consul-dumb -ENV=DEMO -COMPANY=ABC -APPNAME=file2consul-dumb -FILE=data/config/simple/template/basic.prop.txt -uri=http://127.0.0.1:8500 -CACHE=data/{env}.CACHE.b64.txt
 
  # the file2Consul-dump lines wrapped for display when entering it should be one long line. 
 
 -file=name of input paramter file
 
 -uri=uri to reach console server
    
 -appname = variable used for interpolation
 -env =  variable used for interpolation
 -company = variable used for interpolation
 -appname = varabile used for interpolation
 
  other variables can be defined as needed
  variables are not case sensitive.
```

## Basic Interpolation
Basic interpolation allows interpolation of defied environment variables into existing values.  This can allow a single config 
string to be used across multiple environments without requiring separate files.

     {ENV}/DBServer=orcl.master.{ENV}
     
     Assuming the a variable ENV has been defined as UAT01
     the key becomes UAT01.DBServer while the value became 
     orcl.master.UAT01.  Interpolation can occur on either 
     the key and or the data values or both.
     
     Interpolation variables can be defined on the command line using the -varaible notation  There are a few predefined named parameters the program uses for it's own operation such as -uri which indicates the set of URI it should use to talk to consul.   Even these pre-defined variables are available for use in interpolation. 

## Using Ancestor Overrides

Ancestor overrides allow for changes that are specific to an environment that can not be easily handled with interpolation changes.

Sample Usage

TODO:  Add more detail here

## Concatenating multiple Lines

```
{env}/register/vip/members = S1811717127.{env}.{company}.com
  + S1811717128.{env}.{company}.com
  + S1811717129.{env}.{company}.com
  + S1811717130.{env}.{company}.com
```

Starting a line with + indicates you want to add the value on that line into the last key.  This makes it easier to break long values into smaller lines while editing and still save them in a single consul Value for that key.    The values are delimited with \t to make it easier 

## Including External Files

```ini
{ENV}/payments/vip/members=@../payments.vip.members.ini
  
  # When the first character of a value starts 
  # with @ the system will attempt to find a file
  # a file of that name relative to the file where
  # it was mentioned.   It will include that file
  # as the value for that key.   Interpolation is done
  # on both the file name and the file contents if it
  # is found.   If the file is not found then then
  # unmodified string will be kept as it's value. 
  # and a warning message will be printed.  This was
  # designed to make it easier to include larger lists 
  # of values like members in a VIP rather than using
  # the + concatenation semantic above.
```

## Supplying the ACL Token

Will accept token via environment variable `CONSUL_HTTP_TOKEN`

## Build & Setup

This software has been tested to build on Windows10,  Ubuntu,  MacOS. It should run fine on any computer where GO is available.  GO is only needed at build time, You can distribute the executable file without GO present.   

#### Prebuilt Binaries Download

* [Prebuilt binaries from  file2consul wiki](https://raw.githubusercontent.com/wiki/joeatbayes/file2consul/builds/most-recent/README.md)

  We moved the binaries out of the main repository because they were consuming a lot of space due to frequent updates.   We will now place builds when strategic features are completed in the pre-built binaries folder.   You can also build your own as described below.

#### Build Your Own

* Install GO compiler which can be downloaded from https://golang.org/dl/ 

* Download the repository  using GIT  from command line

```sh
git clone https://github.com/joeatbayes/file2consul.git
```

* To Download using your browser, open the following URI  https://github.com/joeatbayes/file2consul/archive/master.zip   Once downloaded save in your desired directory and unzip.    You will need to be able to open a shell at that directory

* Add the directory where you copied the source to your [search PATH](http://www.linfo.org/path_env_var.html)  This can be done temporarily by running the [setGOEnv.bat](setGOEnv.bat) on windows or [setGOEnv.sh](setGOEnv.sh) on Linux or mac.  These are included in the downloaded repository.
* Build the software by running  [makeGO.sh][makeGO.sh] on Linux or running [makeGO.bat](makeGO.bat) on windows.   It should produce several executable files including file2Consul.exe on windows or a executable file2Consul on Linux.
* The file2Consule [executable](http://www.linfo.org/executable.html) can be copied to any location in the [search PATH](http://www.linfo.org/path_env_var.html).  It will always look for files relative to the [current working directory](http://www.informit.com/articles/article.aspx?p=2468330&seqNum=15) unless the paths specified on the command line are [absolute paths](https://www.linux.com/blog/absolute-path-vs-relative-path-linuxunix).   We generally leave our in the same directory where we downloaded the repository to make it easy access our sample input files.

## Main Files



* [makeGo.sh linux](makeGO.sh)  [makeGO.bat Windows](makeGO.bat)  Builds the main executables from the GO Lang source code
* [License](LICENSE.md)
* TODO: Fill This in

## Participating

[Contact me on linkedin](https://linkedin.com/in/joe-ellsworth-68222/) to obtain enhancements or updates.

Please donate or send us a link to your project when you are using file2consul. We love to publish success stories.

You can contribute features by making a fork and submitting a pull request.

You can request new features or a Bug Fix by filing a new Issue on Bitbucket.

I give first priority to features and bug fixes from people willing to pay my hourly rate for the work. 

## Reference

* [git2consul](https://github.com/breser/git2consul) a similar utility but it does not support interpolation or slot based inheritance.

* ##### Consul

  * [Consul download Page](https://www.consul.io/downloads.html)
  * [Consul http API guide](https://www.consul.io/api/index.html),
  * [Consul   kv api](https://www.consul.io/api/kv.html)
  * [Consul ACL security guide](https://www.consul.io/docs/guides/acl.html) 
  * [consul kv by example](https://github.com/JoergM/consul-examples/tree/master/http_api)
  * [Getting started with consul KV](https://learn.hashicorp.com/consul/getting-started/kv.html)  example with command line instead of http
  * [Blog first look of kv interface in consul](http://alesnosek.com/blog/2017/07/15/first-look-at-the-key-value-store-in-consul/)
  * [Tutorials point consul introduction](https://www.tutorialspoint.com/consul/consul_quick_guide.htm)
  * [Consul cheat sheet](https://lzone.de/cheat-sheet/Consul)
  * [consul python api](http://www.voidcn.com/article/p-hjdesarg-up.html)
