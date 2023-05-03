package unarchive

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {

	includes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("src"),
	})
	excludes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("test"),
	})
	model := zipFileDataSourceModel{
		FileName: types.StringValue("master.zip"),
		Includes: includes,
		Excludes: excludes,
	}

	_, err := model.extract()
	assert.Empty(t, err)
}
