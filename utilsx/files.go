package utilsx

import (
	"fmt"
	"io"
	"os"
)

func ReadFiles(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open the file [%s]! Error is [%s]\n", filePath, err.Error())
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file %s. %v", filePath, err.Error())
	}
	return b, nil
}

func EnsureDirExist(dir string) error {
	if IsDirExist(dir) {
		return nil
	}

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("failed to mkdir the dir %s. %v", dir, err.Error())
	}

	return nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDirExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return s.IsDir()
		}

		return false
	}

	return s.IsDir()
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
