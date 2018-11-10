::REM GOPATH must be set to the base directory that 
::REM contains the src directory for input code.  This
::REM is required for all user code where it will look
::REM for libraries by name as sub directories of source
::REM GOPATH is not capable of listing multiple directories
::REM like java and python can.  
::REM see: https://golang.org/doc/code.html  I don't really
::REM like this approach since I have code for private versus
::REM consulting that has to be kept separate and each project
::REM ends up needing to change GOPATH to point at their base

set GOPATH=%cd%
