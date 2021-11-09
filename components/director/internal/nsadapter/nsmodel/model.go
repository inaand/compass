package nsmodel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
)

type System struct {
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
	SystemType   string `json:"type"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	SystemNumber string `json:"systemNumber"`
}

func (s System) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Protocol, validation.Required),
		validation.Field(&s.Host, validation.Required),
		validation.Field(&s.SystemType, validation.Required),
		validation.Field(&s.Description, validation.Required),
		validation.Field(&s.Status, validation.Required),
	)
}

type SCC struct {
	Subaccount     string   `json:"subaccount"`
	LocationID     string   `json:"locationID"`
	ExposedSystems []System `json:"exposedSystems"`
}

func (s SCC) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Subaccount, validation.Required),
		validation.Field(&s.LocationID, validation.Required),
		validation.Field(&s.ExposedSystems, validation.Required),
		validation.Field(&s.ExposedSystems, validation.By(validateSystems)),
	)
}

func validateSystems(value interface{}) error {
	if systems, ok := value.([]System); ok {
		for _, s := range systems {
			if err := s.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

type Report struct {
	ReportType string `json:"type"`
	Value      []SCC  `json:"value"`
}

func (r Report) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ReportType, validation.Required),
		validation.Field(&r.Value, validation.Required),
		validation.Field(&r.Value, validation.By(validateSCCs)),
	)
}

func validateSCCs(value interface{}) error {
	if scc, ok := value.([]SCC); ok {
		for _, s := range scc {
			if err := s.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func ToAppRegisterInput(system System, tenant string, locationID string) model.ApplicationRegisterInput {
	return model.ApplicationRegisterInput{
		Name:         "",
		ProviderName: str.Ptr("SAP"),
		Labels:       map[string]interface{}{"SCC": locationID + system.Host},
		SystemNumber: str.Ptr(system.SystemNumber),
		Status:       system.Status,
	}
}

func ToAppUpdateInput(system System, tenant string, locationID string) model.ApplicationUpdateInput {
	// we should update the description, system type, protocol and system status
	// Where we will store the system type and the protocol ???
	return model.ApplicationUpdateInput{
		Description:         str.Ptr(system.Description),
		IntegrationSystemID: nil,
		StatusCondition:     nil,
		BaseURL:             nil,
		Labels:              nil,
		Status:              system.Status,
	}
}
