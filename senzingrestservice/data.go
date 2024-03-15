package senzingrestservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/senzing-garage/g2-sdk-go/senzing"
	api "github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingresttypedef"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

func licenseData(ctx context.Context, license string) api.OptSzLicenseResponseData {
	parsedLicense, err := senzing.UnmarshalG2productLicenseResponse(ctx, license)
	if err != nil {
		panic(err)
	}
	issueDate, err := time.Parse("2006-01-02", parsedLicense.IssueDate)
	if err != nil {
		panic(err)
	}
	expireDate, err := time.Parse("2006-01-02", parsedLicense.ExpireDate)
	if err != nil {
		panic(err)
	}
	return api.OptSzLicenseResponseData{
		Set: true,
		Value: api.SzLicenseResponseData{
			License: api.OptSzLicenseInfo{
				Set: true,
				Value: api.SzLicenseInfo{
					Customer:       api.NewOptString(parsedLicense.Customer),
					Contract:       api.NewOptString(parsedLicense.Contract),
					LicenseType:    api.NewOptString(parsedLicense.LicenseType),
					LicenseLevel:   api.NewOptString(parsedLicense.LicenseLevel),
					Billing:        api.NewOptString(parsedLicense.Billing),
					IssuanceDate:   api.NewOptDateTime(issueDate),
					ExpirationDate: api.NewOptDateTime(expireDate),
					RecordLimit:    api.NewOptInt64(parsedLicense.RecordLimit),
				},
			},
		},
	}
}

func versionData(ctx context.Context, version string) (r api.OptSzVersionInfo) {
	senzingVersion, err := senzing.UnmarshalG2productVersionResponse(ctx, version)
	if err != nil {
		panic(err)
	}
	nativeApiBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	if err != nil {
		panic(err)
	}
	return api.OptSzVersionInfo{
		Set: true,
		Value: api.SzVersionInfo{
			ApiServerVersion:           api.NewOptString(ApiServerVersion),
			RestApiVersion:             api.NewOptString(RestApiVersion),
			NativeApiVersion:           api.NewOptString(senzingVersion.Version),
			NativeApiBuildVersion:      api.NewOptString(senzingVersion.BuildVersion),
			NativeApiBuildNumber:       api.NewOptString(senzingVersion.BuildVersion),
			NativeApiBuildDate:         api.NewOptDateTime(nativeApiBuildDate),
			ConfigCompatibilityVersion: api.NewOptString(senzingVersion.CompatibilityVersion.ConfigVersion),
		},
	}
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

func (restApiService *SenzingRestServiceImpl) getSzMeta(ctx context.Context, httpMethod api.SzHttpMethod, httpStatusCode int16) api.SzMeta {
	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}
	nativeApiBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	if err != nil {
		panic(err)
	}
	return api.SzMeta{
		Server:                     api.NewOptString("Senzing REST API Server - go"),
		HttpMethod:                 api.NewOptSzHttpMethod(httpMethod),
		HttpStatusCode:             api.NewOptInt16(httpStatusCode),
		Timestamp:                  api.NewOptDateTime(time.Now().UTC()),
		Version:                    api.NewOptString("0.0.0"),
		RestApiVersion:             api.NewOptString("3.4.1"),
		NativeApiVersion:           api.NewOptString(senzingVersion.Version),
		NativeApiBuildVersion:      api.NewOptString(senzingVersion.BuildVersion),
		NativeApiBuildNumber:       api.NewOptString(senzingVersion.BuildNumber),
		NativeApiBuildDate:         api.NewOptDateTime(nativeApiBuildDate),
		ConfigCompatibilityVersion: api.NewOptString(senzingVersion.CompatibilityVersion.ConfigVersion),
		Timings:                    api.NewOptNilSzMetaTimings(map[string]int64{}),
	}
}

func (restApiService *SenzingRestServiceImpl) getMeta(ctx context.Context, httpMethod string, httpStatusCode int32) senzingresttypedef.Meta {
	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}
	nativeApiBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	if err != nil {
		panic(err)
	}
	return senzingresttypedef.Meta{
		Server:                     "Senzing REST API Server - go",
		HTTPMethod:                 httpMethod,
		HTTPStatusCode:             httpStatusCode,
		Timestamp:                  time.Now().UTC().String(),
		Version:                    "0.0.0",
		RestAPIVersion:             "3.4.1",
		NativeAPIVersion:           senzingVersion.Version,
		NativeAPIBuildVersion:      senzingVersion.BuildVersion,
		NativeAPIBuildNumber:       senzingVersion.BuildNumber,
		NativeAPIBuildDate:         nativeApiBuildDate.String(),
		ConfigCompatibilityVersion: senzingVersion.CompatibilityVersion.ConfigVersion,
		Timings:                    senzingresttypedef.Timings{},
	}
}

func (restApiService *SenzingRestServiceImpl) getOptSzLinks(ctx context.Context, uriPath string) api.OptSzLinks {
	var result api.OptSzLinks
	szLinks := api.SzLinks{
		Self:                 api.NewOptString(fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.UrlRoutePrefix, uriPath)),
		OpenApiSpecification: api.NewOptString(fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.UrlRoutePrefix)),
	}
	result = api.NewOptSzLinks(szLinks)
	return result
}

func (restApiService *SenzingRestServiceImpl) getLinks(ctx context.Context, uriPath string) senzingresttypedef.Links {
	return senzingresttypedef.Links{
		Self:                 fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.UrlRoutePrefix, uriPath),
		OpenAPISpecification: fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.UrlRoutePrefix),
	}
}

func (restApiService *SenzingRestServiceImpl) getOptSzMeta(ctx context.Context, httpMethod api.SzHttpMethod, httpStatusCode int16) api.OptSzMeta {
	return api.NewOptSzMeta(restApiService.getSzMeta(ctx, httpMethod, httpStatusCode))
}

func (restApiService *SenzingRestServiceImpl) getOptSzServerInfo(ctx context.Context) api.OptSzServerInfo {
	_ = ctx
	return api.OptSzServerInfo{
		Set: true,
		Value: api.SzServerInfo{
			Concurrency:              api.NewOptInt32(1),
			ActiveConfigId:           api.NewOptInt32(1),
			DynamicConfig:            api.NewOptBool(true),
			ReadOnly:                 api.NewOptBool(true),
			AdminEnabled:             api.NewOptBool(true),
			WebSocketsMessageMaxSize: api.NewOptInt32(100),
			InfoQueueConfigured:      api.NewOptBool(true),
		},
	}
}

func (restApiService *SenzingRestServiceImpl) openApiSpecificationData(ctx context.Context) (r io.Reader) {
	_ = ctx
	return bytes.NewReader(restApiService.OpenApiSpecificationSpec)
}

func (restApiService *SenzingRestServiceImpl) searchEntitiesByGetData(ctx context.Context, input string) io.Reader {
	_ = input

	data := senzingresttypedef.SearchResultsData{
		SearchResults: []senzingresttypedef.SearchResult{senzingresttypedef.SearchResult{AddressData: []string{"Bob was here"}}},
	}

	// append(data.SearchResults, senzingresttypedef.SearchResult{AddressData: []string{"Bob was here"}})

	resultStruct := senzingresttypedef.SearchEntitiesByGetResponse{
		Meta:  restApiService.getMeta(ctx, "GET", http.StatusOK),
		Links: restApiService.getLinks(ctx, "license"),
		Data:  data,
	}

	resultBytes, err := json.Marshal(resultStruct)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(resultBytes)
	// return strings.NewReader(output)
}
