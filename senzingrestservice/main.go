package senzingrestservice

import (
	_ "embed"

	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The SenzingRestService interface is...
type SenzingRestService interface {
	senzingrestapi.Handler
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6503xxxx".
const ComponentID = 9999

// Log message prefix.
const Prefix = "go-rest-api-service.senzingrestservice."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for szconfig implementations.
var IDMessages = map[int]string{
	10: "Enter " + Prefix + "InitializeSenzing().",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

//go:embed openapi.json
var OpenAPISpecificationJSON []byte

//go:embed openapi.yaml
var OpenAPISpecificationYaml []byte
