package extract

import (
	"context"
	"path"
	"testing"

	"github.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/common"
	"github.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/model"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestFromFilterModel(t *testing.T) {
	includes, _ := basetypes.NewListValueFrom(context.TODO(), basetypes.StringType{}, []string{"include"})
	excludes, _ := basetypes.NewListValueFrom(context.TODO(), basetypes.StringType{}, []string{"exclude"})
	filterModel := model.FilterModel{
		Includes: includes,
		Excludes: excludes,
	}

	configFilter := FromFilterModel(filterModel)
	expected := ConfigFilter{
		Includes: common.Patterns{"include"},
		Excludes: common.Patterns{"exclude"},
	}

	assert.EqualValues(t, expected, configFilter)
}

func TestFromUnarchiveDataSourceModel(t *testing.T) {
	includes, _ := basetypes.NewListValueFrom(context.TODO(), types.StringType, []string{"include pattern"})
	excludes, _ := basetypes.NewListValueFrom(context.TODO(), types.StringType, []string{"exclude pattern"})
	m := model.UnarchiveDataSourceModel{
		FileName: types.StringValue("A file name"),
		Output:   types.StringValue("Destination"),
		Filters: []model.FilterModel{
			{
				Includes: includes,
				Excludes: excludes,
			},
		},
		Flat: types.BoolValue(true),
		Type: types.StringValue(".tar.gz"),
	}
	config := FromUnarchiveDataSourceModel(m)
	expected := &Config{
		Name: "A file name",
		Filters: []ConfigFilter{
			{
				Includes: common.Patterns{"include pattern"},
				Excludes: common.Patterns{"exclude pattern"},
			},
		},
		Outdir: "Destination",
		IsFlat: true,
		Type:   ".tar.gz",
	}

	assert.EqualValues(t, expected, config)
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
