package service

import (
	"os"
	"path/filepath"
)

func DownImplants(filename string) ([]byte, error) {
	var er error
	Filepath, _ := filepath.Abs(filepath.Join(".", "implants", filename))
	fileContents, err := os.ReadFile(Filepath)
	if err != nil {
		return nil, er
	}
	return fileContents, nil
}
