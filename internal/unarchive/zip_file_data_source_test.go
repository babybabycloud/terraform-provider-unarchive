package unarchive

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	model := zipFileDataSourceModel{
		FileName: types.StringValue("master.zip"),
	}
	_, err := model.extract(context.TODO())
	assert.Empty(t, err)
}

func TestIncludeFilter(t *testing.T) {
	includes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("Makefile"),
	})
	model := zipFileDataSourceModel{
		FileName: types.StringValue("master.zip"),
		Includes: includes,
	}

	_, err := model.extract(context.TODO())
	assert.Empty(t, err)
}

func TestExcludeFilter(t *testing.T) {
	excludes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("src"),
	})
	model := zipFileDataSourceModel{
		FileName: types.StringValue("master.zip"),
		Excludes: excludes,
	}

	_, err := model.extract(context.TODO())
	assert.Empty(t, err)
}
