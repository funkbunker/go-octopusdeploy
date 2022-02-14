package accounts

import (
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy/resources"

	"github.com/stretchr/testify/require"
)

func TestTokenAccountNew(t *testing.T) {
	accountType := AccountTypeToken
	description := ""
	environmentIDs := []string{}
	name := octopusdeploy.getRandomName()
	spaceID := ""
	tenantedDeploymentMode := resources.TenantedDeploymentMode("Untenanted")
	token := resources.NewSensitiveValue(octopusdeploy.getRandomName())

	account, err := NewTokenAccount(name, token)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	// Resource
	require.Equal(t, "", account.ID)
	require.Equal(t, "", account.ModifiedBy)
	require.Nil(t, account.ModifiedOn)
	require.NotNil(t, account.Links)

	// IResource
	require.Equal(t, "", account.GetID())
	require.Equal(t, "", account.GetModifiedBy())
	require.Nil(t, account.GetModifiedOn())
	require.NotNil(t, account.GetLinks())

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

	// TokenAccount
	require.Equal(t, token, account.Token)
}

func TestTokenAccountNewWithConfigs(t *testing.T) {
	environmentIDs := []string{octopusdeploy.getRandomName(), octopusdeploy.getRandomName()}
	invalidID := octopusdeploy.getRandomName()
	invalidModifiedBy := octopusdeploy.getRandomName()
	invalidModifiedOn := time.Now()
	invalidName := octopusdeploy.getRandomName()
	name := octopusdeploy.getRandomName()
	description := "Description for " + name + " (OK to Delete)"
	spaceID := octopusdeploy.getRandomName()
	tenantedDeploymentMode := resources.TenantedDeploymentMode("Tenanted")
	token := resources.NewSensitiveValue(octopusdeploy.getRandomName())

	options := func(a *TokenAccount) {
		a.Description = description
		a.EnvironmentIDs = environmentIDs
		a.ID = invalidID
		a.ModifiedBy = invalidModifiedBy
		a.ModifiedOn = &invalidModifiedOn
		a.Name = invalidName
		a.SpaceID = spaceID
		a.TenantedDeploymentMode = tenantedDeploymentMode
		a.Token = token
	}

	account, err := NewTokenAccount(name, token, options)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	// Resource
	require.Equal(t, "", account.ID)
	require.Equal(t, "", account.ModifiedBy)
	require.Nil(t, account.ModifiedOn)
	require.NotNil(t, account.Links)

	// IResource
	require.Equal(t, "", account.GetID())
	require.Equal(t, "", account.GetModifiedBy())
	require.Nil(t, account.GetModifiedOn())
	require.NotNil(t, account.GetLinks())

	// account
	require.Equal(t, description, account.Description)
	require.Equal(t, environmentIDs, account.EnvironmentIDs)
	require.Equal(t, name, account.Name)
	require.Equal(t, spaceID, account.SpaceID)
	require.Equal(t, tenantedDeploymentMode, account.TenantedDeploymentMode)

	// IAccount
	require.Equal(t, AccountTypeToken, account.GetAccountType())
	require.Equal(t, description, account.GetDescription())
	require.Equal(t, name, account.GetName())

	// TokenAccount
	require.Equal(t, token, account.Token)
}

func TestTokenAccountSetName(t *testing.T) {
	name := octopusdeploy.getRandomName()
	token := resources.NewSensitiveValue(octopusdeploy.getRandomName())

	account, err := NewTokenAccount(name, token)

	require.NotNil(t, account)
	require.NoError(t, err)
	require.NoError(t, account.Validate())

	require.Equal(t, name, account.Name)
	require.Equal(t, name, account.GetName())
}
