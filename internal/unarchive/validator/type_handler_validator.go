package validator

import (
	"context"
	"fmt"
	"strings"

	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/extract"
	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/model"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type HandlerTypeValidator struct{}

// Description implements github.com/hashicorp/terraform-plugin-framework/schema/validator.Describer
func (h *HandlerTypeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Possible values: %s", strings.Join([]string{
		extract.ZIP,
		extract.TAR,
		extract.TARGZ,
	}, ", "))
}

// MarkdownDescription implements github.com/hashicorp/terraform-plugin-framework/schema/validator.Describer
func (h *HandlerTypeValidator) MarkdownDescription(ctx context.Context) string {
	return h.Description(ctx)
}

// ValidateString implements github.com/hashicorp/terraform-plugin-framework/schema/validator.String
func (h *HandlerTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var model model.UnarchiveDataSourceModel
	diag := req.Config.Get(ctx, &model)
	if diag.HasError() {
		for _, diagnotic := range diag.Errors() {
			resp.Diagnostics.AddError(diagnotic.Summary(), diagnotic.Detail())
		}
		return
	}

	switch model.Type.ValueString() {
	case extract.ZIP, extract.TAR, extract.TARGZ:
	default:
		resp.Diagnostics.AddError(`Invalid attribute "type".`, `The type attribute only supports ".zip", ".tar", ".tar.gz"`)
	}

}
