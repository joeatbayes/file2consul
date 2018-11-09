# file2consul - Load configuration file contents to consul config store. 

Consul provides a fairly simple KV configuration managment system.  This works but it the flat model can make 
supporting multiple enviornments that may contain hundreds of discrete config values onerous.  In many cases a 
new enviornment it actually identical to an existing environment except for a small number of changed parameters.

file2Consul uses a slot based inheritance model to support easier derived enironments.   It will process the values in 
one each parent i the tree sequentially building a new tree where subsequet derived enviornments may replace a subset of
the configured enviornments.

It also supports basic interpolation of enviornment based strings into the Key and values to allow a single set of 
configurtion parameters

## Basic Operation

    file2consul -E ENV=UAT01  -S=../prodConfig,..UAT/config,..TSTConfig,..JOESpecialConfig -sv127.0.0.0.1:8500 -cache=./cache/cache.dta
  -e Sets a local enviornment variable.  
    The system also reads enviornment variables currently defined for the current session.
    
  -p A list of paths that resolve to either a directory or a file.  
     the system will read and process the contents of each in the order
     specified.  If a value for a given key is defined multiple times
     the one encounted on the last file will win.
     
  -s Names the consul server to the values to.
  
  -cache=optional - Name of a file the system will save the full set of key, values after special processing
       to.   This file is used to compare key/values from the last run so it can send only values
       that have changed to consul. 
 

## Basic Interpolation
Basic interpolation allows interpolation of defied environment variables into existing values.  This can allow a single config 
string to be used across multiple enviornments without requiring separate files.

     $^ENV.DBServer=orcl.master.$^ENV
     Assuming the an enviornment variable ENV has been defined as UAT01
     the key becomes UAT01.DBServer while the value became orcle.master.UAT01

## Using the Anceseter Overrides
