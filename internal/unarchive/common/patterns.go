package common

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Patterns []string

// ToPatterns converts string elements in types.List to []string
func ToPatterns(list types.List) Patterns {
	var p Patterns
	if !list.IsNull() {
		PatternsInType := make([]types.String, len(list.Elements()))
		p = make(Patterns, len(list.Elements()))
		list.ElementsAs(context.TODO(), &PatternsInType, false)
		for index, value := range PatternsInType {
			p[index] = value.ValueString()
		}
	}
	return p
}

// DoesNameMatch checkes if the name matches any value in Patterns
func (p Patterns) doesNameMatch(name string) bool {
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

// Include checks if the name is included in the patters
func (p Patterns) Include(name string) bool {
	if len(p) == 0 {
		return true
	}
	return p.doesNameMatch(name)
}

// Exclude checks if the name is not included in the patters
func (p Patterns) Exclude(name string) bool {
	if len(p) == 0 {
		return false
	}
	return p.doesNameMatch(name)
}
