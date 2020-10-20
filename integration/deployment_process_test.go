package integration

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentProcessGet(t *testing.T) {
	octopusClient := getOctopusClient()
	require.NotNil(t, octopusClient)

	project := createTestProject(t, octopusClient, getRandomName())
	defer cleanProject(t, octopusClient, project.ID)

	deploymentProcess, err := octopusClient.DeploymentProcesses.GetByID(project.DeploymentProcessID)

	assert.Equal(t, project.DeploymentProcessID, deploymentProcess.ID)
	assert.NoError(t, err, "there should be error raised getting a projects deployment process")
}

func TestDeploymentProcessGetAll(t *testing.T) {
	octopusClient := getOctopusClient()
	require.NotNil(t, octopusClient)

	project := createTestProject(t, octopusClient, getRandomName())
	defer cleanProject(t, octopusClient, project.ID)

	allDeploymentProcess, err := octopusClient.DeploymentProcesses.GetAll()
	require.NoError(t, err)

	numberOfDeploymentProcesses := len(allDeploymentProcess)

	additionalProject := createTestProject(t, octopusClient, getRandomName())
	defer cleanProject(t, octopusClient, additionalProject.ID)

	allDeploymentProcessAfterCreatingAdditional, err := octopusClient.DeploymentProcesses.GetAll()
	require.NoError(t, err)

	assert.Equal(t, len(allDeploymentProcessAfterCreatingAdditional), numberOfDeploymentProcesses+1, "created an additional project and expected number of deployment processes to increase by 1")
}

func TestDeploymentProcessUpdate(t *testing.T) {
	octopusClient := getOctopusClient()

	project := createTestProject(t, octopusClient, getRandomName())
	defer cleanProject(t, octopusClient, project.ID)

	deploymentProcess, err := octopusClient.DeploymentProcesses.GetByID(project.DeploymentProcessID)

	if err != nil {
		t.Fatalf("Retrieving deployment processes failed when it shouldn't: %s", err)
	}

	deploymentActionWindowService := &octopusdeploy.DeploymentAction{
		Name:       "Install Windows Service",
		ActionType: "Octopus.WindowService",
		Properties: map[string]string{
			"Octopus.Action.WindowService.CreateOrUpdateService":                        "True",
			"Octopus.Action.WindowService.ServiceAccount":                               "LocalSystem",
			"Octopus.Action.WindowService.StartMode":                                    "auto",
			"Octopus.Action.Package.AutomaticallyRunConfigurationTransformationFiles":   "True",
			"Octopus.Action.Package.AutomaticallyUpdateAppSettingsAndConnectionStrings": "True",
			"Octopus.Action.EnabledFeatures":                                            "Octopus.Features.WindowService,Octopus.Features.ConfigurationVariables,Octopus.Features.ConfigurationTransforms,Octopus.Features.SubstituteInFiles",
			"Octopus.Action.Package.FeedId":                                             "feeds-nugetfeed",
			"Octopus.Action.Package.DownloadOnTentacle":                                 "False",
			"Octopus.Action.Package.PackageId":                                          "Newtonsoft.Json",
			"Octopus.Action.WindowService.ServiceName":                                  "My service name",
			"Octopus.Action.WindowService.DisplayName":                                  "my display name",
			"Octopus.Action.WindowService.Description":                                  "my desc",
			"Octopus.Action.WindowService.ExecutablePath":                               "bin\\Myservice.exe",
			"Octopus.Action.SubstituteInFiles.Enabled":                                  "True",
			"Octopus.Action.SubstituteInFiles.TargetFiles":                              "*.sh",
		},
	}

	step1 := &octopusdeploy.DeploymentStep{
		Name: "My First Step",
		Properties: map[string]string{
			"Octopus.Action.TargetRoles": "octopus-server",
		},
	}

	step1.Actions = append(step1.Actions, *deploymentActionWindowService)

	deploymentProcess.Steps = append(deploymentProcess.Steps, *step1)

	updated, err := octopusClient.DeploymentProcesses.Update(*deploymentProcess)

	assert.NoError(t, err, "error when updating deployment process")
	assert.Equal(t, updated.Steps[0].Properties, deploymentProcess.Steps[0].Properties)
	assert.Equal(t, updated.Steps[0].Actions[0].ActionType, deploymentProcess.Steps[0].Actions[0].ActionType)
}
