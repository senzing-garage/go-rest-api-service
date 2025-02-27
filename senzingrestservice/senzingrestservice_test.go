package senzingrestservice

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	senzingRestServiceSingleton *BasicSenzingRestService
	debug                       bool
	logLevel                    = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingRestServiceImpl_AddDataSources(test *testing.T) {
	ctx := context.TODO()
	dataSourceName := fmt.Sprintf("DS-%d", time.Now().Unix())
	testObject := getTestObject(ctx, test)
	request := &senzingrestapi.AddDataSourcesReqApplicationJSON{}
	params := senzingrestapi.AddDataSourcesParams{
		DataSource: []string{dataSourceName},
	}
	response, err := testObject.AddDataSources(ctx, request, params)
	require.NoError(test, err)
	switch responseTyped := response.(type) {
	case *senzingrestapi.SzDataSourcesResponse:
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
	default:
	}
}

func TestSenzingRestServiceImpl_Heartbeat(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	response, err := testObject.Heartbeat(ctx)
	require.NoError(test, err)
	httpMethod, err := response.Meta.Value.HttpMethod.Value.MarshalText()
	require.NoError(test, err)
	assert.Equal(test, "GET", string(httpMethod))
}

func TestSenzingRestServiceImpl_License(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	params := senzingrestapi.LicenseParams{
		WithRaw: senzingrestapi.NewOptBool(false),
	}
	response, err := testObject.License(ctx, params)
	require.NoError(test, err)
	switch responseTyped := response.(type) {
	case *senzingrestapi.SzLicenseResponse:
		recordLimit, _ := responseTyped.Data.Value.License.Value.RecordLimit.Get()
		assert.Equal(test, int64(500), recordLimit)
	default:
	}
}

func TestSenzingRestServiceImpl_OpenAPISpecification(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	var openAPISpecificationBytes []byte
	response, err := testObject.OpenAPISpecification(ctx)
	require.NoError(test, err)
	numBytes, _ := response.Data.Read(openAPISpecificationBytes)
	require.NoError(test, err)

	// testError(test, ctx, err)
	test.Logf(">>>>> %d;  %v\n", numBytes, openAPISpecificationBytes)
}

func TestSenzingRestServiceImpl_Version(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	params := senzingrestapi.VersionParams{
		WithRaw: senzingrestapi.NewOptBool(false),
	}
	response, err := testObject.Version(ctx, params)
	require.NoError(test, err)
	switch responseTyped := response.(type) {
	case *senzingrestapi.SzVersionResponse:
		apiServerVersion, _ := responseTyped.Data.Value.ApiServerVersion.Get()
		assert.Equal(test, "0.0.0", apiServerVersion)
	default:
	}
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) SenzingRestService {
	if senzingRestServiceSingleton == nil {
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		require.NoError(test, err)
		senzingRestServiceSingleton = &BasicSenzingRestService{
			Settings:              settings,
			SenzingInstanceName:   "go-rest-api-service-test",
			SenzingVerboseLogging: int64(0),
		}
		err = senzingRestServiceSingleton.SetLogLevel(ctx, logLevel)
		require.NoError(test, err)
	}
	return senzingRestServiceSingleton
}
