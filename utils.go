package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func fetchSourceFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
func parseFile(fileType string, input []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	var err error

	switch fileType {
	case "json":
		err = json.Unmarshal(input, &result)
	case "yml":
		err = yaml.Unmarshal(input, &result)
	default:
		return nil, ErrInvalidInputType
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func flatInputContent(input map[string]interface{}) map[string]string {
	result := map[string]string{}

	for k, v := range input {
		flatten(k, v, result)
	}

	return result
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkIfString(v interface{}) bool {
	_, isString := v.(string)
	return isString
}

func flatten(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			flatten(prefix+"."+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		flatten(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			flatten(fmt.Sprintf("%s.%d", prefix, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}

func unflatten(flat map[string]interface{}) (map[string]interface{}, error) {
	unflat := map[string]interface{}{}

	for key, value := range flat {
		keyParts := strings.Split(key, ".")

		// Walk the keys until we get to a leaf node.
		m := unflat
		for i, k := range keyParts[:len(keyParts)-1] {
			v, exists := m[k]
			if !exists {
				newMap := map[string]interface{}{}
				m[k] = newMap
				m = newMap
				continue
			}

			innerMap, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("key=%v is not an object", strings.Join(keyParts[0:i+1], "."))
			}
			m = innerMap
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%v already exists", key)
		}
		m[keyParts[len(keyParts)-1]] = value
	}

	return unflat, nil
}

func saveToFile(fileType string, fileContent map[string]string, filePath string) error {
	var outputFileContent []byte
	var err error

	switch fileType {
	case "json":
		outputFileContent, err = json.MarshalIndent(fileContent, "", " ")
	case "yml":
		m2 := make(map[string]interface{}, len(fileContent))
		for k, v := range fileContent {
			m2[k] = v
		}
		outputFileUnflatten, err := unflatten(m2)
		if err != nil {
			return err
		}
		outputFileContent, err = yaml.Marshal(outputFileUnflatten)
	default:
		return ErrInvalidInputType
	}

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, outputFileContent, 0644)
	return err
}
