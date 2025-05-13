package senzingrestservice_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	senzingRestServiceSingleton *senzingrestservice.BasicSenzingRestService
	debug                       bool
	logLevel                    = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingRestServiceImpl_AddDataSources(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
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
				test.Logf(">>>>> AddDataSources %d: ; %-60s %+v\n", index, reflect.TypeOf(value), value)
			}
		}
	default:
	}
}

func TestSenzingRestServiceImpl_Heartbeat(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	response, err := testObject.Heartbeat(ctx)
	require.NoError(test, err)
	httpMethod, err := response.Meta.Value.HttpMethod.Value.MarshalText()
	require.NoError(test, err)
	assert.Equal(test, "GET", string(httpMethod))
}

func TestSenzingRestServiceImpl_License(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	params := senzingrestapi.LicenseParams{
		WithRaw: senzingrestapi.NewOptBool(false),
	}
	response, err := testObject.License(ctx, params)
	require.NoError(test, err)

	switch responseTyped := response.(type) {
	case *senzingrestapi.SzLicenseResponse:
		recordLimit, _ := responseTyped.Data.Value.License.Value.RecordLimit.Get()
		assert.Equal(test, int64(50000), recordLimit)
	default:
	}
}

func TestSenzingRestServiceImpl_OpenAPISpecification(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)

	var openAPISpecificationBytes []byte

	response, err := testObject.OpenAPISpecification(ctx)
	require.NoError(test, err)
	_, err = response.Data.Read(openAPISpecificationBytes)
	require.Error(test, err) // An EOF error.
}

func TestSenzingRestServiceImpl_Version(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
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

func getTestObject(ctx context.Context, t *testing.T) senzingrestservice.SenzingRestService {
	t.Helper()

	if senzingRestServiceSingleton == nil {
		settings, err := settings.BuildSimpleSettingsUsingEnvVars()
		require.NoError(t, err)

		senzingRestServiceSingleton = &senzingrestservice.BasicSenzingRestService{
			Settings:              settings,
			SenzingInstanceName:   "go-rest-api-service-test",
			SenzingVerboseLogging: int64(0),
		}
		err = senzingRestServiceSingleton.SetLogLevel(ctx, logLevel)
		require.NoError(t, err)
	}

	return senzingRestServiceSingleton
}
