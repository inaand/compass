package subject

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/cert"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	tenantEntity "github.com/kyma-incubator/compass/components/director/pkg/tenant"
	"github.com/pkg/errors"
)

const (
	RuntimeType           = "Runtime"
	IntegrationSystemType = "Integration System"
	ApplicationType       = "Application"
)

type subjectConsumerTypeMapping struct {
	Subject            string `json:"subject"`
	ConsumerType       string `json:"consumer_type"`
	InternalConsumerID string `json:"internal_consumer_id"`
	TenantAccessLevel  string `json:"tenant_access_level"`
}

func (s *subjectConsumerTypeMapping) validate() error {
	if len(s.Subject) < 1 {
		return errors.New("subject is not provided")
	}
	if s.ConsumerType != RuntimeType && s.ConsumerType != IntegrationSystemType && s.ConsumerType != ApplicationType {
		return fmt.Errorf("consumer type %s is not valid", s.ConsumerType)
	}
	if s.TenantAccessLevel != string(tenantEntity.Account) &&
		s.TenantAccessLevel != string(tenantEntity.Subaccount) &&
		s.TenantAccessLevel != string(tenantEntity.Customer) {
		return fmt.Errorf("tenant access level %s is not valid", s.TenantAccessLevel)
	}
	return nil
}

type processor struct {
	mappings  []subjectConsumerTypeMapping
	ouPattern string
}

// NewProcessor returns a new subject processor configured with the given subject-to-consumer mapping, and subject organization unit pattern.
// If the subject-to-consumer mapping is invalid, an error is returned.
func NewProcessor(subjectConsumerTypeMappingConfig string, ouPattern string) (*processor, error) {
	mappings, err := unmarshalMappings(subjectConsumerTypeMappingConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "while configuring subject processor")
	}
	for _, m := range mappings {
		if err := m.validate(); err != nil {
			return nil, err
		}
	}
	return &processor{
		mappings:  mappings,
		ouPattern: ouPattern,
	}, nil
}

// AuthIDFromSubjectFunc returns a function able to extract the authentication ID from a given certificate subject.
func (p *processor) AuthIDFromSubjectFunc() func(subject string) string {
	authIDFromMappingFunc := p.authIDFromMappings()
	authIDFromOUsFunc := cert.GetRemainingOrganizationalUnit(p.ouPattern)
	return func(subject string) string {
		if authIDFromMapping := authIDFromMappingFunc(subject); authIDFromMapping != "" {
			return authIDFromMapping
		}
		return authIDFromOUsFunc(subject)
	}
}

// AuthSessionExtraFromSubjectFunc returns a function which returns consumer-specific auth session extra body
// in case the subject matches any of the configured consumers from the mapping.
func (p *processor) AuthSessionExtraFromSubjectFunc() func(context.Context, string) map[string]interface{} {
	return func(ctx context.Context, subject string) map[string]interface{} {
		log.C(ctx).Infof("trying to extract auth session extra from subject %s", subject)
		for _, m := range p.mappings {
			log.C(ctx).Infof("trying to match subject pattern %s", m.Subject)
			if subjectsMatch(subject, m.Subject) {
				log.C(ctx).Infof("pattern matched subject!")
				return cert.GetExtra(m.ConsumerType, m.TenantAccessLevel, m.InternalConsumerID)
			}
		}

		return nil
	}
}

func (p *processor) authIDFromMappings() func(subject string) string {
	return func(subject string) string {
		for _, m := range p.mappings {
			if subjectsMatch(subject, m.Subject) {
				return m.InternalConsumerID
			}
		}
		return ""
	}
}

func unmarshalMappings(mappingsConfig string) ([]subjectConsumerTypeMapping, error) {
	var mappings []subjectConsumerTypeMapping
	if err := json.Unmarshal([]byte(mappingsConfig), &mappings); err != nil {
		return nil, errors.Wrap(err, "while unmarshalling mappings")
	}

	return mappings, nil
}

func subjectsMatch(actualSubject, expectedSubject string) bool {
	expectedSubjectComponents := strings.Split(expectedSubject, ",")

	for _, expectedSubjectComponent := range expectedSubjectComponents {
		if !strings.Contains(actualSubject, strings.TrimSpace(expectedSubjectComponent)) {
			return false
		}
	}

	return true
}