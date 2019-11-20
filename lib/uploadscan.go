package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func ScanFolder() {
	var files []string

	root := "uploaded-images"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files[1:] {


		fmt.Println(file)


	}
}