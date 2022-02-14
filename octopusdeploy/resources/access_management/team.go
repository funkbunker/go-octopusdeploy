package access_management

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy/services"
	"github.com/go-playground/validator/v10"
)

// Teams defines a collection of teams with built-in support for paged results.
type Teams struct {
	Items []*Team `json:"Items"`
	services.PagedResults
}

type Team struct {
	CanBeDeleted           bool                          `json:"CanBeDeleted,omitempty"`
	CanBeRenamed           bool                          `json:"CanBeRenamed,omitempty"`
	CanChangeMembers       bool                          `json:"CanChangeMembers,omitempty"`
	CanChangeRoles         bool                          `json:"CanChangeRoles,omitempty"`
	Description            string                        `json:"Description,omitempty"`
	ExternalSecurityGroups []services.NamedReferenceItem `json:"ExternalSecurityGroups,omitempty"`
	MemberUserIDs          []string                      `json:"MemberUserIds"`
	Name                   string                        `json:"Name" validate:"required"`
	SpaceID                string                        `json:"SpaceId,omitempty"`

	services.Resource
}

func NewTeam(name string) *Team {
	return &Team{
		Name:     name,
		resource: *services.NewResource(),
	}
}

// Validate checks the state of the team and returns an error if invalid.
func (t *Team) Validate() error {
	return validator.New().Struct(t)
}
