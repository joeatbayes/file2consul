#!/bin/sh
# GOPATH must be export to the base directory that 
# contains the src directory for input code.  This
# is required for all user code where it will look
# for libraries by name as sub directories of source
# GOPATH is not capable of listing multiple directories
# like java and python can.  
# see: https://golang.org/doc/code.html  I don't really
# like this approach since I have code for private versus
# consulting that has to be kept separate and each project
# ends up needing to change GOPATH to point at their base

export GOPATH=`pwd`
export PATH=$GOPATH:$PATH
