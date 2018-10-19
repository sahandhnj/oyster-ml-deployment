package filemanager

import (
	"fmt"
	"path"

	"github.com/mholt/archiver"
)

func (f *FileStoreManager) CTarGz(output string, content []string, internal bool) error {
	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// abPath := filepath.Dir(ex)

	output = path.Join(f.DIR, output)
	fmt.Println("writing to " + output)
	if internal {
		for i, file := range content {
			content[i] = path.Join(f.DIR, file)
		}
	}

	return archiver.TarGz.Make(output, content)
}

func (f *FileStoreManager) XTarGz(input string, output string) error {
	output = path.Join(f.DIR, output)
	input = path.Join(f.DIR, input)
	fmt.Println(input, output)
	return archiver.TarGz.Open(input, output)
}
