package endpoint

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type ListOptions struct {
	common.ListOptions

	Organization string
	Project      string
	Environment  string

	Name string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	lo.updateSelfFlags(flags)
	lo.ListOptions.UpdateFlagSet(flags)
}

func (lo *ListOptions) updateSelfFlags(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Name, "component", lo.Name, "Filter by Component Name")
}

func List(options *ListOptions) (*sdk.PaginatedComponentEndpointCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedComponentEndpointCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ComponentEndpointApi.ComponentEndpointList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiComponentEndpointListRequest, options *ListOptions) sdk.ApiComponentEndpointListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}

	if options.Project != "" {
		request = request.Project(options.Project)
	}

	if options.Environment != "" {
		request = request.Environment(options.Environment)
	}

	if options.Name != "" {
		request = request.Name(options.Name)
	}

	return request
}
