package actions

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

)


func ReplaceInAllFiles(old string, new string) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		fmt.Println("----")
		fmt.Println(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		data = bytes.ReplaceAll(data, []byte(old), []byte(new))
		fmt.Println(string(data))
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return err
		}
		return nil
	})
}
