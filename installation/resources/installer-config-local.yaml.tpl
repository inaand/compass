apiVersion: v1
kind: Secret
metadata:
  name: application-connector-certificate-overrides
  namespace: compass-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
type: Opaque
data:
  global.applicationConnectorCa: ""
  global.applicationConnectorCaKey: ""
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-certificate-overrides
  namespace: compass-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.ingress.tlsCrt: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU4RENDQXRpZ0F3SUJBZ0lKQUpBKzlrQzZZZnZlTUEwR0NTcUdTSWIzRFFFQkN3VUFNQmN4RlRBVEJnTlYKQkFNTURDb3VhM2x0WVM1c2IyTmhiREFlRncweE9EQTNNVGd3T0RNNE1UTmFGdzB5T0RBM01UVXdPRE00TVROYQpNQmN4RlRBVEJnTlZCQU1NRENvdWEzbHRZUzVzYjJOaGJEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQCkFEQ0NBZ29DZ2dJQkFPZ05XbVROOU1ranlrUjdvQ0JGSFVqL01EcWh5bml3NEJITGo4ZTdDTFV5dFdwVHZXTkoKU1FiaFZDK3c5NkhHbU50MHZGUTR4OExUa3NNUmorcVZrdkcwKzBDTE1WQm14UjBVdnpZR09QRVRKREtsNTZkcQppMEM1Y2dnU3dkNjcveWxRZmQxc3FEVVBHM0pVZlNOSFFRWSs1SkUwaUpjZXhCQ3cvS200UXlUQXB0aEwzdEgxCm1ZRFBBQ0hUdVpsbGc2RVN4Z2RKY20rVmg5UkRvbWlON250ZjBZVG1xV1VIZXhZUkUrcTFGY0VhUHg5L2Q5QUwKZWd1WXZHTkduMVA4K1F3ZE5DKzVaMEVGa1EyS3RtalRBR09La2xJUG5NQkZubXhGSEtNZzNTa3RkYzYzSU5TNwpjSDM3dDNjb2l4N09HcDllcnRvSFZ5K2U5KzdiYnI3Z0lhWHNORVZUcjRGWXJIMlNPeHI1MjVRQUpHTys1dllQCkNYSlJJdXFNZWU3ZTg3aDRIU3JFVVdWTnhOdElBR1ZNYlRtV2F6aDJsSEFoazJjZVIzUkRaQzB3NGFOd3NRU1UKZE1yMmFCaUFobmFNSG40YnZIZnd1OHFRcjF5aVRUNG00SkltRkZTNDN2bGNETDZmVUJnNnF0U2p6QjF0WTFjTwp2Mnl1QXQzR1lKSENVMFd4R1ZEQk80T0ZncEc3aE5SRHoyRzVFMXQvRWU2VDUxdys2K295eXhib01pK09kVWdNCnBZOGlzcDdBTWhtdlZFOUdtNUhwQkpBRjYvcmE1WUNnNDh1SWtybkdZbm85eVNKcGhLU0JsMFRML01PNVFRaUoKb0hFNnV5ZkJXQjlZZ1NBdDl0MjJqU3FSWkpSdEtrbHB6bFJIRkROWTRwZlRraXlPSlFXYlM4RXJBZ01CQUFHagpQekE5TUFrR0ExVWRFd1FDTUFBd0N3WURWUjBQQkFRREFnWGdNQ01HQTFVZEVRUWNNQnFDQ210NWJXRXViRzlqCllXeUNEQ291YTNsdFlTNXNiMk5oYkRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQWdFQUlYYTlwenlyQ2dzMTRTOHUKZVFZdkorNEFzUE9uT1RGcExkaVl5UkVyNXdyNmJuMXUvMjZxc2FKckpxbkkyRk16SmdEQVRwZEtjbXRHYjBUOQp3S2wrYUJHcFFKcThrUWJwakVGTHhaWDJzaUNrRG82WittaUcrRjRKMHpKa3BKK0JHMS92eGZKbk0zK1ptdXQ5Ck9RV2ZjYTN3UHlhTWRDbGIyZjQwYlRFaFo5Mk9kcWlQMzFMbDlHWExSZmhaNTNsUzF0QWdvUGZoR25NbFY4b2MKWmxuSUROK25wS0Nma2tXUDJZUjlRLy9pa01tM01YRm9RSFppaVJseVZHSGFKZWRLMmNOQzlUYk4xNDFTaWZHZQo3V2FsQVBNcWNOQ3F3YStnN2RFSmR3ZjlRMklJTml0SjlDUVprT1dUZElYY2VHK2lZWWUrQXpmK1NkaHBocVdPCllFcDF6ek40dXI5U2VxU3NSaU9WY0RzVUFSa1M0clgrb0Vzb2hHL1Q5OTcrSDhjR2gzczl6TE84emtwRXZKSmEKS05QT0N5ODhVeEFOV2RRejFLMXRKVVQ2c3hkd0FEcXRJQnNPemhYVjlybDRRNStlZExlcmZPcUtCbUFRMUY5Swo2L1l0ZlNyY0JpeXZEU24wdFJ3OHJLRFVQU1hFNDFldXArOURNeThLVGl6T0RPTXVMSnR2dkJrTEFpNGNYQjVBCjQxMjBEdHdZQXNyNzNZYVl2SW8rWjV2OGZ4TjF3M3IwYS9KOVhZQlg3S3p1OFl4MnNUNWtWM2dNTHFCTXBaa3gKY29FTjNSandDMmV4VHl6dGc1ak1ZN2U4VFJ4OFFTeUxkK0pBd2t1Tm01NlNkcHFHNTE3cktJYkVMNDZzbkd0UgpCYUVOK01GeXNqdDU3ejhKQXJDMzhBMFN5dTQ9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
  global.ingress.tlsKey: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpRd0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1Mwd2dna3BBZ0VBQW9JQ0FRRG9EVnBremZUSkk4cEUKZTZBZ1JSMUkvekE2b2NwNHNPQVJ5NC9IdXdpMU1yVnFVNzFqU1VrRzRWUXZzUGVoeHBqYmRMeFVPTWZDMDVMRApFWS9xbFpMeHRQdEFpekZRWnNVZEZMODJCamp4RXlReXBlZW5hb3RBdVhJSUVzSGV1LzhwVUgzZGJLZzFEeHR5ClZIMGpSMEVHUHVTUk5JaVhIc1FRc1B5cHVFTWt3S2JZUzk3UjlabUF6d0FoMDdtWlpZT2hFc1lIU1hKdmxZZlUKUTZKb2plNTdYOUdFNXFsbEIzc1dFUlBxdFJYQkdqOGZmM2ZRQzNvTG1MeGpScDlUL1BrTUhUUXZ1V2RCQlpFTgppclpvMHdCamlwSlNENXpBUlo1c1JSeWpJTjBwTFhYT3R5RFV1M0I5KzdkM0tJc2V6aHFmWHE3YUIxY3ZudmZ1CjIyNis0Q0dsN0RSRlU2K0JXS3g5a2pzYStkdVVBQ1JqdnViMkR3bHlVU0xxakhudTN2TzRlQjBxeEZGbFRjVGIKU0FCbFRHMDVsbXM0ZHBSd0laTm5Ia2QwUTJRdE1PR2pjTEVFbEhUSzltZ1lnSVoyakI1K0c3eDM4THZLa0s5YwpvazArSnVDU0poUlV1Tjc1WEF5K24xQVlPcXJVbzh3ZGJXTlhEcjlzcmdMZHhtQ1J3bE5Gc1JsUXdUdURoWUtSCnU0VFVRODlodVJOYmZ4SHVrK2RjUHV2cU1zc1c2REl2am5WSURLV1BJcktld0RJWnIxUlBScHVSNlFTUUJldjYKMnVXQW9PUExpSks1eG1KNlBja2lhWVNrZ1pkRXkvekR1VUVJaWFCeE9yc253VmdmV0lFZ0xmYmR0bzBxa1dTVQpiU3BKYWM1VVJ4UXpXT0tYMDVJc2ppVUZtMHZCS3dJREFRQUJBb0lDQUZDN3ZKUlB0M2QzVlRyb1MvaU9NemNmCldaYzhqT1hhbThwMUtRdlRQWjlWQ2hyNUVXNEdwRHFaa0tHYkR6eWdqTFBsZEZSVkFPTCtteFAwK3o0aFZlTjAKRk9vS3cxaDJ1T042UVdBNVgvdzNyYU5WWnpndThFM1BkeVhwNkx0bWFzcmo3elpuUkVwWmZESVZ4UWZPRllobgp2enZwckEvdnEwVW5YbkJwNUNwWVFIUUdTWHFBMlN3Z1dLcHNNQ2wzVVFsc0w2dC9XU29MT3h1VmdGNmg2clBQCnpXUlFuK1MvYW9wdDNLRU82WWVxYXdXNVltVG1hVXE1aytseU82S0w0OVhjSHpqdlowWU8rcjFjWWtRc0RQbVUKejMxdll4amQzOVZKWWtJNi85Y0Fzdmo5YTVXM3ROYVFDZStTRW56Z05oRDJieHo1NnRKdG0xTGwweXpqYTdEVgp0OW9aQVd6WURMOTVnMDFZWVlQMG51c3A4WWlNOXd0OUtLc2NkcUd3RFN5WHNrWFBrcWFLVjNVcHJYdmhFbFVaCkErMmtjcm9VaDNGbEV4bGpwUmtKVkhOZXJ3NHRLRi9oYTFWRjZPdE10eTVQcXV0N0dGQmIvamtWeUg5cnpueWUKTXQyTWVyTTVPazMwd1NuTThISUdTUXpxYlplekJEZlNaUzRzcWdZZnBIMlhtMEs0SjgrRUowQ2hhMXZVSmVNMAoyZ284d00vaHljdmtqTEgxSmM3OEhpaVBTQ01udkpHemUxc2tWdmtRRFhBSFdldzBTUHpUSTZHYjZCb0Y4aVNECm0wZjR2azNoV3NlUWZBaXVZSnlUeUZXNmRhOGE1K2lpSDN4cVRsUUN1MDN1Nmo0U0l4aThJZlNmd0YwQTBldVAKNGtzalZTZVZyT3ExUnlvNUtpR3hBb0lCQVFEOWZtYnl6aW9QdVhRYzl0QXBxMUpSMzErQzlCdFFzcDg5WkZkSQpQaU5xaTJ3NVlVcTA0OFM2Z3VBb3JGOHNObUI3QjhWa1JlclowQ3hub2NHY0tleWdTYWsvME5qVElndk5weGJwCnBGbkFnRjlmbW1oTEl2SlF2REo2Q0ZidDRCQlZIdkJEWlYyQnZqK0k3NUxkK01jN2RPVDdFek1FRjBXcUdzY2MKTUpyNjRXQi9UMkF5dWR1YXlRT2NobmJFQ25FUmdRcHFlbG54MVBraytqbGNvYUs5QjFYUStVOUgzOHppM0FYNApENUxMY0Nhem9YYWlvS0swckNlNE5Ga2hOVXd0TFV0QXhSTXk2aksyZUZudWczUFRlY3N1WktNMElITktqZ0dCCnpGanZVb2tMcFVFb3BJa3FHM09yc0xmanpHZW9jaVFPUXNEdzlUb3lXL0FSOFhmWkFvSUJBUURxV0s2TThVN3EKUXJPeTYzNnpEZlBaZ2ljeXlsUWVoOFdMclBlbW5NeWdQYzR4eWoybnMrTmVRSnNEUmtPT2tWY250SEFYaTcwWgoyT0NCV3dwZHJuTXlSc3RIMU05bjdFNU5TVHZlZDlkU3YxUzRBb3NzS1hDSmgyUHBjYjV0OE9nL3ZGTlNYUlUyClk2aUorWTdOcDBZNDNxSlJOVnlRemd3YmFzaEpiUVdkVFFoVEVhdEVRS2JsUlZSblhlQWRjOXlhNUpHbkRpaTkKbFQrRWEzdFpvN1dha05oeHJkYjVuTkZ3a0xoNEs3ZkFtT3ZzMVBMQWx0SUZqeURCeDEvY1ZHblpDUDBVQmJqZgpkU2FueXBBdVRuMzd1VUwrcXpPVDlYWEZENllGT0x0NWV4d3RxdnUwSzZCNjcvajFFTlJDRk45RnMzWlV5RFFXClZUaDcybFhWU3NLakFvSUJBUUR6ZE5pdXpTNDhWK0thaHJpNXJGNnRYeGkrRG0vRmV5ZlFzSFBiWUVKbmEyd1AKVjgrR0YxS3p4a28vQmYySjJ0ZWlrWDRVcGNtK1UxNnlVUG8vWDB4eFRRMk55cWpUYmRsa005dWZuVWJOeVB6UQpOdDEvZkJxNVMyWTNLWmREY25SOUsrK1k2dHQ1Wmh4akNhUkdKMDVCWGkwa3JmWExNZ2FvTG51WUtWNVBJUEdxCms3TlNSSW9UQ0llOVpxN2Q3U0ZXckZZeW1UdVZOUFByZlo1bHhwOGphTTRVbTd4MnpReGJ2UERHb3o1YXdHV0wKRThGNncwaEF1UzZValVJazBLbE9vamVxQnh3L1JBcGNrUTNlTXNXbEQwNENTb2tyNFJhWlBmVllrY2ZBWWNaWgpOdWR6ZjBKMC9GU0ZTbjN4L0RoNTROV2NGS1IxUnpBVGVaVUJ4cVZSQW9JQkFRRFlDNmZvVWpOQnJ2ckNFVzkrCkhYZlk1Ni9CbUZ4U3hUTHU0U2h6VncwakViZTltVWljQ2pDc1hQMUwySVJCdEdaWU9YWTVqdDlvSzlSV0RSdVMKWUZqZFdmemduU1lWRmZyZUw0emRQVGlxbGEvQjhMNWptVlNoeGNycmxheE02Uk1FWjFlZGtDa1ZPbTFQdmwzVAo1TW5OZGhySXFWeE1OMWxjRVdiU29vclJpUW9Lb3poMHRQSG9YckZBbG9BZVJ3bHpWeE9jb21ZVzJiaDBHUzdmCjVoaHZoZWUxYmVISnY3UXFoWkU3WUhxSU9iTVBaUWJqWEdnRkxmMnlDRitzM2Jtem1DRFJTN0V6ZVdxSXVDdVMKTlZUYU0rSzZyQlRoN0NLRjZUWlNqQW55SmZoRmRlT1ZKNzlNZDEzYWVJaG0zNTB6UWc3dWZKL2drdkorNUR2TApacC9uQW9JQkFBVlg4WHpFTzdMVk1sbENKZFVUdURTdXNPcDNrQVlFZ2dZNFFRM3FNTlcxRnl5WEM3WjBGOWJFCmtTSEhkalJtU2RUbFZueGN2UW1KTS9WL2tJanpNUHhFT3NCS1BVVkR3N3BhOHdiejlGcTRPOCtJb3lqN1ZXclcKMmExL1FNWXlzSGlpTlBzNy8vWUtvMy9rdkhCWUY2SnNkenkyQkVSTkQ0aTlVOWhDN0RqcGxKR3BSNktMTVBsegpNWFJ3VjVTM2V3cnBXZVcxQW5ONC9EKy9zUGlNQTNnS0swSlBFdGVBV1dndEZHTnNBSkJnaFBoUExxQi9CcDUvCkhOeC96M0w0MWtqRnpqOHNWaHMrVDRZYlhiaGF2R2xxc2h5ZldQbnRhV1VOMG15MjU1RFdQUDhWa24yeFNlV2kKVm1hVW5TSDBTZ2tlUENMRnlra25yQzgxU2pXZkRBMD0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo="
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: compass-installation-config-overrides
  namespace: compass-installer
  labels:
    component: compass
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.isLocalEnv: "true"
  global.domainName: "kyma.local"
  global.adminPassword: ""
  global.minikubeIP: ""
  global.ingress.domainName: "kyma.local"
  global.externalServicesMock.enabled: "true"
  global.externalServicesMock.auditlog: "true"
  global.systemFetcher.enabled: "true"
  global.systemFetcher.systemsAPIEndpoint: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/systemfetcher/systems"
  global.systemFetcher.systemsAPIFilterCriteria: "no"
  global.systemFetcher.systemsAPIFilterTenantCriteriaPattern: "tenant=%s"
  global.systemFetcher.systemToTemplateMappings: '[{"Name": "temp1", "SourceKey": ["prop"], "SourceValue": ["val1"] }, {"Name": "temp2", "SourceKey": ["prop"], "SourceValue": ["val2"] }]'
  global.systemFetcher.oauth.client: "client_id"
  global.systemFetcher.oauth.secret: "client_secret"
  global.systemFetcher.oauth.tokenBaseUrl: "compass-external-services-mock"
  global.systemFetcher.oauth.tokenPath: "/secured/oauth/token"
  global.systemFetcher.oauth.tokenEndpointProtocol: "https"
  global.systemFetcher.oauth.scopesClaim: "scopes"
  global.systemFetcher.oauth.tenantHeaderName: "x-zid"
  global.systemFetcher.oauth.tokenRequestTimeout: "10s"
  global.systemFetcher.oauth.skipSSLValidation: "true"
  global.oathkeeper.mutators.authenticationMappingServices.tenant-fetcher.authenticator.enabled: "true"
  global.oathkeeper.mutators.authenticationMappingServices.subscriber.authenticator.enabled: "true"
  gateway.gateway.auditlog.enabled: "true"
  gateway.gateway.auditlog.authMode: "oauth"
  system-broker.http.client.skipSSLValidation: "true"
  connector.http.client.skipSSLValidation: "true"
  operations-controller.http.client.skipSSLValidation: "true"
  global.systemFetcher.http.client.skipSSLValidation: "true"
  global.ordAggregator.http.client.skipSSLValidation: "true"
  global.http.client.skipSSLValidation: "true"
  global.tests.http.client.skipSSLValidation.director: "true"
  global.tests.http.client.skipSSLValidation.ordService: "true"

  global.tenantFetchers.account-fetcher.enabled: "true"
  global.tenantFetchers.account-fetcher.dbPool.maxOpenConnections: "1"
  global.tenantFetchers.account-fetcher.dbPool.maxIdleConnections: "1"
  global.tenantFetchers.account-fetcher.manageSecrets: "true"
  global.tenantFetchers.account-fetcher.secret.name: "compass-account-fetcher-secret"
  global.tenantFetchers.account-fetcher.secret.clientIdKey: "client-id"
  global.tenantFetchers.account-fetcher.secret.clientSecretKey: "client-secret"
  global.tenantFetchers.account-fetcher.secret.oauthUrlKey: "url"
  global.tenantFetchers.account-fetcher.oauth.client: "client_id"
  global.tenantFetchers.account-fetcher.oauth.secret: "client_secret"
  global.tenantFetchers.account-fetcher.oauth.tokenURL: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/secured/oauth/token"
  global.tenantFetchers.account-fetcher.providerName: "external-svc-mock"
  global.tenantFetchers.account-fetcher.schedule: "0 23 * * *"
  global.tenantFetchers.account-fetcher.kubernetes.configMapNamespace: "compass-system"
  global.tenantFetchers.account-fetcher.kubernetes.pollInterval: "2s"
  global.tenantFetchers.account-fetcher.kubernetes.pollTimeout: "1m"
  global.tenantFetchers.account-fetcher.kubernetes.timeout: "2m"
  global.tenantFetchers.account-fetcher.fieldMapping.idField: "guid"
  global.tenantFetchers.account-fetcher.fieldMapping.nameField: "displayName"
  global.tenantFetchers.account-fetcher.fieldMapping.customerIdField: "customerId"
  global.tenantFetchers.account-fetcher.fieldMapping.discriminatorField: ""
  global.tenantFetchers.account-fetcher.fieldMapping.discriminatorValue: ""
  global.tenantFetchers.account-fetcher.fieldMapping.totalPagesField: "totalPages"
  global.tenantFetchers.account-fetcher.fieldMapping.totalResultsField: "totalResults"
  global.tenantFetchers.account-fetcher.fieldMapping.tenantEventsField: "events"
  global.tenantFetchers.account-fetcher.fieldMapping.detailsField: "eventData"
  global.tenantFetchers.account-fetcher.fieldMapping.entityTypeField: "type"
  global.tenantFetchers.account-fetcher.queryMapping.pageNumField: "page"
  global.tenantFetchers.account-fetcher.queryMapping.pageSizeField: "resultsPerPage"
  global.tenantFetchers.account-fetcher.queryMapping.timestampField: "ts"
  global.tenantFetchers.account-fetcher.query.startPage: "1"
  global.tenantFetchers.account-fetcher.query.pageSize: "1000"
  global.tenantFetchers.account-fetcher.shouldSyncSubaccounts: "false"
  global.tenantFetchers.account-fetcher.accountRegion: "local"
  global.tenantFetchers.account-fetcher.endpoints.accountCreated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-create"
  global.tenantFetchers.account-fetcher.endpoints.accountDeleted: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-delete"
  global.tenantFetchers.account-fetcher.endpoints.accountUpdated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-update"

  global.tenantFetchers.subaccount-fetcher.enabled: "true"
  global.tenantFetchers.subaccount-fetcher.dbPool.maxOpenConnections: "1"
  global.tenantFetchers.subaccount-fetcher.dbPool.maxIdleConnections: "1"
  global.tenantFetchers.subaccount-fetcher.manageSecrets: "true"
  global.tenantFetchers.subaccount-fetcher.secret.name: "compass-subaccount-fetcher-secret"
  global.tenantFetchers.subaccount-fetcher.secret.clientIdKey: "client-id"
  global.tenantFetchers.subaccount-fetcher.secret.clientSecretKey: "client-secret"
  global.tenantFetchers.subaccount-fetcher.secret.oauthUrlKey: "url"
  global.tenantFetchers.subaccount-fetcher.oauth.client: "client_id"
  global.tenantFetchers.subaccount-fetcher.oauth.secret: "client_secret"
  global.tenantFetchers.subaccount-fetcher.oauth.tokenURL: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/secured/oauth/token"
  global.tenantFetchers.subaccount-fetcher.providerName: "subaccount-fetcher"
  global.tenantFetchers.subaccount-fetcher.schedule: "0 23 * * *"
  global.tenantFetchers.subaccount-fetcher.kubernetes.configMapNamespace: "compass-system"
  global.tenantFetchers.subaccount-fetcher.kubernetes.pollInterval: "2s"
  global.tenantFetchers.subaccount-fetcher.kubernetes.pollTimeout: "1m"
  global.tenantFetchers.subaccount-fetcher.kubernetes.timeout: "2m"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.idField: "guid"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.nameField: "displayName"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.customerIdField: "customerId"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.discriminatorField: ""
  global.tenantFetchers.subaccount-fetcher.fieldMapping.discriminatorValue: ""
  global.tenantFetchers.subaccount-fetcher.fieldMapping.totalPagesField: "totalPages"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.totalResultsField: "totalResults"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.tenantEventsField: "events"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.detailsField: "eventData"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.entityTypeField: "type"
  global.tenantFetchers.subaccount-fetcher.queryMapping.pageNumField: "page"
  global.tenantFetchers.subaccount-fetcher.queryMapping.pageSizeField: "resultsPerPage"
  global.tenantFetchers.subaccount-fetcher.queryMapping.timestampField: "ts"
  global.tenantFetchers.subaccount-fetcher.query.startPage: "1"
  global.tenantFetchers.subaccount-fetcher.query.pageSize: "1000"
  global.tenantFetchers.subaccount-fetcher.shouldSyncSubaccounts: "true"
  global.tenantFetchers.subaccount-fetcher.accountRegion: "local"
  global.tenantFetchers.subaccount-fetcher.subaccountRegions: "test"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountCreated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-create"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountDeleted: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-delete"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountUpdated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-update"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountMoved: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-move"

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: compass-installation-oidc-kubeconfig-service-overrides
  namespace: compass-installer
  labels:
    component: oidc-kubeconfig-service
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.isLocalEnv: "true"
  global.domainName: "kyma.local"
  global.minikubeIP: ""
  global.ingress.domainName: "kyma.local"
  config.oidc.caFile: "/etc/dex-tls-cert/tls.crt"
