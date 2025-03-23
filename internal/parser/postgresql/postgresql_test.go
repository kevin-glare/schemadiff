package postgresql

import (
	"testing"
)

func TestPostgreSQLParser_Parse(t *testing.T) {
	parser := &PostgreSQLParser{}
	schema, err := parser.Parse("testdata/schema.sql")
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}

	if len(schema.Tables) == 0 {
		t.Error("Схема не должна быть пустой")
	}

	// Проверка наличия таблицы "users"
	usersTable, exists := schema.Tables["users"]
	if !exists {
		t.Error("Таблица 'users' не найдена")
	}

	// Проверка колонок таблицы "users"
	if len(usersTable.Columns) != 3 {
		t.Errorf("Ожидалось 3 колонки, получено %d", len(usersTable.Columns))
	}
}
