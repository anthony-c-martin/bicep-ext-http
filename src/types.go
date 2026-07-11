package main

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/bicep-types/src/bicep-types-go/factory"
	"github.com/Azure/bicep-types/src/bicep-types-go/index"
	"github.com/Azure/bicep-types/src/bicep-types-go/types"
	"github.com/Azure/bicep-types/src/bicep-types-go/writers"
)

const (
	resourceTypeName = "HttpRequest"
	typesFileName    = "types.json"
	extensionName    = "Http"
	extensionVersion = "0.0.1"
)

// allScopes preserves the original behaviour where the resource is valid at any scope.
const allScopes = types.AllExceptExtension | types.ScopeTypeExtension

// buildTypeFiles authors the extension's Bicep type definitions using the
// bicep-types-go library. It returns the serialized index.json content along
// with a map of type file paths to their serialized content.
func buildTypeFiles() (indexContent string, typeFiles map[string]string, err error) {
	fac := factory.NewTypeFactory()

	rawRef := fac.GetReference(fac.CreateStringLiteralType("raw"))
	jsonRef := fac.GetReference(fac.CreateStringLiteralType("json"))
	formatRef := fac.GetReference(fac.CreateUnionType([]types.ITypeReference{rawRef, jsonRef}))

	stringRef := fac.GetReference(fac.CreateStringType())
	intRef := fac.GetReference(fac.CreateIntegerType())
	anyRef := fac.GetReference(fac.CreateAnyType())

	bodyRef := fac.GetReference(fac.CreateObjectType(resourceTypeName, map[string]types.ObjectTypeProperty{
		"uri": {
			Type:        stringRef,
			Flags:       types.TypePropertyFlagsRequired,
			Description: "The HTTP request URI to submit a GET request to.",
		},
		"format": {
			Type:        formatRef,
			Flags:       types.TypePropertyFlagsNone,
			Description: "How to deserialize the response body.",
		},
		"statusCode": {
			Type:        intRef,
			Flags:       types.TypePropertyFlagsReadOnly,
			Description: "The status code of the HTTP request.",
		},
		"body": {
			Type:        anyRef,
			Flags:       types.TypePropertyFlagsReadOnly,
			Description: "The parsed request body.",
		},
	}, nil, nil))

	resource := fac.CreateResourceType(resourceTypeName, bodyRef, allScopes, allScopes, nil)
	resourceRef := fac.GetReference(resource).(types.TypeReference)

	typesJSON, err := writers.NewJSONWriter().WriteTypesToString(fac.GetTypes())
	if err != nil {
		return "", nil, fmt.Errorf("failed to serialize types: %w", err)
	}

	// The library's index serializer keys resources by "type@version"; this
	// extension uses an unversioned resource name, so build the index directly.
	indexJSON, err := json.MarshalIndent(struct {
		Resources         map[string]types.ITypeReference `json:"resources"`
		ResourceFunctions map[string]any                  `json:"resourceFunctions"`
		Settings          index.TypeSettings              `json:"settings"`
	}{
		Resources: map[string]types.ITypeReference{
			resourceTypeName: types.CrossFileTypeReference{RelativePath: typesFileName, Ref: resourceRef.Ref},
		},
		ResourceFunctions: map[string]any{},
		Settings:          index.TypeSettings{Name: extensionName, Version: extensionVersion},
	}, "", "  ")
	if err != nil {
		return "", nil, fmt.Errorf("failed to serialize index: %w", err)
	}

	return string(indexJSON), map[string]string{typesFileName: typesJSON}, nil
}
