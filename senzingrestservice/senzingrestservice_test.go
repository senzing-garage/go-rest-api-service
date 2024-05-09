package senzingrestservice

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	api "github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/stretchr/testify/assert"
)

var (
	senzingRestServiceSingleton SenzingRestService
	debug                       bool = false
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) SenzingRestService {
	_ = ctx
	if senzingRestServiceSingleton == nil {
		senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			test.Errorf("Error: %s", err)
		}
		senzingRestServiceSingleton = &SenzingRestServiceImpl{
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              "go-rest-api-service-test",
			SenzingVerboseLogging:          int64(0),
		}
	}
	return senzingRestServiceSingleton
}

func testError(test *testing.T, ctx context.Context, err error) {
	_ = ctx
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingRestServiceImpl_AddDataSources(test *testing.T) {
	ctx := context.TODO()
	dataSourceName := fmt.Sprintf("DS-%d", time.Now().Unix())
	testObject := getTestObject(ctx, test)
	request := &api.AddDataSourcesReqApplicationJSON{}
	params := api.AddDataSourcesParams{
		DataSource: []string{dataSourceName},
	}
	response, err := testObject.AddDataSources(ctx, request, params)
	testError(test, ctx, err)
	switch responseTyped := response.(type) {
	case *api.SzDataSourcesResponse:
		if debug {
			drillDown := []interface{}{
				response,
				responseTyped,
				responseTyped.Data,
				responseTyped.Data.Value,
				responseTyped.Data.Value.DataSourceDetails,
				responseTyped.Data.Value.DataSourceDetails.Value,
				responseTyped.Data.Value.DataSourceDetails.Value["xxxBob"],
				responseTyped.Data.Value.DataSourceDetails.Value["xxxBob"].DataSourceCode,
				responseTyped.Data.Value.DataSourceDetails.Value["xxxBob"].DataSourceCode.Value,
			}

			for index, value := range drillDown {
				test.Logf(">>>>> %d: %-60s %+v\n", index, reflect.TypeOf(value), value)
			}
		}
	}
}

func TestSenzingRestServiceImpl_Heartbeat(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	response, err := testObject.Heartbeat(ctx)
	testError(test, ctx, err)
	httpMethod, err := response.Meta.Value.HttpMethod.Value.MarshalText()
	testError(test, ctx, err)
	assert.Equal(test, "GET", string(httpMethod))
}

func TestSenzingRestServiceImpl_License(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	params := api.LicenseParams{
		WithRaw: api.NewOptBool(false),
	}
	response, err := testObject.License(ctx, params)
	testError(test, ctx, err)
	switch responseTyped := response.(type) {
	case *api.SzLicenseResponse:
		recordLimit, _ := responseTyped.Data.Value.License.Value.RecordLimit.Get()
		assert.Equal(test, int64(50000), recordLimit)
	}
}

func TestSenzingRestServiceImpl_OpenApiSpecification(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	var openApiSpecificationBytes []byte
	response, err := testObject.OpenApiSpecification(ctx)
	testError(test, ctx, err)
	numBytes, _ := response.Data.Read(openApiSpecificationBytes)
	// testError(test, ctx, err)
	test.Logf(">>>>> %d;  %v\n", numBytes, openApiSpecificationBytes)
}

func TestSenzingRestServiceImpl_Version(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	params := api.VersionParams{
		WithRaw: api.NewOptBool(false),
	}
	response, err := testObject.Version(ctx, params)
	testError(test, ctx, err)
	switch responseTyped := response.(type) {
	case *api.SzVersionResponse:
		apiServerVersion, _ := responseTyped.Data.Value.ApiServerVersion.Get()
		assert.Equal(test, "0.0.0", apiServerVersion)
	}
}
