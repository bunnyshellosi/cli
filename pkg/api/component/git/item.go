package git

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func NewItemOptions(id string) *common.ItemOptions {
	return common.NewItemOptions(id)
}

func Get(options *common.ItemOptions) (*sdk.ComponentGitItem, error) {
	model, resp, err := GetRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func GetRaw(options *common.ItemOptions) (*sdk.ComponentGitItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ComponentGitApi.ComponentGitView(ctx, options.ID)

	return request.Execute()
}
