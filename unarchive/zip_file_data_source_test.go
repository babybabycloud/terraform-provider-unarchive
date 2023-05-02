package unarchive

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestReadFile(t *testing.T) {
	model := &zipFileDataSourceModel{
		FileName: types.StringValue("h.zip"),
	}
	err := model.extract()
	if err != nil {
		t.Fatal("aaa ", err)
	}
}
