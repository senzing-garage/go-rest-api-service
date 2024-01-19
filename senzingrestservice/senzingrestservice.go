package senzingrestservice

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/senzing-garage/g2-sdk-json-type-definition/go/typedef"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	api "github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/g2-sdk-go/senzing"
	"github.com/senzing/go-sdk-abstract-factory/factory"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// SenzingRestServiceImpl is...
type SenzingRestServiceImpl struct {
	api.UnimplementedHandler
	abstractFactory                factory.SdkAbstractFactory
	abstractFactorySyncOnce        sync.Once
	g2configmgrSingleton           g2api.G2configmgr
	g2configmgrSyncOnce            sync.Once
	g2configSingleton              g2api.G2config
	g2configSyncOnce               sync.Once
	g2engineSingleton              g2api.G2engine
	g2engineSyncOnce               sync.Once
	g2productSingleton             g2api.G2product
	g2productSyncOnce              sync.Once
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	isTrace                        bool
	logger                         logging.LoggingInterface
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecificationSpec       []byte
	Port                           int
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int64
	UrlRoutePrefix                 string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var debugOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

var traceOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

// ----------------------------------------------------------------------------
// internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (restApiService *SenzingRestServiceImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if restApiService.logger == nil {
		loggerOptions := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		restApiService.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, loggerOptions...)
		if err != nil {
			panic(err)
		}
	}
	return restApiService.logger
}

// Log message.
func (restApiService *SenzingRestServiceImpl) log(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Debug.
func (restApiService *SenzingRestServiceImpl) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (restApiService *SenzingRestServiceImpl) traceEntry(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (restApiService *SenzingRestServiceImpl) traceExit(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (restApiService *SenzingRestServiceImpl) error(messageNumber int, details ...interface{}) error {
	return restApiService.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

func (restApiService *SenzingRestServiceImpl) getAbstractFactory() factory.SdkAbstractFactory {
	restApiService.abstractFactorySyncOnce.Do(func() {
		if len(restApiService.GrpcTarget) == 0 {
			restApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{}
		} else {
			restApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{
				GrpcDialOptions: restApiService.GrpcDialOptions,
				GrpcTarget:      restApiService.GrpcTarget,
				ObserverOrigin:  restApiService.ObserverOrigin,
				Observers:       restApiService.Observers,
			}
		}
	})
	return restApiService.abstractFactory
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *SenzingRestServiceImpl) getG2config(ctx context.Context) g2api.G2config {
	var err error = nil
	restApiService.g2configSyncOnce.Do(func() {
		restApiService.g2configSingleton, err = restApiService.getAbstractFactory().GetG2config(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2configSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2configSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2configSingleton
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *SenzingRestServiceImpl) getG2configmgr(ctx context.Context) g2api.G2configmgr {
	var err error = nil
	restApiService.g2configmgrSyncOnce.Do(func() {
		restApiService.g2configmgrSingleton, err = restApiService.getAbstractFactory().GetG2configmgr(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2configmgrSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2configmgrSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2configmgrSingleton
}

// Singleton pattern for g2engine.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *SenzingRestServiceImpl) getG2engine(ctx context.Context) g2api.G2engine {
	var err error = nil
	restApiService.g2engineSyncOnce.Do(func() {
		restApiService.g2engineSingleton, err = restApiService.getAbstractFactory().GetG2engine(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2engineSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2engineSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2engineSingleton
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *SenzingRestServiceImpl) getG2product(ctx context.Context) g2api.G2product {
	var err error = nil
	restApiService.g2productSyncOnce.Do(func() {
		restApiService.g2productSingleton, err = restApiService.getAbstractFactory().GetG2product(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2productSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2productSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2productSingleton
}

// --- Senzing convenience ----------------------------------------------------

// Pull the Senzing Configuration from the database into an in-memory copy.
func (restApiService *SenzingRestServiceImpl) getConfigurationHandle(ctx context.Context) (uintptr, error) {
	var err error = nil
	var result uintptr
	var configurationString string
	g2Config := restApiService.getG2config(ctx)
	g2Configmgr := restApiService.getG2configmgr(ctx)
	configID, err := g2Configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		return result, err
	}
	if configID == 0 {
		return g2Config.Create(ctx)
	}
	configurationString, err = g2Configmgr.GetConfig(ctx, configID)
	if err != nil {
		return result, err
	}
	result, err = g2Config.Load(ctx, configurationString)
	if err != nil {
		return result, err
	}
	return result, err
}

// Persist in-memory Senzing Configuration to Senzing database SYS_CFG table.
func (restApiService *SenzingRestServiceImpl) persistConfiguration(ctx context.Context, configurationHandle uintptr) error {
	var err error = nil
	g2Config := restApiService.getG2config(ctx)
	g2Configmgr := restApiService.getG2configmgr(ctx)
	newConfigurationString, err := g2Config.Save(ctx, configurationHandle)
	if err != nil {
		return err
	}
	newConfigId, err := g2Configmgr.AddConfig(ctx, newConfigurationString, "FIXME: description")
	if err != nil {
		return err
	}
	err = g2Configmgr.SetDefaultConfigID(ctx, newConfigId)
	if err != nil {
		return err
	}
	return err
}

func (restApiService *SenzingRestServiceImpl) getSenzingVersion(ctx context.Context) (*typedef.ProductVersionResponse, error) {
	response, err := restApiService.getG2product(ctx).Version(ctx)
	if err != nil {
		return nil, err
	}
	return senzing.UnmarshalProductVersionResponse(ctx, response)
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
	result := "localhost:9999"
	return result
}

// ----------------------------------------------------------------------------
// Interface methods
// See https://github.com/senzing-garage/go-rest-api-service/blob/main/senzingrestpapi/oas_unimplemented_gen.go
// ----------------------------------------------------------------------------

func (restApiService *SenzingRestServiceImpl) AddDataSources(ctx context.Context, req api.AddDataSourcesReq, params api.AddDataSourcesParams) (r api.AddDataSourcesRes, _ error) {
	var err error = nil
	if restApiService.isTrace {
		entryTime := time.Now()
		restApiService.traceEntry(99)
		defer func() { restApiService.traceExit(99, err, time.Since(entryTime)) }()
	}

	// URL parameters.

	dataSources := params.DataSource
	withRaw := params.WithRaw

	// Get Senzing resources.

	g2Config := restApiService.getG2config(ctx)

	// Get current configuration from database into memory.

	configurationHandle, err := restApiService.getConfigurationHandle(ctx)
	if err != nil {
		restApiService.log(9999, dataSources, withRaw, err)
	}

	// Add DataSouces to in-memory version of Senzing Configuration.

	sdkResponses := []string{}
	for _, dataSource := range params.DataSource {
		sdkRequest := fmt.Sprintf(`{"DSRC_CODE": "%s"}`, dataSource)
		sdkResponse, err := g2Config.AddDataSource(ctx, configurationHandle, sdkRequest)
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

	err = g2Config.Close(ctx, configurationHandle)

	fmt.Println(sdkResponses)

	// type SzDataSource struct {
	// 	// The data source code.
	// 	DataSourceCode OptString `json:"dataSourceCode"`
	// 	// The data source ID. The value can be null when used for input in creating a data source to
	// 	// indicate that the data source ID should be auto-generated.
	// 	DataSourceId OptNilInt32 `json:"dataSourceId"`
	// }

	szDataSource := &api.SzDataSource{
		DataSourceCode: api.NewOptString("DataSourceCodeBob"),
		DataSourceId:   api.NewOptNilInt32(1),
	}

	// type SzDataSourcesResponseDataDataSourceDetails map[string]SzDataSource

	szDataSourcesResponseDataDataSourceDetails := &api.SzDataSourcesResponseDataDataSourceDetails{
		"xxxBob": *szDataSource,
	}

	// type OptSzDataSourcesResponseDataDataSourceDetails struct {
	// 	Value SzDataSourcesResponseDataDataSourceDetails
	// 	Set   bool
	// }

	optSzDataSourcesResponseDataDataSourceDetails := &api.OptSzDataSourcesResponseDataDataSourceDetails{
		Value: *szDataSourcesResponseDataDataSourceDetails,
		Set:   true,
	}

	// type SzDataSourcesResponseData struct {
	// 	// The list of data source codes for the configured data sources.
	// 	DataSources []string `json:"dataSources"`
	// 	// The list of `SzDataSource` instances describing the data sources that are configured.
	// 	DataSourceDetails OptSzDataSourcesResponseDataDataSourceDetails `json:"dataSourceDetails"`
	// }

	szDataSourcesResponseData := &api.SzDataSourcesResponseData{
		DataSources:       []string{"Bobber"},
		DataSourceDetails: *optSzDataSourcesResponseDataDataSourceDetails,
	}

	// type OptSzDataSourcesResponseData struct {
	// 	Value SzDataSourcesResponseData
	// 	Set   bool
	// }

	optSzDataSourcesResponseData := &api.OptSzDataSourcesResponseData{
		Value: *szDataSourcesResponseData,
		Set:   true,
	}

	// type SzDataSourcesResponse struct {
	// 	Data OptSzDataSourcesResponseData `json:"data"`
	// }

	r = &api.SzDataSourcesResponse{
		Data: *optSzDataSourcesResponseData,
	}

	// Condensed version of "r"

	r = &api.SzDataSourcesResponse{
		Links: restApiService.getOptSzLinks(ctx, "data-sources"),
		Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data: api.OptSzDataSourcesResponseData{
			Set: true,
			Value: api.SzDataSourcesResponseData{
				DataSources: []string{"Bobber"},
				DataSourceDetails: api.OptSzDataSourcesResponseDataDataSourceDetails{
					Set: true,
					Value: api.SzDataSourcesResponseDataDataSourceDetails{
						"xxxBob": api.SzDataSource{
							DataSourceCode: api.NewOptString("BOBBER5"),
							DataSourceId:   api.NewOptNilInt32(1),
						},
					},
				},
			},
		},
	}

	return r, err
}

func (restApiService *SenzingRestServiceImpl) GetServerInfo(ctx context.Context) (r api.GetServerInfoRes, _ error) {
	var err error = nil
	r = &api.SzServerInfoResponse{
		Links: restApiService.getOptSzLinks(ctx, "server-info"),
		Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data:  restApiService.getServerInfoData(ctx),
	}
	return r, err
}

func (restApiService *SenzingRestServiceImpl) Heartbeat(ctx context.Context) (r *api.SzBaseResponse, _ error) {
	var err error = nil
	r = &api.SzBaseResponse{
		Links: restApiService.getOptSzLinks(ctx, "heartbeat"),
		Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
	}
	return r, err
}

func (restApiService *SenzingRestServiceImpl) License(ctx context.Context, params api.LicenseParams) (r api.LicenseRes, _ error) {
	license, err := restApiService.getG2product(ctx).License(ctx)
	r = &api.SzLicenseResponse{
		Links:   restApiService.getOptSzLinks(ctx, "license"),
		Meta:    restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		RawData: api.OptNilSzLicenseResponseRawData{Set: true, Value: api.SzLicenseResponseRawData{}},
		Data:    licenseData(ctx, license),
	}
	return r, err
}

func (restApiService *SenzingRestServiceImpl) OpenApiSpecification(ctx context.Context) (r api.OpenApiSpecificationOKDefault, _ error) {
	var err error = nil
	r = api.OpenApiSpecificationOKDefault{
		// Links: restApiService.getOptSzLinks(ctx, "specifications/open-api"),
		// Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data: restApiService.openApiSpecificationData(ctx),
	}
	return r, err
}

func (restApiService *SenzingRestServiceImpl) SearchEntitiesByGet(ctx context.Context, params api.SearchEntitiesByGetParams) (r api.SearchEntitiesByGetRes, _ error) {
	response, err := restApiService.getG2engine(ctx).SearchByAttributes(ctx, params.Attrs.Value)
	r = &api.SearchEntitiesByGetOKDefault{
		// Links: restApiService.getOptSzLinks(ctx, "license"),
		// Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data: restApiService.searchEntitiesByGetData(ctx, response),
	}
	return r, err
}

func (restApiService *SenzingRestServiceImpl) Version(ctx context.Context, params api.VersionParams) (r api.VersionRes, _ error) {
	response, err := restApiService.getG2product(ctx).Version(ctx)
	r = &api.SzVersionResponse{
		Links: restApiService.getOptSzLinks(ctx, "version"),
		Meta:  restApiService.getOptSzMeta(ctx, api.SzHttpMethodGET, http.StatusOK),
		Data:  versionData(ctx, response),
	}
	return r, err
}
