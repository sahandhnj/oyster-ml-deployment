package filemanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/sahandhnj/apiclient/util"

	"io"
	"os"
)

const (
	FileStoreManagerDIR = ".oyster"
	ConfigFilePath      = "config.yaml"
	StackFilePath       = "stack.yaml"
)

type FileStoreManager struct {
	DIR        string
	ConfigFile string
}

func NewFileStoreManager() (*FileStoreManager, error) {
	FileStoreManager := &FileStoreManager{
		DIR:        FileStoreManagerDIR,
		ConfigFile: path.Join(FileStoreManagerDIR, ConfigFilePath),
	}

	err := os.MkdirAll(FileStoreManagerDIR, 0755)
	if err != nil {
		return nil, err
	}

	return FileStoreManager, nil
}

func (d *FileStoreManager) ConfigFileExists() (bool, error) {
	return d.FileExists(d.ConfigFile)
}

func (d *FileStoreManager) ReadConfigFile() ([]byte, error) {
	return d.GetFileContent(d.ConfigFile)
}

func (d *FileStoreManager) WriteToConfigFile(content interface{}) error {
	return d.WriteYAMLToFile(d.ConfigFile, content)
}

func (d *FileStoreManager) RemoveDirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
}

func (d *FileStoreManager) GetStackFilePath(stackIdentifier string) string {
	return path.Join(d.DIR, StackFilePath)
}

func (d *FileStoreManager) StoreStackFileFromBytes(stackIdentifier, fileName string, data []byte) error {
	if len(fileName) == 0 {
		fileName = StackFilePath
	}

	stackStorePath := path.Join(d.DIR, fileName)
	r := bytes.NewReader(data)

	err := d.createFileInStore(stackStorePath, r)
	if err != nil {
		return err
	}

	return nil
}

func (d *FileStoreManager) GetFileContent(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *FileStoreManager) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (d *FileStoreManager) WriteToFile(filePath string, content string) error {
	filePath = path.Join(d.DIR, filePath)
	byteContent := []byte(content)

	fmt.Println(filePath)
	return ioutil.WriteFile(filePath, byteContent, 0644)
}

func (d *FileStoreManager) WriteJSONToFile(filePath string, content interface{}) error {
	jsonContent, err := util.MarshalJsonObject(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, jsonContent, 0644)
}

func (d *FileStoreManager) WriteYAMLToFile(filePath string, content interface{}) error {
	yamlContent, err := util.MarshalYamlObject(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, yamlContent, 0644)
}

func (d *FileStoreManager) FileExists(filePath string) (bool, error) {
	_, err := os.Stat(".oyster/config.yaml")

	if err != nil {
		fmt.Println(err)
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *FileStoreManager) CreateDirectoryInStore(name string) error {
	filePath := path.Join(d.DIR, name)
	return os.MkdirAll(filePath, 0700)
}

func (d *FileStoreManager) createFileInStore(filePath string, r io.Reader) error {
	filePath = path.Join(d.DIR, filePath)

	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}
