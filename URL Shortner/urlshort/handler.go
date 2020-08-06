package urlshort

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" //import mysql driver to connect to MySQL database
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler perfoms the same actions as YAMLHandler but parses provided JSON data
func JSONHandler(js []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(js)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

//DatabaseHandler performs same actions as YAMLHandler but retrives patterns from MySQL Database
func DatabaseHandler(fallback http.Handler, databaseName string) (http.HandlerFunc, error) {
	pathMap, err := connectDB(databaseName)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathMap, fallback), nil
}

//ParseYamlJSON is the structure of given patterns
type ParseYamlJSON struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

/*
	Receives YAML byte array and returns list of ParseYamlJSON type after parsing
*/
func parseYAML(yml []byte) ([]ParseYamlJSON, error) {
	var p []ParseYamlJSON
	err1 := yaml.Unmarshal(yml, &p)
	if err1 != nil {
		log.Fatal(err1)
		return nil, err1
	}
	return p, nil
}

/*
	Receives JSON byte array and returns list of ParseYamlJSON type after parsing
*/
func parseJSON(yml []byte) ([]ParseYamlJSON, error) {
	var p []ParseYamlJSON
	err1 := json.Unmarshal(yml, &p)
	if err1 != nil {
		log.Fatal(err1)
		return nil, err1
	}
	return p, nil
}

/*
	Receives list of ParseYamlJSON and returns the converted patterns map
*/
func buildMap(patterns []ParseYamlJSON) map[string]string {
	m := make(map[string]string)

	for _, pattern := range patterns {
		m[pattern.Path] = pattern.URL
	}
	return m
}

/*
	connectDB connects to the database, retrives patterns records and returns them after them to map
*/
func connectDB(databaseName string) (map[string]string, error) {
	db, err := sql.Open("mysql", "manas:CERPM4MaxzMzrGhD@tcp(127.0.0.1:3306)/"+databaseName)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	results, err := db.Query("SELECT * FROM patterns")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	m := make(map[string]string)
	for results.Next() {
		var p ParseYamlJSON
		// for each row, scan the result into our map
		err = results.Scan(&p.Path, &p.URL)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		// add to map
		m[p.Path] = p.URL
	}
	defer db.Close()
	return m, nil

}
