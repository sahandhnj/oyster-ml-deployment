package filemanager

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func (f *FileStoreManager) CopyToStore(source, destination string) error {
	destination = path.Join(f.DIR, destination)

	sourceFileStat, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", source)
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)

	return nil
}

func GetFileContent(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *FileStoreManager) ReadRQLineByLine(filepath string) string {
	filepath = path.Join(f.DIR, filepath)
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var stringArray []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stringArray = append(stringArray, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strings.Join(stringArray, " ")
}
