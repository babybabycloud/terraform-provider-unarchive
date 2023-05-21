package extract

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigExcludeSkip(t *testing.T) {
	conf := &Config{
		Include: func(_ string) bool {
			return true
		},
		Exclude: func(_ string) bool {
			return true
		},
	}

	result := conf.isSkip("any")
	assert.True(t, result)
}

func TestConfigIncludeSkip(t *testing.T) {
	conf := &Config{
		Include: func(_ string) bool {
			return false
		},
		Exclude: func(_ string) bool {
			return false
		},
	}

	result := conf.isSkip("any")
	assert.True(t, result)
}

func TestConfigNoSkip(t *testing.T) {
	conf := &Config{
		Include: func(_ string) bool {
			return true
		},
		Exclude: func(_ string) bool {
			return false
		},
	}

	result := conf.isSkip("any")
	assert.False(t, result)
}

func TestConfigCorrectFileNameFlat(t *testing.T) {
	conf := &Config{
		IsFlat: true,
		Outdir: path.Join("your", "dir"),
	}

	result := conf.correctFileName(path.Join("any", "file.txt"))
	assert.Equal(t, path.Join("your", "dir", "file.txt"), result)
}

func TestConfigCorrectFileNameNoFlat(t *testing.T) {
	conf := &Config{
		IsFlat: false,
		Outdir: path.Join("your", "dir"),
	}

	result := conf.correctFileName(path.Join("any", "file.txt"))
	assert.Equal(t, path.Join("your", "dir", "any", "file.txt"), result)
}
