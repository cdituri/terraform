// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"encoding/json"

	"github.com/kaptinlin/jsonschema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

var schemaCompiler = jsonschema.NewCompiler()

// JSONSchemaFunc validates a JSON instance against a given JSON Schema.
var JSONSchemaFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:        "schema",
			Description: "UTF-8 encoded JSON Schema document.",
			Type:        cty.String,
		},
		{
			Name:        "instance",
			Description: "UTF-8 encoded JSON document.",
			Type:        cty.String,
		},
	},
	Type:         function.StaticReturnType(cty.Bool),
	RefineResult: refineNotNull,
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		var schemaJSON string
		if err := gocty.FromCtyValue(args[0], &schemaJSON); err != nil {
			return cty.UnknownVal(cty.String), function.NewArgError(0, err)
		}

		var instanceJSON string
		if err := gocty.FromCtyValue(args[1], &instanceJSON); err != nil {
			return cty.UnknownVal(cty.String), function.NewArgError(1, err)
		}

		schema, err := schemaCompiler.Compile([]byte(schemaJSON))
		if err != nil {
			return cty.UnknownVal(cty.String), function.NewArgErrorf(
				0,
				"schema must be a valid UTF-8 encoded JSON Schema document: %s",
				err,
			)
		}

		var instance map[string]interface{}
		if err = json.Unmarshal([]byte(instanceJSON), &instance); err != nil {
			return cty.UnknownVal(cty.String), function.NewArgErrorf(
				1,
				"instance must be a valid UTF-8 encoded JSON document: %s",
				err,
			)
		}

		result := schema.Validate(instance)

		return cty.BoolVal(result.IsValid()), nil
	},
})

// JSONSchema validates a JSON instance against a given JSON Schema.
func JSONSchema(schema cty.Value, instance cty.Value) (cty.Value, error) {
	return JSONSchemaFunc.Call([]cty.Value{schema, instance})
}
