package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/dghubble/sling"
)

// Client is an OctopusDeploy for making OctpusDeploy API requests.
type Client struct {
	sling               *sling.Sling
	Accounts            *AccountService
	ActionTemplates     *ActionTemplateService
	APIKeys             *APIKeyService
	Authentication      *AuthenticationService
	Certificates        *CertificateService
	Channels            *ChannelService
	Configuration       *ConfigurationService
	DeploymentProcesses *DeploymentProcessService
	Environments        *EnvironmentService
	Feeds               *FeedService
	Interruptions       *InterruptionsService
	LibraryVariableSets *LibraryVariableSetService
	Lifecycles          *LifecycleService
	Machines            *MachineService
	MachinePolicies     *MachinePolicyService
	Projects            *ProjectService
	ProjectGroups       *ProjectGroupService
	ProjectTriggers     *ProjectTriggerService
	Root                *RootService
	Spaces              *SpaceService
	TagSets             *TagSetService
	Tenants             *TenantService
	Users               *UserService
	Variables           *VariableService
}

// NewClient returns a new
func NewClient(httpClient *http.Client, octopusURL string, apiKey string, spaceName string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if isEmpty(octopusURL) {
		return nil, createInvalidParameterError("NewClient", "octopusURL")
	}

	if isEmpty(apiKey) {
		return nil, createInvalidParameterError("NewClient", "apiKey")
	}

	baseURLWithAPI := strings.TrimRight(octopusURL, "/")

	if isEmpty(spaceName) {
		baseURLWithAPI = fmt.Sprintf("%s/api/", baseURLWithAPI)
	} else {
		baseURLWithAPI = fmt.Sprintf("%s/api/%s/", baseURLWithAPI, spaceName)
	}

	base := sling.New().Client(httpClient).Base(baseURLWithAPI).Set("X-Octopus-ApiKey", apiKey)

	return &Client{
		sling:               base,
		Accounts:            NewAccountService(base.New()),
		ActionTemplates:     NewActionTemplateService(base.New()),
		APIKeys:             NewAPIKeyService(base.New()),
		Authentication:      NewAuthenticationService(base.New()),
		Certificates:        NewCertificateService(base.New()),
		Channels:            NewChannelService(base.New()),
		Configuration:       NewConfigurationService(base.New()),
		DeploymentProcesses: NewDeploymentProcessService(base.New()),
		Environments:        NewEnvironmentService(base.New()),
		Feeds:               NewFeedService(base.New()),
		Interruptions:       NewInterruptionsService(base.New()),
		Machines:            NewMachineService(base.New()),
		MachinePolicies:     NewMachinePolicyService(base.New()),
		LibraryVariableSets: NewLibraryVariableSetService(base.New()),
		Lifecycles:          NewLifecycleService(base.New()),
		Projects:            NewProjectService(base.New()),
		ProjectGroups:       NewProjectGroupService(base.New()),
		ProjectTriggers:     NewProjectTriggerService(base.New()),
		Root:                NewRootService(base.New()),
		Spaces:              NewSpaceService(base.New()),
		TagSets:             NewTagSetService(base.New()),
		Tenants:             NewTenantService(base.New()),
		Users:               NewUserService(base.New()),
		Variables:           NewVariableService(base.New()),
	}, nil
}

func ForSpace(httpClient *http.Client, octopusURL string, apiKey string, space *model.Space) (*Client, error) {
	if httpClient == nil {
		return nil, createInvalidParameterError("ForSpace", "httpClient")
	}

	if isEmpty(octopusURL) {
		return nil, createInvalidParameterError("ForSpace", "octopusURL")
	}

	if isEmpty(apiKey) {
		return nil, createInvalidParameterError("ForSpace", "apiKey")
	}

	if space == nil {
		return nil, createInvalidParameterError("ForSpace", "space")
	}

	baseURLWithAPI := strings.TrimRight(octopusURL, "/")

	const apiPath = "%s/api/%s/"
	baseURLWithAPI = fmt.Sprintf(apiPath, baseURLWithAPI, space.ID)

	base := sling.New().Client(httpClient).Base(baseURLWithAPI).Set("X-Octopus-ApiKey", apiKey)

	return &Client{
		sling:               base,
		Accounts:            NewAccountService(base.New()),
		ActionTemplates:     NewActionTemplateService(base.New()),
		APIKeys:             NewAPIKeyService(base.New()),
		Authentication:      NewAuthenticationService(base.New()),
		Certificates:        NewCertificateService(base.New()),
		Channels:            NewChannelService(base.New()),
		Configuration:       NewConfigurationService(base.New()),
		DeploymentProcesses: NewDeploymentProcessService(base.New()),
		Environments:        NewEnvironmentService(base.New()),
		Feeds:               NewFeedService(base.New()),
		LibraryVariableSets: NewLibraryVariableSetService(base.New()),
		Lifecycles:          NewLifecycleService(base.New()),
		Machines:            NewMachineService(base.New()),
		MachinePolicies:     NewMachinePolicyService(base.New()),
		Projects:            NewProjectService(base.New()),
		ProjectGroups:       NewProjectGroupService(base.New()),
		ProjectTriggers:     NewProjectTriggerService(base.New()),
		Root:                NewRootService(base.New()),
		TagSets:             NewTagSetService(base.New()),
		Tenants:             NewTenantService(base.New()),
		Users:               NewUserService(base.New()),
		Variables:           NewVariableService(base.New()),
	}, nil
}

type APIError struct {
	ErrorMessage  string   `json:"ErrorMessage"`
	Errors        []string `json:"Errors"`
	FullException string   `json:"FullException"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Octopus Deploy Error Response: %v %+v %v", e.ErrorMessage, e.Errors, e.FullException)
}

// APIErrorChecker is a generic error handler for the OctopusDeploy API.
func APIErrorChecker(urlPath string, resp *http.Response, wantedResponseCode int, slingError error, octopusDeployError *APIError) error {
	if octopusDeployError.Errors != nil {
		return fmt.Errorf("octopus deploy api returned an error on endpoint %s - %s", urlPath, octopusDeployError.Errors)
	}

	if slingError != nil {
		return fmt.Errorf("cannot get endpoint %s from server. failure from http client %v", urlPath, slingError)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrItemNotFound
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	if resp.StatusCode != wantedResponseCode {
		return fmt.Errorf("cannot get item from endpoint %s. response from server %s", urlPath, resp.Status)
	}

	return nil
}

// LoadNextPage checks if the next page should be loaded from the API. Returns the new path and a bool if the next page should be checked.
func LoadNextPage(pagedResults model.PagedResults) (string, bool) {
	if pagedResults.Links.PageNext != "" {
		return pagedResults.Links.PageNext, true
	}

	return "", false
}

// Generic OctopusDeploy API Get Function.
func apiGet(sling *sling.Sling, inputStruct interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError("apiGet", "sling")
	}

	getClient := sling.New()

	if getClient == nil {
		return nil, createClientInitializationError("apiGet")
	}

	getClient = getClient.Get(path)

	if getClient == nil {
		return nil, createClientInitializationError("apiGet")
	}

	octopusDeployError := new(APIError)
	resp, err := getClient.Receive(inputStruct, &octopusDeployError)

	if err != nil {
		return nil, err
	}

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return inputStruct, nil
}

// Generic OctopusDeploy API Add Function. Expects a 201 response.
func apiAdd(sling *sling.Sling, inputStruct, resource model.ResourceInterface, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError("apiAdd", "sling")
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError("apiAdd", "path")
	}

	postClient := sling.New()

	if postClient == nil {
		return nil, createClientInitializationError("apiAdd")
	}

	postClient = postClient.Post(path)

	if postClient == nil {
		return nil, createClientInitializationError("apiAdd")
	}

	request := postClient.BodyJSON(inputStruct)

	if request == nil {
		return nil, createClientInitializationError("apiAdd")
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(resource, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusCreated, err, octopusDeployError)

	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return resource, nil
}

// apiPost post to octopus and expect a 200 response code.
func apiPost(sling *sling.Sling, inputStruct, returnStruct interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError("apiPost", "sling")
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError("apiPost", "path")
	}

	postClient := sling.New()

	if postClient == nil {
		return nil, createClientInitializationError("apiPost")
	}

	postClient = postClient.Post(path)

	if postClient == nil {
		return nil, createClientInitializationError("apiPost")
	}

	request := postClient.BodyJSON(inputStruct)

	if request == nil {
		return nil, createClientInitializationError("apiPost")
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(returnStruct, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return returnStruct, nil
}

// Generic OctopusDeploy API Update Function.
func apiUpdate(sling *sling.Sling, inputStruct, returnStruct interface{}, path string) (interface{}, error) {
	if sling == nil {
		return nil, createInvalidParameterError("apiUpdate", "sling")
	}

	if isEmpty(path) {
		return nil, createInvalidParameterError("apiUpdate", "path")
	}

	putClient := sling.New()

	if putClient == nil {
		return nil, createClientInitializationError("apiUpdate")
	}

	putClient = putClient.Put(path)

	if putClient == nil {
		return nil, createClientInitializationError("apiUpdate")
	}

	request := putClient.BodyJSON(inputStruct)

	if request == nil {
		return nil, createClientInitializationError("apiUpdate")
	}

	octopusDeployError := new(APIError)
	resp, err := request.Receive(returnStruct, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

	if apiErrorCheck != nil {
		return nil, apiErrorCheck
	}

	return returnStruct, nil
}

// Generic OctopusDeploy API Delete Function.
func apiDelete(sling *sling.Sling, path string) error {
	if sling == nil {
		return createInvalidParameterError("apiDelete", "sling")
	}

	if isEmpty(path) {
		return createInvalidParameterError("apiDelete", "path")
	}

	deleteClient := sling.New()

	if deleteClient == nil {
		return createClientInitializationError("apiDelete")
	}

	deleteClient = deleteClient.Delete(path)

	if deleteClient == nil {
		return createClientInitializationError("apiDelete")
	}

	octopusDeployError := new(APIError)
	resp, err := deleteClient.Receive(nil, &octopusDeployError)

	apiErrorCheck := APIErrorChecker(path, resp, http.StatusOK, err, octopusDeployError)

	if apiErrorCheck != nil {
		return apiErrorCheck
	}

	return nil
}

// ErrItemNotFound is an OctopusDeploy error returned an item cannot be found.
var ErrItemNotFound = errors.New("cannot find the item")
