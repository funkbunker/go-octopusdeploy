package accounts

import (
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy/resources"

	"github.com/stretchr/testify/require"
)

func TestAmazonWebServicesAccountNew(t *testing.T) {
	accessKey := octopusdeploy.getRandomName()
	accountType := AccountTypeAmazonWebServicesAccount
	description := resources.emptyString
	environmentIDs := []string{}
	name := octopusdeploy.getRandomName()
	secretKey := resources.NewSensitiveValue(octopusdeploy.getRandomName())
	spaceID := resources.emptyString
	tenantedDeploymentMode := resources.TenantedDeploymentMode("Untenanted")

	account, err := NewAmazonWebServicesAccount(name, accessKey, secretKey)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	// account
	require.Equal(t, description, account.Description)
	require.Equal(t, environmentIDs, account.EnvironmentIDs)
	require.Equal(t, name, account.Name)
	require.Equal(t, spaceID, account.SpaceID)
	require.Equal(t, tenantedDeploymentMode, account.TenantedDeploymentMode)

	// IAccount
	require.Equal(t, accountType, account.GetAccountType())
	require.Equal(t, description, account.GetDescription())
	require.Equal(t, name, account.GetName())

	// AmazonWebServicesAccount
	require.Equal(t, accessKey, account.AccessKey)
	require.Equal(t, secretKey, account.SecretKey)
}

func TestAmazonWebServicesAccountNewWithConfigs(t *testing.T) {
	accessKey := octopusdeploy.getRandomName()
	accountType := AccountTypeAmazonWebServicesAccount
	environmentIDs := []string{"environment-id-1", "environment-id-2"}
	invalidID := octopusdeploy.getRandomName()
	invalidModifiedBy := octopusdeploy.getRandomName()
	invalidModifiedOn := time.Now()
	invalidName := octopusdeploy.getRandomName()
	name := octopusdeploy.getRandomName()
	description := "Description for " + name + " (OK to Delete)"
	secretKey := resources.NewSensitiveValue(octopusdeploy.getRandomName())
	spaceID := octopusdeploy.getRandomName()
	tenantedDeploymentMode := resources.TenantedDeploymentMode("Tenanted")

	options := func(a *AmazonWebServicesAccount) {
		a.Description = description
		a.EnvironmentIDs = environmentIDs
		a.ID = invalidID
		a.Name = invalidName
		a.SecretKey = secretKey
		a.SpaceID = spaceID
		a.TenantedDeploymentMode = tenantedDeploymentMode
	}

	account, err := NewAmazonWebServicesAccount(name, accessKey, secretKey, options)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NoError(t, account.Validate())

	// account
	require.Equal(t, description, account.Description)
	require.Equal(t, environmentIDs, account.EnvironmentIDs)
	require.Equal(t, name, account.Name)
	require.Equal(t, spaceID, account.SpaceID)
	require.Equal(t, tenantedDeploymentMode, account.TenantedDeploymentMode)

	// IAccount
	require.Equal(t, accountType, account.GetAccountType())
	require.Equal(t, description, account.GetDescription())
	require.Equal(t, name, account.GetName())

	// AmazonWebServicesAccount
	require.Equal(t, accessKey, account.AccessKey)
	require.Equal(t, secretKey, account.SecretKey)
}
