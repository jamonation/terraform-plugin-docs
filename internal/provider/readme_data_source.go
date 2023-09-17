// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

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
	//client *http.Client
	client *string
}

// ReadmeDataSourceModel describes the data source data model.
type ReadmeDataSourceModel struct {
	Intro       types.String `tfsdk:"intro" hcl:"intro" cty:"intro"`
	Body        types.String `tfsdk:"body" hcl:"body" cty:"body"`
	Description types.String `tfsdk:"description" hcl:"description" cty:"description"`
	Image       types.String `tfsdk:"image" hcl:"image" cty:"image"`
}

// Readme describes the data source data model.
type Readme struct {
	Intro       string `tfsdk:"intro" hcl:"intro"`
	Body        string `tfsdk:"body" hcl:"body"`
	Description string `tfsdk:"description" hcl:"description"`
	Image       string `tfsdk:"image" hcl:"image"`
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
		},
	}
}

func (d *ReadmeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*string)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ReadmeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReadmeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var readme Readme
	//fileName := fmt.Sprintf("%s/%s", invokerModule, "README.hcl")
	err := hclsimple.DecodeFile("/Users/jamon/chainguard/images/images/zot/README.hcl", nil, &readme)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unabled to parse README.hcl",
			fmt.Sprintf("%v. Please report this issue to the provider developers.", err),
		)
	}

	tflog.Trace(ctx, fmt.Sprintf("%#v\n", readme))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &readme)...)
}
