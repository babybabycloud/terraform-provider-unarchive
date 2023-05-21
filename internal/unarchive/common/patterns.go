package common

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type patterns []string

// ToPatterns converts string elements in types.List to []string
func ToPatterns(list types.List) patterns {
	var p patterns
	if !list.IsNull() {
		patternsInType := make([]types.String, len(list.Elements()))
		p = make(patterns, len(list.Elements()))
		list.ElementsAs(context.TODO(), &patternsInType, false)
		for index, value := range patternsInType {
			p[index] = value.ValueString()
		}
	}
	return p
}

// DoesNameMatch checkes if the name matches any value in patterns
func (p patterns) DoesNameMatch(name string) bool {
	for _, value := range p {
		matched, err := regexp.MatchString(value, name)

		if err != nil {
			tflog.Warn(context.Background(), fmt.Sprintf("%s. Ignore it", err.Error()))
		}
		if matched {
			return true
		}
	}
	return false
}
