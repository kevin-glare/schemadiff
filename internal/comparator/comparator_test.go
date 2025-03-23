package comparator

import (
	"testing"

	"schemadiff/internal/models"
)

func TestCompareSchemas(t *testing.T) {
	schema1 := &models.Schema{
		Tables: map[string]models.Table{
			"users": {
				Name: "users",
				Columns: map[string]models.Column{
					"id":       {Name: "id", Type: "SERIAL"},
					"username": {Name: "username", Type: "VARCHAR(50)", NotNull: true},
				},
			},
		},
	}

	schema2 := &models.Schema{
		Tables: map[string]models.Table{
			"users": {
				Name: "users",
				Columns: map[string]models.Column{
					"id":       {Name: "id", Type: "SERIAL"},
					"username": {Name: "username", Type: "VARCHAR(50)", NotNull: true},
				},
			},
		},
	}

	areEqual, differences := CompareSchemas(schema1, schema2)
	if !areEqual {
		t.Errorf("Схемы должны быть идентичны, различия: %v", differences)
	}
}
