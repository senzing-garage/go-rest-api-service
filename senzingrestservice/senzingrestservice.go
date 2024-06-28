package senzingrestservice

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

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
	szConfigSingleton        senzing.SzConfig
	szConfigSyncOnce         sync.Once
	szProductSingleton       senzing.SzProduct
	szProductSyncOnce        sync.Once
	URLRoutePrefix           string
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
			&logging.OptionCallerSkip{Value: 3},
		}
		restApiService.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, loggerOptions...)
		if err != nil {
			panic(err)
		}
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
			restApiService.abstractFactory, err = szfactorycreator.CreateCoreAbstractFactory(restApiService.SenzingInstanceName, restApiService.Settings, restApiService.SenzingVerboseLogging, senzing.SzInitializeWithDefaultConfiguration)
			if err != nil {
				panic(err)
			}
		} else {
			grpcConnection, err := grpc.NewClient(restApiService.GrpcTarget, restApiService.GrpcDialOptions...)
			if err != nil {
				panic(err)
			}
			restApiService.abstractFactory, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.abstractFactory
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getG2config(ctx context.Context) senzing.SzConfig {
	var err error
	restApiService.szConfigSyncOnce.Do(func() {
		restApiService.szConfigSingleton, err = restApiService.getAbstractFactory(ctx).CreateSzConfig(ctx)
		if err != nil {
			panic(err)
		}
	})
	return restApiService.szConfigSingleton
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getG2configmgr(ctx context.Context) senzing.SzConfigManager {
	var err error
	restApiService.szConfigManagerSyncOnce.Do(func() {
		restApiService.szConfigManagerSingleton, err = restApiService.getAbstractFactory(ctx).CreateSzConfigManager(ctx)
		if err != nil {
			panic(err)
		}
	})
	return restApiService.szConfigManagerSingleton
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *BasicSenzingRestService) getG2product(ctx context.Context) senzing.SzProduct {
	var err error
	restApiService.szProductSyncOnce.Do(func() {
		restApiService.szProductSingleton, err = restApiService.getAbstractFactory(ctx).CreateSzProduct(ctx)
		if err != nil {
			panic(err)
		}
	})
	return restApiService.szProductSingleton
}

// --- Misc -------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) getOptSzLinks(ctx context.Context, uriPath string) senzingrestapi.OptSzLinks {
	var result senzingrestapi.OptSzLinks
	szLinks := senzingrestapi.SzLinks{
		Self:                 senzingrestapi.NewOptString(fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.URLRoutePrefix, uriPath)),
		OpenApiSpecification: senzingrestapi.NewOptString(fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.URLRoutePrefix)),
	}
	result = senzingrestapi.NewOptSzLinks(szLinks)
	return result
}

func (restApiService *BasicSenzingRestService) getOptSzMeta(ctx context.Context, httpMethod senzingrestapi.SzHttpMethod, httpStatusCode int16) senzingrestapi.OptSzMeta {
	var result senzingrestapi.OptSzMeta

	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}

	nativeAPIBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	if err != nil {
		panic(err)
	}

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

// Pull the Senzing Configuration from the database into an in-memory copy.
func (restApiService *BasicSenzingRestService) getConfigurationHandle(ctx context.Context) (uintptr, error) {
	var err error
	var result uintptr
	var configurationString string
	szConfig := restApiService.getG2config(ctx)
	szConfigManager := restApiService.getG2configmgr(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		return result, err
	}
	if configID == 0 {
		return szConfig.CreateConfig(ctx)
	}
	configurationString, err = szConfigManager.GetConfig(ctx, configID)
	if err != nil {
		return result, err
	}
	result, err = szConfig.ImportConfig(ctx, configurationString)
	if err != nil {
		return result, err
	}
	return result, err
}

// Persist in-memory Senzing Configuration to Senzing database SYS_CFG table.
func (restApiService *BasicSenzingRestService) persistConfiguration(ctx context.Context, configurationHandle uintptr) error {
	var err error
	szConfig := restApiService.getG2config(ctx)
	szConfigManager := restApiService.getG2configmgr(ctx)
	newConfigurationString, err := szConfig.ExportConfig(ctx, configurationHandle)
	if err != nil {
		return err
	}
	newConfigID, err := szConfigManager.AddConfig(ctx, newConfigurationString, "FIXME: description")
	if err != nil {
		return err
	}
	err = szConfigManager.SetDefaultConfigID(ctx, newConfigID)
	if err != nil {
		return err
	}
	return err
}

func (restApiService *BasicSenzingRestService) getSenzingVersion(ctx context.Context) (*typedef.SzProductGetVersionResponse, error) {
	aresponse, err := restApiService.getG2product(ctx).GetVersion(ctx)
	if err != nil {
		return nil, err
	}
	return response.SzProductGetVersion(ctx, aresponse)
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

// ----------------------------------------------------------------------------
// Interface methods
// See https://github.com/senzing-garage/go-rest-api-service/blob/main/senzingrestpapi/oas_unimplemented_gen.go
// ----------------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) AddDataSources(ctx context.Context, req senzingrestapi.AddDataSourcesReq, params senzingrestapi.AddDataSourcesParams) (r senzingrestapi.AddDataSourcesRes, _ error) {
	var err error
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

	szConfig := restApiService.getG2config(ctx)

	// Get current configuration from database into memory.

	configurationHandle, err := restApiService.getConfigurationHandle(ctx)
	if err != nil {
		restApiService.log(9999, dataSources, withRaw, err)
	}

	// Add DataSouces to in-memory version of Senzing Configuration.

	sdkResponses := []string{}
	for _, dataSource := range params.DataSource {
		sdkResponse, err := szConfig.AddDataSource(ctx, configurationHandle, dataSource)
		if err != nil {
			return r, err
		}
		sdkResponses = append(sdkResponses, sdkResponse)
	}

	// Persist in-memory Senzing Configuration to Senzing database SYS_CFG table.

	err = restApiService.persistConfiguration(ctx, configurationHandle)
	if err != nil {
		restApiService.log(9999, dataSources, withRaw, err)
	}

	// Construct response.

	// Retrieve all DataSources

	// rawData, err := g2Config.ListDataSources(ctx, configurationHandle)
	// if err != nil {
	// 	return r, err
	// }

	// fmt.Printf(">>>>>> ListDataSources: %s\n", rawData)

	err = szConfig.CloseConfig(ctx, configurationHandle)

	fmt.Println(sdkResponses)

	// type SzDataSource struct {
	// 	// The data source code.
	// 	DataSourceCode OptString `json:"dataSourceCode"`
	// 	// The data source ID. The value can be null when used for input in creating a data source to
	// 	// indicate that the data source ID should be auto-generated.
	// 	DataSourceId OptNilInt32 `json:"dataSourceId"`
	// }

	szDataSource := &senzingrestapi.SzDataSource{
		DataSourceCode: senzingrestapi.NewOptString("DataSourceCodeBob"),
		DataSourceId:   senzingrestapi.NewOptNilInt32(1),
	}

	// type SzDataSourcesResponseDataDataSourceDetails map[string]SzDataSource

	szDataSourcesResponseDataDataSourceDetails := &senzingrestapi.SzDataSourcesResponseDataDataSourceDetails{
		"xxxBob": *szDataSource,
	}

	// type OptSzDataSourcesResponseDataDataSourceDetails struct {
	// 	Value SzDataSourcesResponseDataDataSourceDetails
	// 	Set   bool
	// }

	optSzDataSourcesResponseDataDataSourceDetails := &senzingrestapi.OptSzDataSourcesResponseDataDataSourceDetails{
		Value: *szDataSourcesResponseDataDataSourceDetails,
		Set:   true,
	}

	// type SzDataSourcesResponseData struct {
	// 	// The list of data source codes for the configured data sources.
	// 	DataSources []string `json:"dataSources"`
	// 	// The list of `SzDataSource` instances describing the data sources that are configured.
	// 	DataSourceDetails OptSzDataSourcesResponseDataDataSourceDetails `json:"dataSourceDetails"`
	// }

	szDataSourcesResponseData := &senzingrestapi.SzDataSourcesResponseData{
		DataSources:       []string{"Bobber"},
		DataSourceDetails: *optSzDataSourcesResponseDataDataSourceDetails,
	}

	// type OptSzDataSourcesResponseData struct {
	// 	Value SzDataSourcesResponseData
	// 	Set   bool
	// }

	optSzDataSourcesResponseData := &senzingrestapi.OptSzDataSourcesResponseData{
		Value: *szDataSourcesResponseData,
		Set:   true,
	}

	// type SzDataSourcesResponse struct {
	// 	Data OptSzDataSourcesResponseData `json:"data"`
	// }

	r = &senzingrestapi.SzDataSourcesResponse{
		Data: *optSzDataSourcesResponseData,
	}

	// Condensed version of "r"

	r = &senzingrestapi.SzDataSourcesResponse{
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

	return r, err
}

func (restApiService *BasicSenzingRestService) Heartbeat(ctx context.Context) (r *senzingrestapi.SzBaseResponse, _ error) {
	var err error
	r = &senzingrestapi.SzBaseResponse{
		Links: restApiService.getOptSzLinks(ctx, "heartbeat"),
		Meta:  restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
	}
	return r, err
}

func (restApiService *BasicSenzingRestService) License(ctx context.Context, params senzingrestapi.LicenseParams) (r senzingrestapi.LicenseRes, _ error) {
	_ = params
	aresponse, err := restApiService.getG2product(ctx).GetLicense(ctx)
	if err != nil {
		return nil, err
	}
	parsedResponse, err := response.SzProductGetLicense(ctx, aresponse)
	if err != nil {
		return nil, err
	}
	issueDate, err := time.Parse("2006-01-02", parsedResponse.IssueDate)
	if err != nil {
		panic(err)
	}
	expireDate, err := time.Parse("2006-01-02", parsedResponse.ExpireDate)
	if err != nil {
		panic(err)
	}
	r = &senzingrestapi.SzLicenseResponse{
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
	return r, err
}

func (restApiService *BasicSenzingRestService) OpenAPISpecification(ctx context.Context) (r senzingrestapi.OpenAPISpecificationOKDefault, _ error) {
	var err error
	_ = ctx
	r = senzingrestapi.OpenAPISpecificationOKDefault{
		// Links: restApiService.getOptSzLinks(ctx, "specifications/open-api"),
		// Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data: bytes.NewReader(restApiService.OpenAPISpecificationSpec),
	}
	return r, err
}

func (restApiService *BasicSenzingRestService) Version(ctx context.Context, params senzingrestapi.VersionParams) (r senzingrestapi.VersionRes, _ error) {
	_ = params
	parsedResponse, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}
	nativeAPIBuildDate, err := time.Parse("2006-01-02", parsedResponse.BuildDate)
	if err != nil {
		panic(err)
	}
	r = &senzingrestapi.SzVersionResponse{
		Links: restApiService.getOptSzLinks(ctx, "version"),
		Meta:  restApiService.getOptSzMeta(ctx, senzingrestapi.SzHttpMethodGET, http.StatusOK),
		Data: senzingrestapi.OptSzVersionInfo{
			Set: true,
			Value: senzingrestapi.SzVersionInfo{
				ApiServerVersion:           senzingrestapi.NewOptString("0.0.0"),
				RestApiVersion:             senzingrestapi.NewOptString("3.4.1"),
				NativeApiVersion:           senzingrestapi.NewOptString(parsedResponse.Version),
				NativeApiBuildVersion:      senzingrestapi.NewOptString(parsedResponse.BuildVersion),
				NativeApiBuildNumber:       senzingrestapi.NewOptString(parsedResponse.BuildVersion),
				NativeApiBuildDate:         senzingrestapi.NewOptDateTime(nativeAPIBuildDate),
				ConfigCompatibilityVersion: senzingrestapi.NewOptString(parsedResponse.CompatibilityVersion.ConfigVersion),
			},
		},
	}
	return r, err
}
