// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const readmeTemplate = `name        = "%s"
image       = ""
intro       = ""
body        = ""
description = ""
`

// Ensure provider defined types fully satisfy framework interfaces.
var (
	//_ datasource.DataSource              = &ReadmeDataSource{}
	_ datasource.DataSourceWithConfigure = &ReadmeDataSource{}
)

func NewReadmeDataSource() datasource.DataSource {
	return &ReadmeDataSource{}
}

// ReadmeDataSource defines the data source implementation.
type ReadmeDataSource struct {
	popts ProviderOpts
}

// ReadmeDataSourceModel describes the data source data model.
type ReadmeDataSourceModel struct {
	Body        types.String `tfsdk:"body" hcl:"body" cty:"body"`
	Description types.String `tfsdk:"description" hcl:"description" cty:"description"`
	Image       types.String `tfsdk:"image" hcl:"image" cty:"image"`
	Intro       types.String `tfsdk:"intro" hcl:"intro" cty:"intro"`
	Name        string       `tfsdk:"name"`
}

// Readme describes the data source data model.
type Readme struct {
	Body        string `tfsdk:"body" hcl:"body"`
	Description string `tfsdk:"description" hcl:"description"`
	Image       string `tfsdk:"image" hcl:"image"`
	Intro       string `tfsdk:"intro" hcl:"intro"`
	Name        string `tfsdk:"name" hcl:"name"`
}

func (d *ReadmeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_readme"
}

func (d *ReadmeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"intro": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"body": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"image": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

// copied from https://github.com/chainguard-dev/terraform-provider-apko/blob/55ea67c749a662a8c27f64c5f6d47576308a997d/internal/provider/config_data_source.go#L94-L106
func (d *ReadmeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	popts, ok := req.ProviderData.(*ProviderOpts)
	if !ok || popts == nil {
		resp.Diagnostics.AddError("Client Error", "invalid provider data")
		return
	}
	d.popts = *popts
}

func (d *ReadmeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReadmeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var readme Readme

	tflog.Trace(ctx, fmt.Sprintf("got repos: %v", data.Name))

	fileName := fmt.Sprintf("images/%s/%s", data.Name, "README.hcl")
	_, err := os.Stat(fileName)
	if err != nil {
		resp.Diagnostics.AddError(
			"README.hcl error",
			fmt.Sprintf("Missing or inaccessible file %v.\nCopy /tmp/README.hcl for an example template.", fileName),
		)
		sampleReadme := fmt.Sprintf(readmeTemplate, data.Name)
		_ = os.WriteFile("/tmp/README.hcl", []byte(sampleReadme), os.FileMode(0o644))
		return
	}

	err = hclsimple.DecodeFile(fileName, nil, &readme)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse README.hcl",
			fmt.Sprintf("%v. URL HERE", err),
		)
	}
	readme.Name = data.Name

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readme)...)
}
