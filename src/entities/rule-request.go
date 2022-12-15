package entities

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cjlapao/common-go/helper"
)

// RuleRequest struct
type RuleRequest struct {
	Name      string `json:"name"`
	SQLFilter string `json:"sqlFilter"`
	SQLAction string `json:"sqlAction"`
}

func (r *RuleRequest) IsValid() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if r.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Queue name is null"
		errorResponse.Message = "Queue name cannot be null"
		return false, &errorResponse
	}

	return true, nil
}

func (r *RuleRequest) FromFile(filePath string) error {
	fileExists := helper.FileExists(filePath)

	if !fileExists {
		err := errors.New("file " + filePath + " was not found")
		return err
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, r)
	if err != nil {
		return err
	}

	return nil
}
