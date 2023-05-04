package unarchive

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type patterns []string

func toPatterns(list types.List) patterns {
	var p patterns
	if !list.IsNull() {
		patternsInType := make([]types.String, len(list.Elements()))
		p = make(patterns, len(list.Elements()))
		for index, value := range patternsInType {
			p[index] = value.ValueString()
		}
	}
	return p
}

func (p patterns) doesNameMatchPatterns(name string) bool {

	if len(p) == 0 {
		return true
	}

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

func (p patterns) doesNotNameMatchPatterns(name string) bool {
	if len(p) == 0 {
		return true
	}

	return !p.doesNameMatchPatterns(name)
}
