package senzingrestservice

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/response"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-json-type-definition/go/typedef"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicSenzingRestService is...
type BasicSenzingRestService struct {
	abstractFactory         senzing.SzAbstractFactory
	abstractFactorySyncOnce sync.Once
	senzingrestapi.UnimplementedHandler
	GrpcDialOptions          []grpc.DialOption
	GrpcTarget               string
	isTrace                  bool
	logger                   logging.Logging
	LogLevelName             string
	ObserverOrigin           string
	Observers                []observer.Observer
	OpenAPISpecificationSpec []byte
	Port                     int
	Settings                 string
	SenzingInstanceName      string
	SenzingVerboseLogging    int64
	szConfigManagerSingleton senzing.SzConfigManager
	szConfigManagerSyncOnce  sync.Once
	szProductSingleton       senzing.SzProduct
	szProductSyncOnce        sync.Once
	URLRoutePrefix           string
}

const (
	optionCallerSkip = 3
)

// ----------------------------------------------------------------------------
// Interface methods
// See https://github.com/senzing-garage/go-rest-api-service/blob/main/senzingrestpapi/oas_unimplemented_gen.go
// ----------------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) AddDataSources(
	ctx context.Context,
	req senzingrestapi.AddDataSourcesReq,
	params senzingrestapi.AddDataSourcesParams,
) (senzingrestapi.AddDataSourcesRes, error) {
	var (
		err    error
		result senzingrestapi.AddDataSourcesRes
	)

	_ = ctx
	_ = req

	if restApiService.isTrace {
		entryTime := time.Now()

		restApiService.traceEntry(99)

		defer func() { restApiService.traceExit(99, err, time.Since(entryTime)) }()
	}

	// URL parameters.

	dataSources := params.DataSource
	withRaw := params.WithRaw

	// Get Senzing resources.

	szConfig := restApiService.getSzConfig(ctx)

	// Add DataSouces to in-memory version of Senzing Configuration.

	sdkResponses := []string{}

	for _, dataSource := range params.DataSource {
		sdkResponse, err := szConfig.RegisterDataSource(ctx, dataSource)
		if err != nil {
			return result, wraperror.Errorf(err, "RegisterDataSource")
		}

		sdkResponses = append(sdkResponses, sdkResponse)
	}

	// Persist in-memory Senzing Configuration to Senzing database SYS_CFG table.

	err = restApiService.persistConfiguration(ctx, szConfig)
	if err != nil {
		restApiService.log(9998, dataSources, withRaw, err)
	}

	// Construct response.

	// Retrieve all DataSources

	// rawData, err := szConfig.ListDataSources(ctx, configurationHandle)
	// if err != nil {
	// 	return r, err
	// }

	// fmt.Printf(">>>>>> ListDataSources: %s\n", rawData)

	fmt.Println(sdkResponses) //nolint

	// type SzDataSource struct {
	// 	// The data source code.
	// 	DataSourceCode OptString `json:"dataSourceCode"`
	// 	// The data source ID. The value can be null when used for input in creating a data source to
	// 	// indicate that the data source ID should be auto-generated.
	// 	DataSourceId OptNilInt32 `json:"dataSourceId"`
	// }

	// szDataSource := &senzingrestapi.SzDataSource{
	// 	DataSourceCode: senzingrestapi.NewOptString("DataSourceCodeBob"),
	// 	DataSourceId:   senzingrestapi.NewOptNilInt32(1),
	// }

	// type SzDataSourcesResponseDataDataSourceDetails map[string]SzDataSource

	// szDataSourcesResponseDataDataSourceDetails := &senzingrestapi.SzDataSourcesResponseDataDataSourceDetails{
	// 	"xxxBob": *szDataSource,
	// }

	// type OptSzDataSourcesResponseDataDataSourceDetails struct {
	// 	Value SzDataSourcesResponseDataDataSourceDetails
	// 	Set   bool
	// }

	// optSzDataSourcesResponseDataDataSourceDetails := &senzingrestapi.OptSzDataSourcesResponseDataDataSourceDetails{
	// 	Value: *szDataSourcesResponseDataDataSourceDetails,
	// 	Set:   true,
	// }

	// type SzDataSourcesResponseData struct {
	// 	// The list of data source codes for the configured data sources.
	// 	DataSources []string `json:"dataSources"`
	// 	// The list of `SzDataSource` instances describing the data sources that are configured.
	// 	DataSourceDetails OptSzDataSourcesResponseDataDataSourceDetails `json:"dataSourceDetails"`
	// }

	// szDataSourcesResponseData := &senzingrestapi.SzDataSourcesResponseData{
	// 	DataSources:       []string{"Bobber"},
	// 	DataSourceDetails: *optSzDataSourcesResponseDataDataSourceDetails,
	// }

	// type OptSzDataSourcesResponseData struct {
	// 	Value SzDataSourcesResponseData
	// 	Set   bool
	// }

	// optSzDataSourcesResponseData := &senzingrestapi.OptSzDataSourcesResponseData{
	// 	Value: *szDataSourcesResponseData,
	// 	Set:   true,
	// }

	// type SzDataSourcesResponse struct {
	// 	Data OptSzDataSourcesResponseData `json:"data"`
	// }

	// result = &senzingrestapi.SzDataSourcesResponse{
	// 	Data: *optSzDataSourcesResponseData,
	// }

	// Condensed version of "r"

	result = &senzingrestapi.SzDataSourcesResponse{
		Links: restApiService.getOptSzLinks(ctx, "data-sources"),
		Meta:  restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
		Data: senzingrestapi.OptSzDataSourcesResponseData{
			Set: true,
			Value: senzingrestapi.SzDataSourcesResponseData{
				DataSources: []string{"Bobber"},
				DataSourceDetails: senzingrestapi.OptSzDataSourcesResponseDataDataSourceDetails{
					Set: true,
					Value: senzingrestapi.SzDataSourcesResponseDataDataSourceDetails{
						"xxxBob": senzingrestapi.SzDataSource{
							DataSourceCode: senzingrestapi.NewOptString("BOBBER5"),
							DataSourceId:   senzingrestapi.NewOptNilInt32(1),
						},
					},
				},
			},
		},
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (restApiService *BasicSenzingRestService) Heartbeat(
	ctx context.Context,
) (*senzingrestapi.SzBaseResponse, error) {
	var err error

	result := &senzingrestapi.SzBaseResponse{
		Links: restApiService.getOptSzLinks(ctx, "heartbeat"),
		Meta:  restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (restApiService *BasicSenzingRestService) License(
	ctx context.Context,
	params senzingrestapi.LicenseParams,
) (senzingrestapi.LicenseRes, error) {
	_ = params

	aresponse, err := restApiService.getSzproduct(ctx).GetLicense(ctx)
	if err != nil {
		return nil, wraperror.Errorf(err, "GetLicense")
	}

	parsedResponse, err := response.SzProductGetLicense(ctx, aresponse)
	if err != nil {
		return nil, wraperror.Errorf(err, "SzProductGetLicense")
	}

	issueDate, err := time.Parse("2006-01-02", parsedResponse.IssueDate)
	panicOnError(err)
	expireDate, err := time.Parse("2006-01-02", parsedResponse.ExpireDate)
	panicOnError(err)

	result := &senzingrestapi.SzLicenseResponse{
		Links:   restApiService.getOptSzLinks(ctx, "license"),
		Meta:    restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
		RawData: senzingrestapi.OptNilSzLicenseResponseRawData{},
		Data: senzingrestapi.OptSzLicenseResponseData{
			Set: true,
			Value: senzingrestapi.SzLicenseResponseData{
				License: senzingrestapi.OptSzLicenseInfo{
					Set: true,
					Value: senzingrestapi.SzLicenseInfo{
						Customer:       senzingrestapi.NewOptString(parsedResponse.Customer),
						Contract:       senzingrestapi.NewOptString(parsedResponse.Contract),
						LicenseType:    senzingrestapi.NewOptString(parsedResponse.LicenseType),
						LicenseLevel:   senzingrestapi.NewOptString(parsedResponse.LicenseLevel),
						Billing:        senzingrestapi.NewOptString(parsedResponse.Billing),
						IssuanceDate:   senzingrestapi.NewOptDateTime(issueDate),
						ExpirationDate: senzingrestapi.NewOptDateTime(expireDate),
						RecordLimit:    senzingrestapi.NewOptInt64(parsedResponse.RecordLimit),
					},
				},
			},
		},
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (restApiService *BasicSenzingRestService) OpenAPISpecification(
	ctx context.Context,
) (senzingrestapi.OpenAPISpecificationOKDefault, error) {
	var err error

	_ = ctx
	result := senzingrestapi.OpenAPISpecificationOKDefault{
		// Links: restApiService.getOptSzLinks(ctx, "specifications/open-api"),
		// Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data: bytes.NewReader(restApiService.OpenAPISpecificationSpec),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (restApiService *BasicSenzingRestService) Version(
	ctx context.Context,
	params senzingrestapi.VersionParams,
) (senzingrestapi.VersionRes, error) {
	_ = params
	parsedResponse, err := restApiService.getSenzingVersion(ctx)
	panicOnError(err)
	nativeAPIBuildDate, err := time.Parse("2006-01-02", parsedResponse.BuildDate)
	result := &senzingrestapi.SzVersionResponse{
		Links: restApiService.getOptSzLinks(ctx, "version"),
		Meta:  restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
		Data: senzingrestapi.OptSzVersionInfo{
			Set: true,
			Value: senzingrestapi.SzVersionInfo{
				ApiServerVersion:      senzingrestapi.NewOptString("0.0.0"),
				RestApiVersion:        senzingrestapi.NewOptString("3.4.1"),
				NativeApiVersion:      senzingrestapi.NewOptString(parsedResponse.Version),
				NativeApiBuildVersion: senzingrestapi.NewOptString(parsedResponse.BuildVersion),
				NativeApiBuildNumber:  senzingrestapi.NewOptString(parsedResponse.BuildVersion),
				NativeApiBuildDate:    senzingrestapi.NewOptDateTime(nativeAPIBuildDate),
				ConfigCompatibilityVersion: senzingrestapi.NewOptString(
					parsedResponse.CompatibilityVersion.ConfigVersion,
				),
			},
		},
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	_ = ctx

	restApiService.LogLevelName = logLevelName
	if logLevelName == "TRACE" {
		restApiService.isTrace = true
	}

	return err
}

// ----------------------------------------------------------------------------
// internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (restApiService *BasicSenzingRestService) getLogger() logging.Logging {
	var err error

	if restApiService.logger == nil {
		loggerOptions := []interface{}{
			logging.OptionCallerSkip{Value: optionCallerSkip},
			logging.OptionLogLevel{Value: restApiService.LogLevelName},
		}
		restApiService.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, loggerOptions...)
		panicOnError(err)
	}

	return restApiService.logger
}

// Log message.
func (restApiService *BasicSenzingRestService) log(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (restApiService *BasicSenzingRestService) traceEntry(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (restApiService *BasicSenzingRestService) traceExit(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

func (restApiService *BasicSenzingRestService) getAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error

	_ = ctx

	restApiService.abstractFactorySyncOnce.Do(func() {
		if len(restApiService.GrpcTarget) == 0 {
			restApiService.abstractFactory, err = szfactorycreator.CreateCoreAbstractFactory(
				restApiService.SenzingInstanceName,
				restApiService.Settings,
				restApiService.SenzingVerboseLogging,
				senzing.SzInitializeWithDefaultConfiguration,
			)
			panicOnError(err)
		} else {
			grpcConnection, err := grpc.NewClient(restApiService.GrpcTarget, restApiService.GrpcDialOptions...)
			panicOnError(err)
			restApiService.abstractFactory, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
			panicOnError(err)
		}
	})

	return restApiService.abstractFactory
}

// Singleton pattern for szconfig.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getSzConfig(ctx context.Context) senzing.SzConfig {
	szConfigManager := restApiService.getSzConfigmgr(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	panicOnError(err)
	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	panicOnError(err)

	return szConfig
}

// Singleton pattern for szconfigmanager.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getSzConfigmgr(ctx context.Context) senzing.SzConfigManager {
	var err error

	restApiService.szConfigManagerSyncOnce.Do(func() {
		restApiService.szConfigManagerSingleton, err = restApiService.getAbstractFactory(ctx).CreateConfigManager(ctx)
		panicOnError(err)
	})

	return restApiService.szConfigManagerSingleton
}

// Singleton pattern for szproduct.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getSzproduct(ctx context.Context) senzing.SzProduct {
	var err error

	restApiService.szProductSyncOnce.Do(func() {
		restApiService.szProductSingleton, err = restApiService.getAbstractFactory(ctx).CreateProduct(ctx)
		panicOnError(err)
	})

	return restApiService.szProductSingleton
}

// --- Misc -------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) getOptSzLinks(
	ctx context.Context,
	uriPath string,
) senzingrestapi.OptSzLinks {
	var result senzingrestapi.OptSzLinks

	szLinks := senzingrestapi.SzLinks{
		Self: senzingrestapi.NewOptString(
			fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.URLRoutePrefix, uriPath),
		),
		OpenApiSpecification: senzingrestapi.NewOptString(
			fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.URLRoutePrefix),
		),
	}
	result = senzingrestapi.NewOptSzLinks(szLinks)

	return result
}

func (restApiService *BasicSenzingRestService) getOptSzMeta(
	ctx context.Context,
	httpMethod senzingrestapi.SzHttpMethod,
	httpStatusCode int16,
) senzingrestapi.OptSzMeta {
	var result senzingrestapi.OptSzMeta

	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	panicOnError(err)
	nativeAPIBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	panicOnError(err)

	szMeta := senzingrestapi.SzMeta{
		Server:                     senzingrestapi.NewOptString("Senzing REST API Server - go"),
		HttpMethod:                 senzingrestapi.NewOptSzHttpMethod(httpMethod),
		HttpStatusCode:             senzingrestapi.NewOptInt16(httpStatusCode),
		Timestamp:                  senzingrestapi.NewOptDateTime(time.Now().UTC()),
		Version:                    senzingrestapi.NewOptString("0.0.0"),
		RestApiVersion:             senzingrestapi.NewOptString("3.4.1"),
		NativeApiVersion:           senzingrestapi.NewOptString(senzingVersion.Version),
		NativeApiBuildVersion:      senzingrestapi.NewOptString(senzingVersion.BuildVersion),
		NativeApiBuildNumber:       senzingrestapi.NewOptString(senzingVersion.BuildNumber),
		NativeApiBuildDate:         senzingrestapi.NewOptDateTime(nativeAPIBuildDate),
		ConfigCompatibilityVersion: senzingrestapi.NewOptString(senzingVersion.CompatibilityVersion.ConfigVersion),
		Timings:                    senzingrestapi.NewOptNilSzMetaTimings(map[string]int64{}),
	}
	result = senzingrestapi.NewOptSzMeta(szMeta)

	return result
}

// --- Senzing convenience ----------------------------------------------------

// Persist in-memory Senzing Configuration to Senzing database SYS_CFG table.
func (restApiService *BasicSenzingRestService) persistConfiguration(
	ctx context.Context,
	szConfig senzing.SzConfig,
) error {
	var err error

	szConfigManager := restApiService.getSzConfigmgr(ctx)

	newConfigurationString, err := szConfig.Export(ctx)
	if err != nil {
		return wraperror.Errorf(err, "Export")
	}

	newConfigID, err := szConfigManager.RegisterConfig(ctx, newConfigurationString, "FIXME: description")
	if err != nil {
		return wraperror.Errorf(err, "RegisterConfig")
	}

	err = szConfigManager.SetDefaultConfigID(ctx, newConfigID)
	if err != nil {
		return wraperror.Errorf(err, "SetDefaultConfigID")
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (restApiService *BasicSenzingRestService) getSenzingVersion(
	ctx context.Context,
) (*typedef.SzProductGetVersionResponse, error) {
	var (
		result *typedef.SzProductGetVersionResponse
		err    error
	)

	aresponse, err := restApiService.getSzproduct(ctx).GetVersion(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "getSzproduct")
	}

	result, err = response.SzProductGetVersion(ctx, aresponse)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// --- Debug ------------------------------------------------------------------

// func printContextInternals(ctx interface{}, inner bool) {
// 	contextValues := reflect.ValueOf(ctx).Elem()
// 	contextKeys := reflect.TypeOf(ctx).Elem()

// 	if !inner {
// 		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
// 	}

// 	if contextKeys.Kind() == reflect.Struct {
// 		for i := 0; i < contextValues.NumField(); i++ {
// 			reflectValue := contextValues.Field(i)
// 			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

// 			reflectField := contextKeys.Field(i)

// 			if reflectField.Name == "Context" {
// 				printContextInternals(reflectValue.Interface(), true)
// 			} else {
// 				fmt.Printf("field name: %+v\n", reflectField.Name)
// 				fmt.Printf("value: %+v\n", reflectValue.Interface())
// 			}
// 		}
// 	} else {
// 		fmt.Printf("context is empty (int)\n")
// 	}
// }

// type serverURLKey struct{}

// func (c *Client) requestURL(ctx context.Context) *url.URL {
// 	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
// 	if !ok {
// 		return c.serverURL
// 	}
// 	return u
// }

// func requestURL(ctx context.Context) *url.URL {
// 	// FIXME: See https://github.com/ogen-go/ogen/issues/930
// 	u, _ := ctx.Value(serverURLKey{}).(*url.URL)
// 	log.Printf("CONTEXT %+v", ctx)
// 	// u, _ := ctx.Value()

// 	return u
// }

func getHostname(ctx context.Context) string {
	_ = ctx
	result := "localhost:9999"

	return result
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
