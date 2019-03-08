package hlsp

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// Parse receives either a YAML or JSON AsyncAPI document, and tries to parse it.
func Parse(yamlOrJSONDocument []byte) (json.RawMessage, *ParserError) {
	jsonDocument, err := yaml.YAMLToJSON(yamlOrJSONDocument)
	if err != nil {
		return nil, &ParserError{
			errorMessage: err.Error(),
		}
	}
	return ParseJSON(jsonDocument)
}

// ParseJSON receives a JSON AsyncAPI document.
// It parses the document and checks if it's valid AsyncAPI.
// Skips specification extensions and schemas validation.
// If validation fails, the Parser/Validator should trigger an error.
// Produces a beautified version of the document in JSON Schema Draft 07.
func ParseJSON(jsonDocument []byte) (json.RawMessage, *ParserError) {
	schemaLoader := gojsonschema.NewReferenceLoader("file://../asyncapi/2.0.0/schema.json")
	documentLoader := gojsonschema.NewBytesLoader(jsonDocument)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, &ParserError{
			errorMessage: err.Error(),
		}
	}

	if result.Valid() {
		beautifiedDoc, err := Beautify(jsonDocument)
		if err != nil {
			panic(err)
		}
		return beautifiedDoc, nil
	}

	return jsonDocument, &ParserError{
		errorMessage:  "[Invalid AsyncAPI document] Check out err.ParsingErrors() for more information.",
		parsingErrors: result.Errors(),
	}
}
