package resource_mapper

import "github.com/cloudfoundry-incubator/app-restarter/api"

//go:generate counterfeiter . PaginatedRequester
type PaginatedRequester interface {
	Do(filter api.Filter, params map[string]interface{}) ([][]byte, error)
}
