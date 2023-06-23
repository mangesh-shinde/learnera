package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadCoursesFile() []byte {
	currentDir, _ := os.Getwd()
	fmt.Println(currentDir)
	dataFilePath := filepath.Join(currentDir, "data", "courses.json")
	file, err := os.Open(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	bytesData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return bytesData
}
