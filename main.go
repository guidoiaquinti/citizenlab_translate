//
// Automatically fetch and translate https://github.com/CitizenLabDotCo/citizenlab
// language file using "AWS Translate" service
//
// TODO
// - handle transient failures
//
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// const SOURCE_FILE_URL = []string{
// 	"https://raw.githubusercontent.com/CitizenLabDotCo/citizenlab/master/front/app/translations/en.json",
// }

func main() {

	// Check if AWS env variables are set

	// Fetch translation files
	sourceFile, err := getSourceFile("https://raw.githubusercontent.com/CitizenLabDotCo/citizenlab/master/front/app/translations/en.json")
	if err != nil {
		log.Fatalf("Error: %s", err)
		panic(err)
	}

	// Translate
	fmt.Println(len(sourceFile))

	// Save output
}

func getSourceFile(url string) (map[string]string, error) {
	// Fetch file
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result map[string]string
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func translateFile(sourceFile map[string]string) map[string]string {
	var result map[string]string

	for k, v := range sourceFile {
		fmt.Printf("KEY: %s\n", k)
		fmt.Printf("VALUE: %s\n\n", v)
	}
}
