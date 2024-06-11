package filestorage

import (
	"bufio"
	"fmt"
	"os"
)

type FileStorage struct {
	fileNeme string
}

func New(fileNeme string) *FileStorage {
	return &FileStorage{
		fileNeme: fileNeme,
	}
}

func (f *FileStorage) Save(data []byte) error {
	file, err := os.OpenFile(fmt.Sprintf("./%s", f.fileNeme), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) GetAll() ([]string, error) {
	file, err := os.OpenFile(fmt.Sprintf("./%s", f.fileNeme), os.O_RDONLY, 0o600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
