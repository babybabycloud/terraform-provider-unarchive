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
	assert.Equal(t, patterns(elements), p)

}

func TestDoesNameMatchTrue(t *testing.T) {
	elements := patterns{
		"src",
		"Makefile",
		`^main\.py$`,
	}

	result := elements.DoesNameMatch("go/src")
	assert.True(t, result)
}

func TestDoesNameMatchFalse(t *testing.T) {
	elements := patterns{
		"src",
		"Makefile",
		`^main\.py$`,
	}

	result := elements.DoesNameMatch("go/hello.py")
	assert.False(t, result)
}
