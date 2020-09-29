package Sender

import (
	"os"
	"path/filepath"
)

type DataPack struct {
	FileName string
	FileSize int64
	Data     []byte
}

var dataPack DataPack

func Send(dir string, ip string) {
	var fileArry []string
	if isFile(dir) {
		dataPack.FileName = dir
	} else {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fileArry = append(fileArry, path)
			}
			return nil
		})
	}
}

func isFile(dir string) bool {
	f, _ := os.Stat(dir)
	if f.IsDir() {
		return false
	}
	return true
}

func (d *DataPack) Pack() {

}
