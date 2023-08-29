package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var err error
	cwd, _ := os.Getwd()
	file := cwd + "/yafgo/yafgo.txt"
	err = Create(file, []byte(`yafgo`))
	assert.NoError(t, err)

	assert.Equal(t, 1, GetLineNum(file))
	assert.True(t, Exists(file))
	assert.True(t, Remove(file))
	assert.True(t, Remove(cwd+"/yafgo"))
}

func TestExtension(t *testing.T) {
	file := "path/to/file.go"
	extension := Extension(file)
	assert.NotEmpty(t, extension)
	assert.Equal(t, "go", extension)

	baseName := BaseName(file)
	assert.Equal(t, "file.go", baseName)

	baseNameWithoutExt := BaseNameWithoutExtentsion(file)
	assert.Equal(t, "file", baseNameWithoutExt)
}
