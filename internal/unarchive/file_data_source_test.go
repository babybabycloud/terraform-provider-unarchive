package unarchive

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue("master.zip"),
	}
	info := model.extract(context.TODO())
	assert.Empty(t, info.err)
}

func TestIncludeFilter(t *testing.T) {
	includes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("Makefile"),
	})
	model := fileDataSourceModel{
		FileName: types.StringValue("master.zip"),
		Includes: includes,
	}

	info := model.extract(context.TODO())
	assert.Empty(t, info.err)
}

func TestExcludeFilter(t *testing.T) {
	excludes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("src"),
	})
	model := fileDataSourceModel{
		FileName: types.StringValue("master.zip"),
		Excludes: excludes,
	}

	info := model.extract(context.TODO())
	assert.Empty(t, info.err)
}
