package filemanager

import (
	"path"

	"github.com/mholt/archiver"
)

func (f *FileStoreManager) CTarGz(output string, content []string, internal bool) error {
	output = path.Join(f.DIR, output)

	if internal {
		for i, file := range content {
			content[i] = path.Join(f.DIR, file)
		}
	}

	return archiver.Zip.Make(output, content)
}

func (f *FileStoreManager) XTarGz(input string, output string) error {
	output = path.Join(f.DIR, output)
	input = path.Join(f.DIR, input)

	return archiver.Zip.Open(input, output)
}
