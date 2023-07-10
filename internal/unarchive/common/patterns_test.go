package common

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestToPatternNil(t *testing.T) {
	patterns := ToPatterns(types.ListNull(types.StringType))
	assert.Empty(t, patterns)
}

func TestToPattern(t *testing.T) {
	elements := []string{
		"src",
		"Makefile",
		`^main\.py$`,
	}
	list, _ := types.ListValueFrom(context.TODO(), types.StringType, elements)
	p := ToPatterns(list)
	assert.Equal(t, Patterns(elements), p)

}

func TestDoesNameMatchTrue(t *testing.T) {
	elements := Patterns{
		"src",
		"Makefile",
		`^main\.py$`,
	}

	result := elements.doesNameMatch("go/src")
	assert.True(t, result)
}

func TestDoesNameMatchFalse(t *testing.T) {
	elements := Patterns{
		"src",
		"Makefile",
		`^main\.py$`,
	}

	result := elements.doesNameMatch("go/hello.py")
	assert.False(t, result)
}

func TestIncludeZeroLength(t *testing.T) {
	p := make(Patterns, 0)
	assert.True(t, p.Include("any"))
}

func TestIncludeMatch(t *testing.T) {
	p := Patterns{
		"http",
		"net",
	}

	assert.True(t, p.Include("net.py"))
}

func TestExcludeZeroLength(t *testing.T) {
	p := make(Patterns, 0)
	assert.False(t, p.Exclude("any"))
}

func TestExcludeMatch(t *testing.T) {
	p := Patterns{
		"http",
		"net",
	}

	assert.True(t, p.Exclude("http"))
}
