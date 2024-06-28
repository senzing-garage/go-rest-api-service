package senzingrestservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingresttypedef"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// func licenseData(ctx context.Context, license string) senzingrestapi.OptSzLicenseResponseData {
// 	parsedLicense, err := response.SzProductGetLicense(ctx, license)
// 	if err != nil {
// 		panic(err)
// 	}
// 	issueDate, err := time.Parse("2006-01-02", parsedLicense.IssueDate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	expireDate, err := time.Parse("2006-01-02", parsedLicense.ExpireDate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return senzingrestapi.OptSzLicenseResponseData{
// 		Set: true,
// 		Value: senzingrestapi.SzLicenseResponseData{
// 			License: senzingrestapi.OptSzLicenseInfo{
// 				Set: true,
// 				Value: senzingrestapi.SzLicenseInfo{
// 					Customer:       senzingrestapi.NewOptString(parsedLicense.Customer),
// 					Contract:       senzingrestapi.NewOptString(parsedLicense.Contract),
// 					LicenseType:    senzingrestapi.NewOptString(parsedLicense.LicenseType),
// 					LicenseLevel:   senzingrestapi.NewOptString(parsedLicense.LicenseLevel),
// 					Billing:        senzingrestapi.NewOptString(parsedLicense.Billing),
// 					IssuanceDate:   senzingrestapi.NewOptDateTime(issueDate),
// 					ExpirationDate: senzingrestapi.NewOptDateTime(expireDate),
// 					RecordLimit:    senzingrestapi.NewOptInt64(parsedLicense.RecordLimit),
// 				},
// 			},
// 		},
// 	}
// }

// func versionData(ctx context.Context, version string) (r senzingrestapi.OptSzVersionInfo) {
// 	senzingVersion, err := response.SzProductGetVersion(ctx, version)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nativeAPIBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return senzingrestapi.OptSzVersionInfo{
// 		Set: true,
// 		Value: senzingrestapi.SzVersionInfo{
// 			ApiServerVersion:           senzingrestapi.NewOptString(APIServerVersion),
// 			RestApiVersion:             senzingrestapi.NewOptString(RestAPIVersion),
// 			NativeApiVersion:           senzingrestapi.NewOptString(senzingVersion.Version),
// 			NativeApiBuildVersion:      senzingrestapi.NewOptString(senzingVersion.BuildVersion),
// 			NativeApiBuildNumber:       senzingrestapi.NewOptString(senzingVersion.BuildVersion),
// 			NativeApiBuildDate:         senzingrestapi.NewOptDateTime(nativeAPIBuildDate),
// 			ConfigCompatibilityVersion: senzingrestapi.NewOptString(senzingVersion.CompatibilityVersion.ConfigVersion),
// 		},
// 	}
// }

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

func (restApiService *BasicSenzingRestService) getSzMeta(ctx context.Context, httpMethod senzingrestapi.SzHttpMethod, httpStatusCode int16) senzingrestapi.SzMeta {
	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}
	nativeAPIBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
	if err != nil {
		panic(err)
	}
	return senzingrestapi.SzMeta{
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
}

func (restApiService *BasicSenzingRestService) getMeta(ctx context.Context, httpMethod string, httpStatusCode int32) senzingresttypedef.Meta {
	senzingVersion, err := restApiService.getSenzingVersion(ctx)
	if err != nil {
		panic(err)
	}
	nativeAPIBuildDate, err := time.Parse("2006-01-02", senzingVersion.BuildDate)
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
		NativeAPIBuildDate:         nativeAPIBuildDate.String(),
		ConfigCompatibilityVersion: senzingVersion.CompatibilityVersion.ConfigVersion,
		Timings:                    senzingresttypedef.Timings{},
	}
}

func (restApiService *BasicSenzingRestService) getOptSzLinks(ctx context.Context, uriPath string) senzingrestapi.OptSzLinks {
	var result senzingrestapi.OptSzLinks
	szLinks := senzingrestapi.SzLinks{
		Self:                 senzingrestapi.NewOptString(fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.URLRoutePrefix, uriPath)),
		OpenApiSpecification: senzingrestapi.NewOptString(fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.URLRoutePrefix)),
	}
	result = senzingrestapi.NewOptSzLinks(szLinks)
	return result
}

func (restApiService *BasicSenzingRestService) getLinks(ctx context.Context, uriPath string) senzingresttypedef.Links {
	return senzingresttypedef.Links{
		Self:                 fmt.Sprintf("http://%s/%s/%s", getHostname(ctx), restApiService.URLRoutePrefix, uriPath),
		OpenAPISpecification: fmt.Sprintf("http://%s/%s/swagger_spec", getHostname(ctx), restApiService.URLRoutePrefix),
	}
}

func (restApiService *BasicSenzingRestService) getOptSzServerInfo(ctx context.Context) senzingrestapi.OptSzServerInfo {
	_ = ctx
	return senzingrestapi.OptSzServerInfo{
		Set: true,
		Value: senzingrestapi.SzServerInfo{
			Concurrency:              senzingrestapi.NewOptInt32(1),
			ActiveConfigId:           senzingrestapi.NewOptInt32(1),
			DynamicConfig:            senzingrestapi.NewOptBool(true),
			ReadOnly:                 senzingrestapi.NewOptBool(true),
			AdminEnabled:             senzingrestapi.NewOptBool(true),
			WebSocketsMessageMaxSize: senzingrestapi.NewOptInt32(100),
			InfoQueueConfigured:      senzingrestapi.NewOptBool(true),
		},
	}
}

// func (restApiService *BasicSenzingRestService) openAPISpecificationData(ctx context.Context) (r io.Reader) {
// 	_ = ctx
// 	return bytes.NewReader(restApiService.OpenAPISpecificationSpec)
// }

func (restApiService *BasicSenzingRestService) searchEntitiesByGetData(ctx context.Context, input string) io.Reader {
	_ = input

	data := senzingresttypedef.SearchResultsData{
		SearchResults: []senzingresttypedef.SearchResult{{AddressData: []string{"Bob was here"}}},
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
