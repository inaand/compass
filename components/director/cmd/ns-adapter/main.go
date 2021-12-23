package main

import (
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/director/internal/domain/api"
	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/domain/apptemplate"
	"github.com/kyma-incubator/compass/components/director/internal/domain/auth"
	bundleutil "github.com/kyma-incubator/compass/components/director/internal/domain/bundle"
	"github.com/kyma-incubator/compass/components/director/internal/domain/bundlereferences"
	"github.com/kyma-incubator/compass/components/director/internal/domain/document"
	"github.com/kyma-incubator/compass/components/director/internal/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/fetchrequest"
	"github.com/kyma-incubator/compass/components/director/internal/domain/integrationsystem"
	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/domain/labeldef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal/domain/spec"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/kyma-incubator/compass/components/director/internal/domain/webhook"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/nsadapter/adapter"
	"github.com/kyma-incubator/compass/components/director/internal/nsadapter/handler"
	"github.com/kyma-incubator/compass/components/director/internal/nsadapter/httputil"
	"github.com/kyma-incubator/compass/components/director/internal/nsadapter/nsmodel"
	"github.com/kyma-incubator/compass/components/director/internal/systemfetcher"
	"github.com/kyma-incubator/compass/components/director/internal/uid"
	"github.com/kyma-incubator/compass/components/director/pkg/accessstrategy"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/certloader"
	"github.com/kyma-incubator/compass/components/director/pkg/correlation"
	directorHandler "github.com/kyma-incubator/compass/components/director/pkg/handler"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/normalizer"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/signal"
	"github.com/kyma-incubator/compass/components/director/pkg/str"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
	"net/http"
	"os"
	"strings"
)

const appTemplateName = "S4HANA"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	conf := adapter.Configuration{}
	err := envconfig.InitWithPrefix(&conf, "APP")
	exitOnError(err, "while reading ns adapter configuration")

	transact, closeFunc, err := persistence.Configure(ctx, conf.Database)
	exitOnError(err, "Error while establishing the connection to the database")

	certCache, err := certloader.StartCertLoader(ctx, conf.ExternalClientCertSecret)
	exitOnError(err, "Failed to initialize certificate loader")

	uidSvc := uid.NewService()

	tenantConv := tenant.NewConverter()
	tenantRepo := tenant.NewRepository(tenantConv)

	authConverter := auth.NewConverter()
	frConverter := fetchrequest.NewConverter(authConverter)
	versionConverter := version.NewConverter()
	docConverter := document.NewConverter(frConverter)
	webhookConverter := webhook.NewConverter(authConverter)
	specConverter := spec.NewConverter(frConverter)
	apiConverter := api.NewConverter(versionConverter, specConverter)
	eventAPIConverter := eventdef.NewConverter(versionConverter, specConverter)
	labelDefConverter := labeldef.NewConverter()
	labelConverter := label.NewConverter()
	intSysConverter := integrationsystem.NewConverter()
	bundleConverter := bundleutil.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, bundleConverter)
	runtimeConverter := runtime.NewConverter()
	bundleReferenceConv := bundlereferences.NewConverter()

	runtimeRepo := runtime.NewRepository(runtimeConverter)
	applicationRepo := application.NewRepository(appConverter)
	labelRepo := label.NewRepository(labelConverter)
	labelDefRepo := labeldef.NewRepository(labelDefConverter)
	webhookRepo := webhook.NewRepository(webhookConverter)
	apiRepo := api.NewRepository(apiConverter)
	eventAPIRepo := eventdef.NewRepository(eventAPIConverter)
	specRepo := spec.NewRepository(specConverter)
	docRepo := document.NewRepository(docConverter)
	fetchRequestRepo := fetchrequest.NewRepository(frConverter)
	intSysRepo := integrationsystem.NewRepository(intSysConverter)
	bundleRepo := bundleutil.NewRepository(bundleConverter)
	bundleReferenceRepo := bundlereferences.NewRepository(bundleReferenceConv)

	labelSvc := label.NewLabelService(labelRepo, labelDefRepo, uidSvc)
	assignmentConv := scenarioassignment.NewConverter()
	scenarioAssignmentRepo := scenarioassignment.NewRepository(assignmentConv)
	scenariosSvc := labeldef.NewService(labelDefRepo, labelRepo, scenarioAssignmentRepo, tenantRepo, uidSvc, conf.DefaultScenarioEnabled)
	fetchRequestSvc := fetchrequest.NewService(fetchRequestRepo, &http.Client{Timeout: conf.ClientTimeout}, accessstrategy.NewDefaultExecutorProvider(certCache))
	specSvc := spec.NewService(specRepo, fetchRequestRepo, uidSvc, fetchRequestSvc)
	bundleReferenceSvc := bundlereferences.NewService(bundleReferenceRepo, uidSvc)
	apiSvc := api.NewService(apiRepo, uidSvc, specSvc, bundleReferenceSvc)
	eventAPISvc := eventdef.NewService(eventAPIRepo, uidSvc, specSvc, bundleReferenceSvc)
	docSvc := document.NewService(docRepo, fetchRequestRepo, uidSvc)
	bundleSvc := bundleutil.NewService(bundleRepo, apiSvc, eventAPISvc, docSvc, uidSvc)
	appSvc := application.NewService(&normalizer.DefaultNormalizator{}, nil, applicationRepo, webhookRepo, runtimeRepo, labelRepo, intSysRepo, labelSvc, scenariosSvc, bundleSvc, uidSvc)

	appTemplateConverter := apptemplate.NewConverter(appConverter, webhookConverter)
	appTemplateRepo := apptemplate.NewRepository(appTemplateConverter)
	appTemplateSvc := apptemplate.NewService(appTemplateRepo, webhookRepo, uidSvc)

	tntSvc := tenant.NewService(tenantRepo, uidSvc)

	defer func() {
		err := closeFunc()
		exitOnError(err, "Error while closing the connection to the database")
	}()

	registerAppTemplate(ctx, transact, appTemplateSvc)

	err = calculateTemplateMappings(ctx, conf, transact)
	exitOnError(err, "while calculating template mappings")

	h := handler.NewHandler(appSvc, appConverter, appTemplateSvc, tntSvc, transact)

	handlerWithTimeout, err := directorHandler.WithTimeoutWithErrorMessage(h, conf.ServerTimeout, httputil.GetTimeoutMessage())
	exitOnError(err, "Failed configuring timeout on handler")

	router := mux.NewRouter()

	router.Use(correlation.AttachCorrelationIDToContext())
	router.NewRoute().
		Methods(http.MethodPost).
		Path("/api/v1/notifications").
		Handler(handlerWithTimeout)
	router.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	setValidationMessages()

	server := &http.Server{
		Addr:              conf.Address,
		Handler:           router,
		ReadHeaderTimeout: conf.ServerTimeout,
	}
	ctx, err = log.Configure(ctx, conf.Log)
	exitOnError(err, "while configuring logger")

	log.C(ctx).Infof("API listening on %s", conf.Address)
	exitOnError(server.ListenAndServe(), "on starting HTTP server")
}

func registerAppTemplate(ctx context.Context, transact persistence.Transactioner, appTemplateSvc apptemplate.ApplicationTemplateService) {
	tx, err := transact.Begin()
	if err != nil {
		exitOnError(err, "Error while beginning transaction")
	}
	defer transact.RollbackUnlessCommitted(ctx, tx)
	ctxWithTx := persistence.SaveToContext(ctx, tx)

	appTemplate := model.ApplicationTemplateInput{
		Name:        appTemplateName,
		Description: str.Ptr("Template for systems pushed from Notifications Service"),
		ApplicationInputJSON: `{
									"name": "on-premise-system",
									"description": "{{description}}",
									"providerName": "SAP",
									"labels": {"scc": {"Subaccount":"{{subaccount}}", "LocationId":"{{location-id}}", "Host":"{{host}}"}, "applicationType":"{{system-type}}", "systemProtocol": "{{protocol}}" },
									"systemNumber": "{{system-number}}",
									"systemStatus": "{{system-status}}"
								}`,
		Placeholders: []model.ApplicationTemplatePlaceholder{
			{
				Name:        "description",
				Description: str.Ptr("description of the system"),
			},
			{
				Name:        "subaccount",
				Description: str.Ptr("subaccount to which the scc is connected"),
			},
			{
				Name:        "location-id",
				Description: str.Ptr("location id of the scc"),
			},
			{
				Name:        "system-type",
				Description: str.Ptr("type of the system"),
			},
			{
				Name:        "host",
				Description: str.Ptr("host of the system"),
			},
			{
				Name:        "protocol",
				Description: str.Ptr("protocol of the system"),
			},
			{
				Name:        "system-number",
				Description: str.Ptr("unique identification of the system"),
			},
			{
				Name:        "system-status",
				Description: str.Ptr("describes whether the system is reachable or not"),
			},
		},
		AccessLevel: model.GlobalApplicationTemplateAccessLevel, //TODO check proper access level
	}

	_, err = appTemplateSvc.GetByName(ctxWithTx, appTemplateName)
	if err != nil {
		if !strings.Contains(err.Error(), "Object not found") {
			exitOnError(err, fmt.Sprintf("error while getting application template with name: %s", appTemplateName))
		}

		templateId, err := appTemplateSvc.Create(ctxWithTx, appTemplate)
		if err != nil {
			exitOnError(err, fmt.Sprintf("error while registering application template with name: %s", appTemplateName))
		}
		fmt.Println(fmt.Sprintf("successfully registered application template with id: %s", templateId))
	}

	if err := tx.Commit(); err != nil {
		exitOnError(err, "while committing transaction")
	}
}

func exitOnError(err error, context string) {
	if err != nil {
		wrappedError := errors.Wrap(err, context)
		log.D().Fatal(wrappedError)
	}
}

func calculateTemplateMappings(ctx context.Context, cfg adapter.Configuration, transact persistence.Transactioner) error {
	var systemToTemplateMappings []systemfetcher.TemplateMapping
	if err := json.Unmarshal([]byte(cfg.SystemToTemplateMappings), &systemToTemplateMappings); err != nil {
		return errors.Wrap(err, "failed to read system template mappings")
	}

	authConverter := auth.NewConverter()
	versionConverter := version.NewConverter()
	frConverter := fetchrequest.NewConverter(authConverter)
	webhookConverter := webhook.NewConverter(authConverter)
	specConverter := spec.NewConverter(frConverter)
	apiConverter := api.NewConverter(versionConverter, specConverter)
	eventAPIConverter := eventdef.NewConverter(versionConverter, specConverter)
	docConverter := document.NewConverter(frConverter)
	bundleConverter := bundleutil.NewConverter(authConverter, apiConverter, eventAPIConverter, docConverter)
	appConverter := application.NewConverter(webhookConverter, bundleConverter)
	appTemplateConv := apptemplate.NewConverter(appConverter, webhookConverter)
	appTemplateRepo := apptemplate.NewRepository(appTemplateConv)
	webhookRepo := webhook.NewRepository(webhookConverter)

	uidSvc := uid.NewService()
	appTemplateSvc := apptemplate.NewService(appTemplateRepo, webhookRepo, uidSvc)

	tx, err := transact.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	for index, tm := range systemToTemplateMappings {
		appTemplate, err := appTemplateSvc.GetByName(ctx, tm.Name)
		if err != nil && !apperrors.IsNotFoundError(err) {
			return err
		}
		systemToTemplateMappings[index].ID = appTemplate.ID
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	nsmodel.Mappings = systemToTemplateMappings
	return nil
}

func setValidationMessages() {
	validation.ErrRequired = validation.ErrRequired.SetMessage("the value is required")
	validation.ErrNotNilRequired = validation.ErrNotNilRequired.SetMessage("the value can not be nil")
}
