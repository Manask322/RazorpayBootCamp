package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./urlshort"
)

/*
	Receives file path as an argument and returns byte array of the read file
*/
func readFile(filepath string) []byte {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}
	return f
}

/*
	Main function of the program
*/
func main() {

	// Flags for the program
	yamlFilePath := flag.String("ymlFile", "yaml", "path to of the yaml file .yaml")
	jsonFilePath := flag.String("jsonFile", "json", "path of the json file .json")
	database := flag.Bool("database", false, "use database to retrive path")
	databaseName := flag.String("databaseName", "golang", "Enter the databaseName")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	//read yaml or json files if specified
	yamlFile := readFile(*yamlFilePath)
	jsonFile := readFile(*jsonFilePath)

	handler := mapHandler // choose maphandler as the default handler
	var err error

	// perform appropriate action based on flags specified by the user
	if jsonFile != nil {
		handler, err = urlshort.JSONHandler(jsonFile, mapHandler)
	} else if yamlFile != nil {
		handler, err = urlshort.YAMLHandler(yamlFile, mapHandler)
	} else if *database {
		handler, err = urlshort.DatabaseHandler(mapHandler, *databaseName)
	}

	if err != nil {
		log.Fatal(err) // print custom handler error if any.
	}

	fmt.Println("Starting the server on :3000")
	serveErr := http.ListenAndServe(":3000", handler)
	log.Fatal(serveErr) // log server error if any
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
