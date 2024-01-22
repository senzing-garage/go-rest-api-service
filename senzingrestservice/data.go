package senzingrestservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	api "github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingresttypedef"
	"github.com/senzing/g2-sdk-go/senzing"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

func licenseData(ctx context.Context, license string) api.OptSzLicenseResponseData {
	parsedLicense, err := senzing.UnmarshalProductLicenseResponse(ctx, license)
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

func searchEntitiesByGetDataXX(ctx context.Context, input string) io.Reader {

	output := `{"meta":{"server":"Senzing POC API Server","httpMethod":"GET","httpStatusCode":200,"timings":{"overall":2215,"enqueued":1,"nativeAPI":2189,"nativeAPI:engine.searchByAttributes":2189,"processRawData":23},"timestamp":"2024-01-19T16:59:32.660Z","version":"","restApiVersion":"3.4.1","nativeApiVersion":"3.8.0","nativeApiBuildVersion":"3.8.0.23303","nativeApiBuildNumber":"2023_10_30__10_45","nativeApiBuildDate":"2023-10-30T17:45:00.000Z","configCompatibilityVersion":"10","pocServerVersion":"3.5.1-dirty","pocApiVersion":"3.5.1"},"links":{"self":"http://localhost:8251/entities?attrs=%7B%22NAME_FULL%22%3A%22robert%20smith%22%2C%22NAME_TYPE%22%3A%22PRIMARY%22%2C%22COMPANY_NAME_ORG%22%3A%22robert%20smith%22%7D","openApiSpecification":"http://localhost:8251/specifications/open-api"},"data":{"searchResults":[{"entityId":1,"entityName":"Robert Smith","recordSummaries":[{"dataSource":"CUSTOMERS","recordCount":4,"topRecordIds":["1001","1002","1003","1004"]}],"addressData":["MAILING: 123 Main Street, Las Vegas NV 89132","HOME: 1515 Adela Lane Las Vegas NV 89111"],"characteristicData":["RECORD_TYPE: PERSON","DOB: 11/12/1979","DOB: 12/11/1978"],"identifierData":["EMAIL: bsmith@work.com"],"nameData":["PRIMARY: Robert Smith","PRIMARY: B Smith"],"phoneData":["MOBILE: 702-919-1300","HOME: 702-919-1300"],"otherData":["DATE: 1/5/15","STATUS: Inactive","AMOUNT: 400","DATE: 4/9/16","AMOUNT: 300","DATE: 1/2/18","STATUS: Active","AMOUNT: 100","DATE: 3/10/17","AMOUNT: 200"],"features":{"ADDRESS":[{"primaryId":40,"primaryValue":"1515 Adela Lane Las Vegas NV 89111","usageType":"HOME","duplicateValues":["1515 Adela Ln Las Vegas NV 89132"],"featureDetails":[{"internalId":40,"featureValue":"1515 Adela Lane Las Vegas NV 89111"},{"internalId":20,"featureValue":"1515 Adela Ln Las Vegas NV 89132"}]},{"primaryId":53,"primaryValue":"123 Main Street, Las Vegas NV 89132","usageType":"MAILING","featureDetails":[{"internalId":53,"featureValue":"123 Main Street, Las Vegas NV 89132"}]}],"DOB":[{"primaryId":19,"primaryValue":"11/12/1979","featureDetails":[{"internalId":19,"featureValue":"11/12/1979"}]},{"primaryId":2,"primaryValue":"12/11/1978","duplicateValues":["11/12/1978"],"featureDetails":[{"internalId":2,"featureValue":"12/11/1978"},{"internalId":39,"featureValue":"11/12/1978"}]}],"EMAIL":[{"primaryId":3,"primaryValue":"bsmith@work.com","featureDetails":[{"internalId":3,"featureValue":"bsmith@work.com"}]}],"NAME":[{"primaryId":18,"primaryValue":"B Smith","usageType":"PRIMARY","featureDetails":[{"internalId":18,"featureValue":"B Smith"}]},{"primaryId":52,"primaryValue":"Robert Smith","usageType":"PRIMARY","duplicateValues":["Bob J Smith","Bob Smith"],"featureDetails":[{"internalId":52,"featureValue":"Robert Smith"},{"internalId":1,"featureValue":"Bob J Smith"},{"internalId":38,"featureValue":"Bob Smith"}]}],"PHONE":[{"primaryId":41,"primaryValue":"702-919-1300","usageType":"HOME","featureDetails":[{"internalId":41,"featureValue":"702-919-1300"}]},{"primaryId":41,"primaryValue":"702-919-1300","usageType":"MOBILE","featureDetails":[{"internalId":41,"featureValue":"702-919-1300"}]}],"RECORD_TYPE":[{"primaryId":16,"primaryValue":"PERSON","featureDetails":[{"internalId":16,"featureValue":"PERSON"}]}]},"records":[{"dataSource":"CUSTOMERS","recordId":"1004","addressData":["HOME: 1515 Adela Ln Las Vegas NV 89132"],"characteristicData":["DOB: 11/12/1979","RECORD_TYPE: PERSON"],"identifierData":["EMAIL: bsmith@work.com"],"nameData":["PRIMARY: Smith B"],"otherData":["DATE: 1/5/15","STATUS: Inactive","AMOUNT: 400"],"originalSourceData":{"ADDR_CITY":"Las Vegas","ADDR_LINE1":"1515 Adela Ln","ADDR_POSTAL_CODE":"89132","ADDR_STATE":"NV","ADDR_TYPE":"HOME","AMOUNT":"400","DATA_SOURCE":"CUSTOMERS","DATE":"1/5/15","DATE_OF_BIRTH":"11/12/1979","EMAIL_ADDRESS":"bsmith@work.com","PRIMARY_NAME_FIRST":"B","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1004","RECORD_TYPE":"PERSON","STATUS":"Inactive"},"lastSeenTimestamp":"2024-01-19T16:52:47.007Z","matchLevel":0},{"dataSource":"CUSTOMERS","recordId":"1003","characteristicData":["DOB: 12/11/1978","RECORD_TYPE: PERSON"],"identifierData":["EMAIL: bsmith@work.com"],"nameData":["PRIMARY: Smith Bob J"],"otherData":["DATE: 4/9/16","STATUS: Inactive","AMOUNT: 300"],"originalSourceData":{"AMOUNT":"300","DATA_SOURCE":"CUSTOMERS","DATE":"4/9/16","DATE_OF_BIRTH":"12/11/1978","EMAIL_ADDRESS":"bsmith@work.com","PRIMARY_NAME_FIRST":"Bob","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_MIDDLE":"J","RECORD_ID":"1003","RECORD_TYPE":"PERSON","STATUS":"Inactive"},"lastSeenTimestamp":"2024-01-19T16:52:47.005Z","matchLevel":1,"matchKey":"+NAME+DOB+EMAIL","resolutionRuleCode":"SF1_PNAME_CSTAB"},{"dataSource":"CUSTOMERS","recordId":"1001","addressData":["MAILING: 123 Main Street, Las Vegas NV 89132"],"characteristicData":["DOB: 12/11/1978","RECORD_TYPE: PERSON"],"identifierData":["EMAIL: bsmith@work.com"],"nameData":["PRIMARY: Smith Robert"],"phoneData":["HOME: 702-919-1300"],"otherData":["DATE: 1/2/18","STATUS: Active","AMOUNT: 100"],"originalSourceData":{"ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","ADDR_TYPE":"MAILING","AMOUNT":"100","DATA_SOURCE":"CUSTOMERS","DATE":"1/2/18","DATE_OF_BIRTH":"12/11/1978","EMAIL_ADDRESS":"bsmith@work.com","PHONE_NUMBER":"702-919-1300","PHONE_TYPE":"HOME","PRIMARY_NAME_FIRST":"Robert","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1001","RECORD_TYPE":"PERSON","STATUS":"Active"},"lastSeenTimestamp":"2024-01-19T16:52:51.928Z","matchLevel":1,"matchKey":"+NAME+DOB+PHONE+EMAIL","resolutionRuleCode":"SF1_SNAME_CFF_CSTAB"},{"dataSource":"CUSTOMERS","recordId":"1002","addressData":["HOME: 1515 Adela Lane Las Vegas NV 89111"],"characteristicData":["DOB: 11/12/1978","RECORD_TYPE: PERSON"],"nameData":["PRIMARY: Smith Bob"],"phoneData":["MOBILE: 702-919-1300"],"otherData":["DATE: 3/10/17","STATUS: Inactive","AMOUNT: 200"],"originalSourceData":{"ADDR_CITY":"Las Vegas","ADDR_LINE1":"1515 Adela Lane","ADDR_POSTAL_CODE":"89111","ADDR_STATE":"NV","ADDR_TYPE":"HOME","AMOUNT":"200","DATA_SOURCE":"CUSTOMERS","DATE":"3/10/17","DATE_OF_BIRTH":"11/12/1978","PHONE_NUMBER":"702-919-1300","PHONE_TYPE":"MOBILE","PRIMARY_NAME_FIRST":"Bob","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1002","RECORD_TYPE":"PERSON","STATUS":"Inactive"},"lastSeenTimestamp":"2024-01-19T16:52:47.005Z","matchLevel":1,"matchKey":"+NAME+DOB+ADDRESS","resolutionRuleCode":"CNAME_CFF_CEXCL"}],"partial":false,"lastSeenTimestamp":"2024-01-19T16:52:51.928Z","matchLevel":2,"matchKey":"+NAME","resolutionRuleCode":"CNAME","resultType":"POSSIBLE_MATCH","bestNameScore":100,"featureScores":{"NAME":[{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robert Smith","score":100,"nameScoringDetails":{"fullNameScore":100}}]}},{"entityId":48,"entityName":"Robert E Smith Sr","recordSummaries":[{"dataSource":"CUSTOMERS","recordCount":1,"topRecordIds":["1005"]},{"dataSource":"WATCHLIST","recordCount":1,"topRecordIds":["1006"]}],"addressData":["MAILING: 123 E Main St Henderson NV 89132","MAILING: 123 Main St, Las Vegas"],"characteristicData":["RECORD_TYPE: PERSON","DOB: 3/31/1954"],"identifierData":["DRLIC: 112233 NV"],"nameData":["PRIMARY: Robert E Smith Sr"],"otherData":["DATE: 7/16/19","STATUS: Active","AMOUNT: 500","DATE: 1/3/17","CATEGORY: Fraud"],"features":{"ADDRESS":[{"primaryId":382,"primaryValue":"123 E Main St Henderson NV 89132","usageType":"MAILING","featureDetails":[{"internalId":382,"featureValue":"123 E Main St Henderson NV 89132"}]},{"primaryId":1010,"primaryValue":"123 Main St, Las Vegas","usageType":"MAILING","featureDetails":[{"internalId":1010,"featureValue":"123 Main St, Las Vegas"}]}],"DOB":[{"primaryId":1009,"primaryValue":"3/31/1954","featureDetails":[{"internalId":1009,"featureValue":"3/31/1954"}]}],"DRLIC":[{"primaryId":383,"primaryValue":"112233 NV","featureDetails":[{"internalId":383,"featureValue":"112233 NV"}]}],"NAME":[{"primaryId":1008,"primaryValue":"Robert E Smith Sr","usageType":"PRIMARY","duplicateValues":["Robbie Smith"],"featureDetails":[{"internalId":1008,"featureValue":"Robert E Smith Sr"},{"internalId":381,"featureValue":"Robbie Smith"}]}],"RECORD_TYPE":[{"primaryId":16,"primaryValue":"PERSON","featureDetails":[{"internalId":16,"featureValue":"PERSON"}]}]},"records":[{"dataSource":"CUSTOMERS","recordId":"1005","addressData":["MAILING: 123 E Main St Henderson NV 89132"],"characteristicData":["RECORD_TYPE: PERSON"],"identifierData":["DRLIC: 112233 NV"],"nameData":["PRIMARY: Smith Robbie"],"otherData":["DATE: 7/16/19","STATUS: Active","AMOUNT: 500"],"originalSourceData":{"ADDR_CITY":"Henderson","ADDR_LINE1":"123 E Main St","ADDR_POSTAL_CODE":"89132","ADDR_STATE":"NV","ADDR_TYPE":"MAILING","AMOUNT":"500","DATA_SOURCE":"CUSTOMERS","DATE":"7/16/19","DRIVERS_LICENSE_NUMBER":"112233","DRIVERS_LICENSE_STATE":"NV","PRIMARY_NAME_FIRST":"Robbie","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1005","RECORD_TYPE":"PERSON","STATUS":"Active"},"lastSeenTimestamp":"2024-01-19T16:52:51.973Z","matchLevel":0},{"dataSource":"WATCHLIST","recordId":"1006","addressData":["MAILING: 123 Main St, Las Vegas"],"characteristicData":["DOB: 3/31/1954","RECORD_TYPE: PERSON"],"identifierData":["DRLIC: 112233 NV"],"nameData":["PRIMARY: Smith Sr Robert E"],"otherData":["STATUS: Active","DATE: 1/3/17","CATEGORY: Fraud"],"originalSourceData":{"ADDR_LINE1":"123 Main St, Las Vegas ","ADDR_TYPE":"MAILING","CATEGORY":"Fraud","DATA_SOURCE":"WATCHLIST","DATE":"1/3/17","DATE_OF_BIRTH":"3/31/1954","DRIVERS_LICENSE_NUMBER":"112233","DRIVERS_LICENSE_STATE":"NV","PRIMARY_NAME_FIRST":"Robert","PRIMARY_NAME_LAST":"Smith Sr","PRIMARY_NAME_MIDDLE":"E","RECORD_ID":"1006","RECORD_TYPE":"PERSON","STATUS":"Active"},"lastSeenTimestamp":"2024-01-19T16:52:52.502Z","matchLevel":1,"matchKey":"+NAME+DRLIC","resolutionRuleCode":"SF1_CNAME"}],"partial":false,"lastSeenTimestamp":"2024-01-19T16:52:52.502Z","matchLevel":2,"matchKey":"+NAME","resolutionRuleCode":"CNAME","resultType":"POSSIBLE_MATCH","bestNameScore":98,"featureScores":{"NAME":[{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robert E Smith Sr","score":93,"nameScoringDetails":{"fullNameScore":93}},{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robbie Smith","score":97,"nameScoringDetails":{"fullNameScore":97}},{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robert E Smith Sr","score":98,"nameScoringDetails":{"fullNameScore":98,"orgNameScore":98}},{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robbie Smith","score":82,"nameScoringDetails":{"fullNameScore":82,"orgNameScore":82}}]}},{"entityId":107,"entityName":"Robert Smith","recordSummaries":[{"dataSource":"WATCHLIST","recordCount":1,"topRecordIds":["1008"]}],"characteristicData":["RECORD_TYPE: PERSON"],"identifierData":["EMAIL: robert.smith@email.com"],"nameData":["PRIMARY: Robert Smith"],"otherData":["STATUS: Active","DATE: 3/5/19","CATEGORY: Fraud"],"features":{"EMAIL":[{"primaryId":802,"primaryValue":"robert.smith@email.com","featureDetails":[{"internalId":802,"featureValue":"robert.smith@email.com"}]}],"NAME":[{"primaryId":52,"primaryValue":"Robert Smith","usageType":"PRIMARY","featureDetails":[{"internalId":52,"featureValue":"Robert Smith"}]}],"RECORD_TYPE":[{"primaryId":16,"primaryValue":"PERSON","featureDetails":[{"internalId":16,"featureValue":"PERSON"}]}]},"records":[{"dataSource":"WATCHLIST","recordId":"1008","characteristicData":["RECORD_TYPE: PERSON"],"identifierData":["EMAIL: robert.smith@email.com"],"nameData":["PRIMARY: Smith Robert"],"otherData":["STATUS: Active","DATE: 3/5/19","CATEGORY: Fraud"],"originalSourceData":{"CATEGORY":"Fraud","DATA_SOURCE":"WATCHLIST","DATE":"3/5/19","EMAIL_ADDRESS":"robert.smith@email.com","PRIMARY_NAME_FIRST":"Robert","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1008","RECORD_TYPE":"PERSON","STATUS":"Active"},"lastSeenTimestamp":"2024-01-19T16:52:52.306Z","matchLevel":0}],"partial":false,"lastSeenTimestamp":"2024-01-19T16:52:52.306Z","matchLevel":2,"matchKey":"+NAME","resolutionRuleCode":"CNAME","resultType":"POSSIBLE_MATCH","bestNameScore":100,"featureScores":{"NAME":[{"featureType":"NAME","inboundFeature":"robert smith","candidateFeature":"Robert Smith","score":100,"nameScoringDetails":{"fullNameScore":100}}]}}]}}`

	return strings.NewReader(output)
}

func versionData(ctx context.Context, version string) (r api.OptSzVersionInfo) {
	senzingVersion, err := senzing.UnmarshalProductVersionResponse(ctx, version)
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
	return bytes.NewReader(restApiService.OpenApiSpecificationSpec)
}

func (restApiService *SenzingRestServiceImpl) searchEntitiesByGetData(ctx context.Context, input string) io.Reader {

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
