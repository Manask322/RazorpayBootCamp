## Steps to run to program

* ### Enter `go run main.go --help` to view all the flags available
    * `-database=<bool>` choose to use database of not
        * if yes then use `-databaseName=<string:database name>` to give database name (use given dump file to create testing database)
    * `-jsonFile=<path>` specify path to JSON file to retrive URL patterns
    * `-ymlFile=<path>` specify path to YAML file to retrive URL patterns
* ### To build the application switch to top directory and enter `go build .`
    * To run the build file enter 
        * `.\URL_Shortner.exe` for Windows
        * `./URL_Shortner` for linux or Mac
