package test

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func getBaseSchemaVersion() (string, error) {
	file, err := os.Open(schemaPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
		return "", err
	}
	imports := data["imports"].([]interface{})
	return imports[1].(string), nil
}
