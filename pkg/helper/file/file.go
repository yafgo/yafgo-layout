package file

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"
)

func Create(file string, content ...[]byte) (err error) {
	err = os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	// Write file content
	if len(content) > 0 {
		_, err = f.Write(content[0])
	}

	return
}

func PutContent(file string, data []byte) error {
	err := os.WriteFile(file, data, 0644)
	return err
}

func Exists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Remove(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		dir, err := os.ReadDir(file)

		if err != nil {
			return false
		}

		for _, d := range dir {
			err := os.RemoveAll(path.Join([]string{file, d.Name()}...))
			if err != nil {
				return false
			}
		}
	}

	err = os.Remove(file)
	return err == nil
}

func Contain(file string, search string) bool {
	if Exists(file) {
		data, err := os.ReadFile(file)
		if err != nil {
			return false
		}
		return strings.Contains(string(data), search)
	}

	return false
}

// BaseName
//
//	"/path/to/file.ext" => "file.ext"
func BaseName(file string) string {
	return path.Base(file)
}

// BaseNameWithoutExtentsion
//
//	"/path/to/file.ext" => "file"
func BaseNameWithoutExtentsion(file string) string {
	baseName := path.Base(file)
	return strings.TrimSuffix(baseName, path.Ext(file))
}

// Extentsion
//
//	"/path/to/file.ext" => "ext"
func Extension(file string) string {
	return strings.ReplaceAll(path.Ext(file), ".", "")
}

func GetLineNum(file string) (total int) {
	f, _ := os.OpenFile(file, os.O_RDONLY, 0444)
	buf := bufio.NewReader(f)
	defer func() {
		f.Close()
	}()

	for {
		_, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				total++

				break
			}
		} else {
			total++
		}
	}

	return
}
