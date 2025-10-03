package main

import (
	"os"
	"fmt"
	"strings"

	"path/filepath"
)

// if the video is already downloaded, then we can send the existing video
func checkVideoAlreadyDownloaded(dir string, filename string) (string, error) {
	var filePath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			// Get filename without extension since we don't necessarily know
			// what the ext will be (mp4, webm, etc)
			nameWithoutExt := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			if nameWithoutExt == filename {
				filePath = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("Failed to find existing file with base name of %s", filename)
	}
	return filePath, nil
}
