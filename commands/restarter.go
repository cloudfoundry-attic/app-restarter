package commands

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/app-restarter/api"
)

type AppRestarter interface {
	Restart(string) ([]string, error)
}

type appRestarter struct {
	cli api.Connection
}

func NewAppRestarter(cli api.Connection) AppRestarter {
	return &appRestarter{
		cli: cli,
	}
}

func (r *appRestarter) Restart(appGuid string) ([]string, error) {
	output, err := r.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"state":"STOPPED"}`)
	if err != nil {
		return output, err
	}

	if err = checkError(output[0]); err != nil {
		return output, err
	}

	output, err = r.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"state":"STARTED"}`)
	if err != nil {
		return output, err
	}

	if err = checkError(output[0]); err != nil {
		return output, err
	}

	return output, nil
}

type apiError struct {
	Code        int64  `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
}

func checkError(jsonRsp string) error {
	b := []byte(jsonRsp)
	theError := apiError{}
	err := json.Unmarshal(b, &theError)
	if err != nil {
		return err
	}

	if theError.ErrorCode != "" || theError.Code != 0 {
		return errors.New(theError.ErrorCode + " - " + theError.Description)
	}

	return nil
}
