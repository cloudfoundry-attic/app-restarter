package models

import "encoding/json"

type Applications []Application

type ApplicationMetadata struct {
	Guid string `json:"guid"`
}

const (
	Started = "STARTED"
	Stopped = "STOPPED"
)

type ApplicationEntity struct {
	Name      string `json:"name"`
	Diego     bool
	State     string `json:"state"`
	SpaceGuid string `json:"space_guid"`
}

type ApplicationsResponse struct {
	Resources Applications `json:"resources"`
}

type Application struct {
	ApplicationEntity   `json:"entity"`
	ApplicationMetadata `json:"metadata"`
}

type ApplicationsParser struct{}

func (a ApplicationsParser) Parse(body []byte) (Applications, error) {
	var response ApplicationsResponse
	var emptyApplications Applications

	err := json.Unmarshal(body, &response)
	if err != nil {
		return emptyApplications, err
	}

	return response.Resources, nil
}
