package jsonschema_test

import (
	"testing"

	jsonschema "github.com/bmeg/jsonschema/v6"
)

func TestValidateFastMatchesValidate(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	const schemaURL = "https://example.org/fast-validation.json"
	if err := compiler.AddResource(schemaURL, map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type":    "object",
		"required": []any{"id"},
		"properties": map[string]any{
			"id":   map[string]any{"type": "string"},
			"name": map[string]any{"type": "string"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatal(err)
	}

	valid := map[string]any{"id": "row-1", "name": "valid"}
	if err := sch.Validate(valid); err != nil {
		t.Fatalf("rich validation rejected valid input: %v", err)
	}
	if err := sch.ValidateFast(valid); err != nil {
		t.Fatalf("fast validation rejected valid input: %v", err)
	}

	invalid := map[string]any{"id": 42}
	if err := sch.Validate(invalid); err == nil {
		t.Fatal("rich validation accepted invalid input")
	}
	if err := sch.ValidateFast(invalid); err == nil {
		t.Fatal("fast validation accepted invalid input")
	}
}

func BenchmarkValidate(b *testing.B) {
	benchmarkValidation(b, false)
}

func BenchmarkValidateFast(b *testing.B) {
	benchmarkValidation(b, true)
}

func benchmarkValidation(b *testing.B, fast bool) {
	b.Helper()
	compiler := jsonschema.NewCompiler()
	const schemaURL = "https://example.org/benchmark-validation.json"
	if err := compiler.AddResource(schemaURL, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string"},
		},
	}); err != nil {
		b.Fatal(err)
	}
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		b.Fatal(err)
	}
	value := map[string]any{"id": "row-1"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if fast {
			err = sch.ValidateFast(value)
		} else {
			err = sch.Validate(value)
		}
		if err != nil {
			b.Fatal(err)
		}
	}
}
