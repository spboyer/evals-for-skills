package jsonrpc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spboyer/waza/schemas"
)

// evalSchema is the compiled JSON Schema for eval.yaml files.
// Compiled once at package init time so every validation call reuses it.
var evalSchema *jsonschema.Schema

// taskSchema is the compiled JSON Schema for task YAML files.
var taskSchema *jsonschema.Schema

func init() {
	evalSchema = mustCompileSchema(schemas.EvalSchemaJSON, "eval.schema.json")
	taskSchema = mustCompileSchema(schemas.TaskSchemaJSON, "task.schema.json")
}

func mustCompileSchema(raw string, name string) *jsonschema.Schema {
	var schemaDoc any
	if err := json.Unmarshal([]byte(raw), &schemaDoc); err != nil {
		panic(fmt.Sprintf("failed to parse embedded %s: %v", name, err))
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(name, schemaDoc); err != nil {
		panic(fmt.Sprintf("failed to add %s resource: %v", name, err))
	}

	sch, err := compiler.Compile(name)
	if err != nil {
		panic(fmt.Sprintf("failed to compile %s: %v", name, err))
	}
	return sch
}

// validateEvalSchema validates a generic value against the eval JSON schema
// and returns a slice of human-readable error strings.
func validateEvalSchema(instance any) []string {
	return validateAgainstSchema(evalSchema, instance)
}

// validateTaskSchema validates a generic value against the task JSON schema
// and returns a slice of human-readable error strings.
func validateTaskSchema(instance any) []string {
	return validateAgainstSchema(taskSchema, instance)
}

func validateAgainstSchema(schema *jsonschema.Schema, instance any) []string {
	err := schema.Validate(instance)
	if err == nil {
		return nil
	}
	ve, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return []string{fmt.Sprintf("schema: %v", err)}
	}
	var errs []string
	collectSchemaErrors(ve, &errs)
	return errs
}

// collectSchemaErrors recursively collects leaf validation errors.
func collectSchemaErrors(ve *jsonschema.ValidationError, errs *[]string) {
	if len(ve.Causes) == 0 {
		var parts []string
		for _, s := range ve.InstanceLocation {
			parts = append(parts, s)
		}
		loc := "/"
		if len(parts) > 0 {
			loc = "/" + strings.Join(parts, "/")
		}
		*errs = append(*errs, fmt.Sprintf("schema: %s: %v", loc, ve.ErrorKind))
		return
	}
	for _, c := range ve.Causes {
		collectSchemaErrors(c, errs)
	}
}
