package resource_mapper

import (
	"github.com/cloudfoundry-incubator/app-restarter/api"
	"github.com/cloudfoundry-incubator/app-restarter/models"
)

//go:generate counterfeiter . SpacesParser
type SpacesParser interface {
	Parse([]byte) (models.Spaces, error)
}

func Spaces(spacesParser SpacesParser, paginatedRequester PaginatedRequester) (models.Spaces, error) {
	var noSpaces models.Spaces

	filter := api.Filters{}

	params := map[string]interface{}{
		"inline-relations-depth": 1,
	}

	responseBodies, err := paginatedRequester.Do(filter, params)
	if err != nil {
		return noSpaces, err
	}

	var spaces models.Spaces

	for _, nextBody := range responseBodies {
		apps, err := spacesParser.Parse(nextBody)
		if err != nil {
			return noSpaces, err
		}

		spaces = append(spaces, apps...)
	}

	return spaces, nil
}
