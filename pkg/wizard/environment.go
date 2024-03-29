package wizard

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/sdk"
)

func (w *Wizard) SelectEnvironment() (*sdk.EnvironmentCollection, error) {
	return w.selectEnvironment(1)
}

func (w *Wizard) selectEnvironment(page int32) (*sdk.EnvironmentCollection, error) {
	model, err := w.getEnvironments(page)
	if err != nil {
		return nil, err
	}

	embedded, ok := model.GetEmbeddedOk()
	if !ok {
		return nil, ErrEmptyListing
	}

	collectionItems := embedded.GetItem()

	items := []string{}
	for _, item := range collectionItems {
		items = append(items, fmt.Sprintf("%s (%s)", item.GetName(), item.GetId()))
	}

	currentPage, totalPages := getPaginationInfo(model)

	index, newPage, err := chooseOrNavigate("Select Environment", items, currentPage, totalPages)
	if err != nil {
		return nil, err
	}

	if newPage != nil {
		return w.selectEnvironment(*newPage)
	}

	if index != nil {
		return &collectionItems[*index], nil
	}

	panic("Something went wrong...")
}

func (w *Wizard) getEnvironments(page int32) (*sdk.PaginatedEnvironmentCollection, error) {
	listOptions := environment.NewListOptions()
	listOptions.Page = page
	listOptions.Organization = w.profile.Context.Organization
	listOptions.Project = w.profile.Context.Project
	listOptions.Profile = w.profile

	return environment.List(listOptions)
}
