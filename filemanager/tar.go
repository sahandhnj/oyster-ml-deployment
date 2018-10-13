package filemanager

import (
	"archive/tar"
	"bytes"
)

func (f *FileStoreManager) TarFile(fileName string) (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	content, err := f.GetFileContent(fileName)
	if err != nil {
		return nil, err
	}

	tarHeader := &tar.Header{
		Name: fileName,
		Size: int64(len(content)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(content)
	if err != nil {
		return nil, err
	}

	tarReader := bytes.NewReader(buf.Bytes())

	return tarReader, nil
}
