package config

import (
	"os"

	"github.com/aws/jsii-runtime-go"
)

func mustReadFromFile(path string) string {
	out, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func loadFromFile(path string) *string {
	return jsii.String(mustReadFromFile(path))
}
