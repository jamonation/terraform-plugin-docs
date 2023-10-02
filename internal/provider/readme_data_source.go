// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const readmeTemplate = `name        = "%s"
image       = "%s"
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
	ImagePath   types.String `tfsdk:"image_path"`
	FileName    types.String `tfsdk:"file_name"`
}

// Readme describes the data source data model.
type Readme struct {
	Body        string `tfsdk:"body" hcl:"body"`
	Description string `tfsdk:"description" hcl:"description"`
	Image       string `tfsdk:"image" hcl:"image"`
	Intro       string `tfsdk:"intro" hcl:"intro"`
	Name        string `tfsdk:"name" hcl:"name"`
}

type completeReadme struct {
	Body        string `tfsdk:"body" hcl:"body"`
	Description string `tfsdk:"description" hcl:"description"`
	Image       string `tfsdk:"image" hcl:"image"`
	Intro       string `tfsdk:"intro" hcl:"intro"`
	Name        string `tfsdk:"name" hcl:"name"`
	ImagePath   string `tfsdk:"image_path"`
	FileName    string `tfsdk:"file_name"`
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
			"file_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"image_path": schema.StringAttribute{
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
	imagePath, fileName := getFileName(data)

	fullPath := fmt.Sprintf("%s/%s", imagePath, fileName)
	_, err := os.Stat(fullPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"README.hcl error",
			fmt.Sprintf("%s missing or inaccessible.\nRun `cp /tmp/%s %s` and edit it.", fullPath, fileName, fullPath),
		)
		name := strings.Replace(data.Name, ".", "-", -1)
		image := fmt.Sprintf("cgr.dev/chainguard/%s", name)
		sampleReadme := fmt.Sprintf(readmeTemplate, name, image)
		_ = os.WriteFile("/tmp/"+fileName, []byte(sampleReadme), os.FileMode(0o644))
		return
	}

	err = hclsimple.DecodeFile(fullPath, nil, &readme)
	if err != nil {
		errSummary := fmt.Sprintf("Unable to parse %s", fullPath)
		resp.Diagnostics.AddError(
			errSummary,
			fmt.Sprintf("%v", err),
		)
	}

	fullReadme := &completeReadme{
		Body:        readme.Body,
		Description: readme.Description,
		Image:       readme.Image,
		Intro:       readme.Intro,
		Name:        readme.Name,
		ImagePath:   imagePath,
		FileName:    strings.Replace(fileName, ".hcl", ".md", -1),
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &fullReadme)...)
}

// checks for variants like cert-manager.acmesolver, ensure they get their own template
func getFileName(data ReadmeDataSourceModel) (string, string) {
	var imagePath, fileName string

	splitName := strings.Split(data.Name, ".")
	imagePath = fmt.Sprintf("images/%s", splitName[0])
	if len(splitName) > 1 {
		variant := splitName[len(splitName)-1]
		fileName = fmt.Sprintf("README.%s.hcl", variant)
	} else {
		fileName = fmt.Sprintf("README.hcl")
	}
	return imagePath, fileName
}
