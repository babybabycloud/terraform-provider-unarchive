package unarchive

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type patterns []string

func (p patterns) doesNameMatchPatterns(name string) bool {
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
	return !p.doesNameMatchPatterns(name)
}
