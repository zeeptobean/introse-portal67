package service

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func checkEventContracts(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	file, err := os.Open("../../contracts/events/v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	document, err := jsonschema.UnmarshalJSON(file)
	if err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	if err := compiler.AddResource("v1.schema.json", document); err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile("v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, owner := range []string{"identity", "academic", "registration", "affairs"} {
		query := "SELECT type, payload FROM " + pgx.Identifier{owner, "outbox"}.Sanitize()
		rows, err := pool.Query(ctx, query)
		if err != nil {
			t.Fatal(err)
		}
		type eventRow struct {
			Kind    string
			Payload []byte
		}
		batch, err := pgx.CollectRows(rows, pgx.RowToStructByPos[eventRow])
		if err != nil {
			t.Fatal(err)
		}
		if len(batch) == 0 {
			t.Errorf("no %s outbox events to validate", owner)
		}
		for _, event := range batch {
			value, err := jsonschema.UnmarshalJSON(bytes.NewReader(event.Payload))
			if err != nil {
				t.Fatalf("%s is invalid JSON: %v", event.Kind, err)
			}
			if err := schema.Validate(value); err != nil {
				t.Errorf("%s does not match the v1 contract: %v", event.Kind, err)
			}
		}
	}
}
