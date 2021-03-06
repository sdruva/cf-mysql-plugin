package cfmysql

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	sdkModels "code.cloudfoundry.org/cli/plugin/models"
	pluginModels "github.com/andreasf/cf-mysql-plugin/cfmysql/models"
	"github.com/andreasf/cf-mysql-plugin/cfmysql/resources"
)

//go:generate counterfeiter . ApiClient
type ApiClient interface {
	GetServiceBindings(cliConnection plugin.CliConnection) ([]pluginModels.ServiceBinding, error)
	GetServiceInstances(cliConnection plugin.CliConnection) ([]pluginModels.ServiceInstance, error)
	GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error)
}

func NewApiClient(httpClient HttpWrapper) *apiClient {
	return &apiClient{
		httpClient: httpClient,
	}
}

type apiClient struct {
	httpClient HttpWrapper
}

func (self *apiClient) GetServiceInstances(cliConnection plugin.CliConnection) ([]pluginModels.ServiceInstance, error) {
	var err error
	var allInstances []pluginModels.ServiceInstance
	nextUrl := "/v2/service_instances"

	for nextUrl != "" {
		instanceResponse, err := self.getFromCfApi(nextUrl, cliConnection)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve service instances: %s", err)
		}

		var instances []pluginModels.ServiceInstance
		nextUrl, instances, err = deserializeInstances(instanceResponse)
		allInstances = append(allInstances, instances...)
	}

	return allInstances, err
}

func (self *apiClient) getFromCfApi(path string, cliConnection plugin.CliConnection) ([]byte, error) {
	endpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		return nil, fmt.Errorf("Unable to get API endpoint: %s", err)
	}

	accessToken, err := cliConnection.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("Unable to get access token: %s", err)
	}

	sslDisabled, err := cliConnection.IsSSLDisabled()
	if err != nil {
		return nil, fmt.Errorf("Unable to check SSL status: %s", err)
	}
	return self.httpClient.Get(endpoint+path, accessToken, sslDisabled)
}

func deserializeInstances(jsonResponse []byte) (string, []pluginModels.ServiceInstance, error) {
	paginatedResources := new(resources.PaginatedServiceInstanceResources)
	err := json.Unmarshal(jsonResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("Unable to deserialize service instances: %s", err)
	}

	return paginatedResources.NextUrl, paginatedResources.ToModel(), nil
}

func (self *apiClient) GetServiceBindings(cliConnection plugin.CliConnection) ([]pluginModels.ServiceBinding, error) {
	var allBindings []pluginModels.ServiceBinding
	nextUrl := "/v2/service_bindings"

	for nextUrl != "" {
		bindingsResp, err := self.getFromCfApi(nextUrl, cliConnection)
		if err != nil {
			return nil, fmt.Errorf("Unable to call service bindings endpoint: %s", err)
		}

		var bindings []pluginModels.ServiceBinding
		nextUrl, bindings, err = deserializeBindings(bindingsResp)
		if err != nil {
			return nil, fmt.Errorf("Unable to deserialize service bindings: %s", err)
		}

		allBindings = append(allBindings, bindings...)
	}

	return allBindings, nil
}

func deserializeBindings(bindingResponse []byte) (string, []pluginModels.ServiceBinding, error) {
	paginatedResources := new(resources.PaginatedServiceBindingResources)
	err := json.Unmarshal(bindingResponse, paginatedResources)

	if err != nil {
		return "", nil, fmt.Errorf("Unable to deserialize service bindings: %s", err)
	}

	bindings, err := paginatedResources.ToModel()
	return paginatedResources.NextUrl, bindings, err
}

func (self *apiClient) GetStartedApps(cliConnection plugin.CliConnection) ([]sdkModels.GetAppsModel, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve apps: %s", err)
	}

	startedApps := make([]sdkModels.GetAppsModel, 0, len(apps))

	for _, app := range apps {
		if app.State == "started" {
			startedApps = append(startedApps, app)
		}
	}

	return startedApps, nil
}
