#!/usr/bin/env bash

INTERNAL_TENANT_ID=$(docker exec -i ${POSTGRES_CONTAINER} psql -qtAX -U "${DB_USER}" -h "${DB_HOST}" -p "${DB_PORT}" -d "${DB_NAME}" -c "SELECT id FROM business_tenant_mappings WHERE external_tenant = '3e64ebae-38b5-46a0-b1ed-9ccee153a0ae'")
echo -e "${GREEN}Internal Tenant ID for default tenant from dump: $INTERNAL_TENANT_ID${NC}"


HEADER=$(echo "{ \"alg\": \"none\", \"typ\": \"JWT\" }" | base64 | tr '/+' '_-' | tr -d '=')
PAYLOAD=$(echo "{ \"scopes\": \"tenant:write fetch-request.auth:read webhooks.auth:read application.auths:read application.webhooks:read application_template.webhooks:read document.fetch_request:read event_spec.fetch_request:read api_spec.fetch_request:read runtime.auths:read integration_system.auths:read bundle.instance_auths:read bundle.instance_auths:read application:read automatic_scenario_assignment:write automatic_scenario_assignment:read health_checks:read application:write runtime:write label_definition:write label_definition:read runtime:read tenant:read formation:write\", \"tenant\":\"{\\\"consumerTenant\\\":\\\"$INTERNAL_TENANT_ID\\\",\\\"externalTenant\\\":\\\"3e64ebae-38b5-46a0-b1ed-9ccee153a0ae\\\"}\" }" | base64 | tr '/+' '_-' | tr -d '=')
JWT_TOKEN="$HEADER.$PAYLOAD."

echo -e "${GREEN}Use the following JWT token when requesting Director as default tenant:${NC}"
echo $JWT_TOKEN