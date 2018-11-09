# file2consul - Loads config file contents into consul. 

Consul provides a fairly simple KV configuration management system.  This works but the flat model can make supporting multiple environments that may contain hundreds of discrete configuration values onerous.  
#### Environments are actually quite similar
In many cases a new environment is actually identical to an existing environment except for a small number of changed parameters. In other instances a configuration value may only changed based on a predictable value such as a environment name such. For example in a TST environment the DB server may be at test-orcle-main.abc.com while in the PROD environment it may be at prod-orcle-main.abc.com. 
#### Minimize Cost & Oportunity for errors 
file2consul seeks to minimize the total opportunity for manually introduced errors and the amount of work in managing configuration parameters by optimizing for use of these two patterns.  

file2Consul uses interpolation and inheritance to allow a smaller set of configuration parameters to be used whenever values or keys change in predictable ways between environments. 

Interpolation allows a single setting 
to be changed automatically between 
environments.  A slot based inheritance model to
support easier derived environments.   It will
process the values in each parent in the tree 
sequentially building a new tree where subsequent 
derived environments may replace a subset of
the configured environments.


## Basic Operation

    file2consul -E ENV=UAT01  -S=../prodConfig,..UAT/config,..TSTConfig,..JOESpecialConfig -sv127.0.0.0.1:8500 -cache=./cache/cache.dta
  -e Sets a local environment variable.  
    The system also reads environment variables currently defined for the current session.
    
  -p A list of paths that resolve to either a directory or a file.  
     the system will read and process the contents of each in the order
     specified.  If a value for a given key is defined multiple times
     the one encountered on the last file will win.
     
  -s Names the consul server to the values to.
  
  -cache=optional - Name of a file the system will save the full set of key, values after special processing
       to.   This file is used to compare key/values from the last run so it can send only values
       that have changed to consul. 
 
 -runPull optional if present and when the source path is a directory the the system will run a git pull in that directory to fetch most recent copy of the config settings.

## Basic Interpolation
Basic interpolation allows interpolation of defied environment variables into existing values.  This can allow a single config 
string to be used across multiple environments without requiring separate files.

     $^ENV.DBServer=orcl.master.$^ENV
     
     Assuming the an environment variable ENV has been defined as UAT01
     the key becomes UAT01.DBServer while the value became orcl.master.UAT01.  Interpolation can occur on either the key
     and or the data values.

## Using the Ancestor Overrides


## Build & Setup

## Reference

* [git2consul](https://github.com/breser/git2consul) a similar utility but it does not support interpolation slot based inheritance.
* [Consul http API guide](https://www.consul.io/api/index.html) 
* [Consul ACL security guide](https://www.consul.io/docs/guides/acl.html) 
