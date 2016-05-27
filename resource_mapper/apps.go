package resource_mapper

import (
	"fmt"
	"github.com/cloudfoundry-incubator/app-restarter/api"
	"github.com/cloudfoundry-incubator/app-restarter/models"
)

type AppsGetterFunc func(ApplicationsParser, PaginatedRequester) (models.Applications, error)

type ApplicationsParser interface {
	Parse([]byte) (models.Applications, error)
}

type AppsGetter struct {
	OrganizationGuid string
	SpaceGuid        string
}

type OrgNotFoundErr struct {
	OrganizationName string
}

func (e OrgNotFoundErr) Error() string {
	return fmt.Sprintf("Organization not found: %s", e.OrganizationName)
}

type SpaceNotFoundErr struct {
	SpaceName string
}

func (e SpaceNotFoundErr) Error() string {
	return fmt.Sprintf("Space not found: %s", e.SpaceName)
}

func NewAppsGetterFunc(
	cliConnection api.Connection,
	orgName string,
	spaceName string,
) (AppsGetterFunc, error) {
	command := AppsGetter{}

	if orgName != "" {
		org, err := cliConnection.GetOrg(orgName)
		if err != nil || org.Guid == "" {
			return nil, OrgNotFoundErr{OrganizationName: orgName}
		}
		command.OrganizationGuid = org.Guid
	} else if spaceName != "" {
		space, err := cliConnection.GetSpace(spaceName)
		if err != nil || space.Guid == "" {
			return nil, SpaceNotFoundErr{SpaceName: spaceName}
		}
		command.SpaceGuid = space.Guid
	}

	var appsGetterFunc = command.Apps

	return appsGetterFunc, nil
}

func (c AppsGetter) Apps(
	appsParser ApplicationsParser,
	paginatedRequester PaginatedRequester,
) (models.Applications, error) {
	var noApps models.Applications

	filter := api.Filters{}

	if c.OrganizationGuid != "" {
		filter = append(
			filter,
			api.EqualFilter{
				Name:  "organization_guid",
				Value: c.OrganizationGuid,
			},
		)
	} else if c.SpaceGuid != "" {
		filter = append(
			filter,
			api.EqualFilter{
				Name:  "space_guid",
				Value: c.SpaceGuid,
			},
		)
	}

	params := map[string]interface{}{}

	responseBodies, err := paginatedRequester.Do(filter, params)
	if err != nil {
		return noApps, err
	}

	var applications models.Applications

	for _, nextBody := range responseBodies {
		apps, err := appsParser.Parse(nextBody)
		if err != nil {
			return noApps, err
		}

		applications = append(applications, apps...)
	}

	return applications, nil
}
