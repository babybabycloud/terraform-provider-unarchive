package unarchive

import (
	"context"

	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/extract"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const DEFAULT_DIR_MODE = 0740

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &unarchiveDataSource{}
)

// NewUnarchiveDataSource is a helper function to simplify the provider implementation.
func NewUnarchiveDataSource() datasource.DataSource {
	return &unarchiveDataSource{}
}

// unarchiveDataSource is the data source implementation.
type unarchiveDataSource struct{}

// Metadata returns the data source type name.
func (d *unarchiveDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema defines the schema for the data source.
func (d *unarchiveDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"file_name": schema.StringAttribute{
				Required: true,
			},
			"output": schema.StringAttribute{
				Optional: true,
			},
			"includes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"excludes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"flat": schema.BoolAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"file_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *unarchiveDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model unarchiveDataSourceModel
	diag := req.Config.Get(ctx, &model)
	if diag.HasError() {
		for _, diagnotic := range diag.Errors() {
			resp.Diagnostics.AddError(diagnotic.Summary(), diagnotic.Detail())
		}
		return
	}

	conf := &extract.Config{
		Ctx:     ctx,
		Name:    model.FileName.ValueString(),
		Include: model.includePatterns(),
		Exclude: model.excludePatterns(),
		Outdir:  model.decideOutputDir(),
		IsFlat:  model.isFlat(),
		Type:    model.Type.ValueString(),
	}
	info := extract.Extract(conf)

	if info.Err != nil {
		resp.Diagnostics.AddError(info.Msg, info.Err.Error())
		return
	}

	model.addFileNames(info.FileNames)
	resp.State.Set(ctx, &model)
}
