package parser

import (
	"fmt"
	"schemadiff/internal/models"
	"schemadiff/internal/parser/dbml"
	"schemadiff/internal/parser/postgresql"
)

// Parser интерфейс для парсеров
type Parser interface {
	Parse(file string) (*models.Schema, error)
}

// NewParser создаёт парсер для указанного формата
func NewParser(format string) (Parser, error) {
	switch format {
	case "postgresql":
		return &postgresql.PostgreSQLParser{}, nil
	case "dbml":
		return &dbml.DBMLParser{}, nil
	default:
		return nil, fmt.Errorf("неподдерживаемый формат: %s", format)
	}
}
