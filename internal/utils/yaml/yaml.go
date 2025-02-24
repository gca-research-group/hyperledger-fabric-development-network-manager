package yaml

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse[T any](data T) string {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err := encoder.Encode(&data)

	if err != nil {
		panic(fmt.Sprintf("Error marshalling to YAML: %v", err))
	}

	encoder.Close()

	content := buf.String()

	content = strings.ReplaceAll(content, "''", "'")
	content = strings.ReplaceAll(content, "'\"", "\"")
	content = strings.ReplaceAll(content, "\"'", "\"")

	return content
}

func Write(data string, filename string) {

	// data = strings.Replace(data, "- Name: Orderer", "- &Orderer\n    Name: Orderer", 1)

	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %v\n", err))
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		panic(fmt.Sprintf("Error writing to file: %v\n", err))
	}
}
