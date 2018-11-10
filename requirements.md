

# file2Consul requirements

Requirements for file to consul under consideration for future work

# Agreed for Implementation

* Update data/config/simple samples for a larger example set of properties.
* Update documentation to better explain inheritance override versus interpolation.
* Document what happens when a value that should not be interpolated such as {radomstring} occurs in file and does not have a overlapping variable defined.

# Under Consideration 

* Allow a variable starting with value starting with @ to name a file.  If the file exists it's contents will be read and substituted for the value. The file name will be subjected to interpolation before attempting to open and the contents will also be subjected to interpolation
* Add Leading + semantic to allow longer strings to be defined inline
* Allow a processing Directive @INCLUDE= to cause a file to be read at this time and processed as if it were include in the source file.  NOTE: Need to think about this use case we already have the ability to process files in order so should be able to do the same thing by breaking the files up.
* Ability to suppress processing of files in directory that do not end with specific extensions such as .txt or .ini
* Ability to process locally defined environment variables in addition to variables defined on the command line.
* Ability to support yaml style input in addition to property file syntax
* Allow the list of files to be defined in a external file rather than on the command line since this can make the command line rather long and complex.
* Ability to trigger a git clone during the processing. 
* Ability to log the values that changed and were saved during each pass with a timestamp.  Since we do not save every variable every time  may need to to reconstruct chronology of what is updated.
* Ability to specify file processing order in configuration string when processing all files in directory.   Most likely path for this is to process files in ASCII sorted order then use file naming convention to control processing order.
* Ability to redefine variables defined in command line to by overriding a key.   This allows interpolation and file processing order to change the value of variables through the run.     NOTE:  Need Use Case
* Ability to use multi-threading when submitting values to Consul
* Ability to have a @ style macro to explicitly change value of defined variable.
* Ability to process properties files style [group] and redefine CURRENT_GROUP as a variable that will be interpolation
* 

# Done

* Save value set to cache file and allow to optionally re-use when determining what to send to the consul server
* Ability to process multiple input files and use them all to determine a set of key values saved in consul.
* Ability to use files processed latter in the sequence to override values defined earlier in the processing.
* Ability to interpolate defined variables into configuration patterns without manual editing.
* Ability to override a defined variable defined in another environment.
* Ability to create new keys using interpolate defined variables
* Ability to update multiple consul servers during a single run.
* Ability to generate errors when failure to save to consul or failure to read files. 