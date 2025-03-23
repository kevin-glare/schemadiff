package models

// Column описывает структуру колонки
type Column struct {
	Name    string
	Type    string
	NotNull bool
	Default string
	Comment string
}

// Index описывает структуру индекса
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

// Constraint описывает структуру констрейна
type Constraint struct {
	Name       string
	Type       string // PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK
	Columns    []string
	References string // Для FOREIGN KEY
}

// Table описывает структуру таблицы
type Table struct {
	Name        string
	Columns     map[string]Column
	Indexes     map[string]Index
	Constraints map[string]Constraint
	Comment     string
}

// Schema описывает структуру схемы базы данных
type Schema struct {
	Tables map[string]Table
}
