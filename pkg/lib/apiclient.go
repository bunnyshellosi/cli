package lib

import (
	"io"
	"net/http"
	"os"

	bunnysdk "bunnyshell.com/sdk"
)

func GetOrganizations() (*bunnysdk.PaginatedOrganizationCollection, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().OrganizationApi.OrganizationList(ctx)

	return request.Execute()
}

func GetOrganization(organizationID string) (*bunnysdk.OrganizationItem, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().OrganizationApi.OrganizationView(ctx, organizationID)

	return request.Execute()
}

func GetProjects(organization string) (*bunnysdk.PaginatedProjectCollection, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().ProjectApi.ProjectList(ctx)
	if organization != "" {
		request = request.Organization(organization)
	}

	return request.Execute()
}

func GetProject(projectID string) (*bunnysdk.ProjectItem, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().ProjectApi.ProjectView(ctx, projectID)

	return request.Execute()
}

func GetEnvironments(projectID string) (*bunnysdk.PaginatedEnvironmentCollection, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().EnvironmentApi.EnvironmentList(ctx)
	if projectID != "" {
		request = request.Project(projectID)
	}

	return request.Execute()
}

func GetEnvironment(environmentID string) (*bunnysdk.EnvironmentItem, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().EnvironmentApi.EnvironmentView(ctx, environmentID)

	return request.Execute()
}

func GetComponents(environment, operationStatus string) (
	*bunnysdk.PaginatedComponentCollection,
	*http.Response,
	error,
) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().ComponentApi.ComponentList(ctx)
	if environment != "" {
		request = request.Environment(environment)
	}

	if operationStatus != "" {
		request = request.OperationStatus(operationStatus)
	}

	return request.Execute()
}

func GetComponent(componentID string) (*bunnysdk.ComponentItem, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().ComponentApi.ComponentView(ctx, componentID)

	return request.Execute()
}

func GetComponentResources(componentID string) ([]bunnysdk.ComponentResourceItem, *http.Response, error) {
	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().ComponentApi.ComponentResources(ctx, componentID)

	return request.Execute()
}

func DownloadEnvironmentKubeConfig(kubeConfigPath, environmentID string) error {
	kubeConfigFile, err := os.Create(kubeConfigPath)
	if err != nil {
		return err
	}
	defer kubeConfigFile.Close()

	ctx, cancel := GetContext()
	defer cancel()

	request := GetAPI().EnvironmentApi.EnvironmentKubeConfig(ctx, environmentID)

	_, resp, err := request.Execute()
	if err != nil && err.Error() != "undefined response type" {
		return err
	}

	_, err = io.Copy(kubeConfigFile, resp.Body)

	return err
}
