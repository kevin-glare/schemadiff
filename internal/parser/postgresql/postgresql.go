package postgresql

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"schemadiff/internal/models"
)

// PostgreSQLParser реализует парсер для PostgreSQL
type PostgreSQLParser struct{}

// Parse парсит SQL-файл и возвращает схему
func (p *PostgreSQLParser) Parse(file string) (*models.Schema, error) {
	sql, err := readSQLFile(file)
	if err != nil {
		return nil, err
	}

	statements := splitSQLStatements(sql)
	schema := &models.Schema{
		Tables: make(map[string]models.Table),
	}

	for _, stmt := range statements {
		if strings.HasPrefix(strings.ToUpper(stmt), "CREATE TABLE") {
			table, err := parseCreateTable(stmt)
			if err != nil {
				log.Printf("Ошибка парсинга CREATE TABLE: %v", err)
				continue
			}
			schema.Tables[table.Name] = table
		}
	}

	return schema, nil
}

// readSQLFile читает SQL-файл и удаляет комментарии и пустые строки
func readSQLFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %v", err)
	}

	sql := string(data)
	sql = removeComments(sql)
	sql = removeEmptyLines(sql)

	return sql, nil
}

// removeComments удаляет комментарии из SQL
func removeComments(sql string) string {
	reSingleLine := regexp.MustCompile(`--.*`)
	sql = reSingleLine.ReplaceAllString(sql, "")

	reMultiLine := regexp.MustCompile(`/\*.*?\*/`)
	sql = reMultiLine.ReplaceAllString(sql, "")

	return sql
}

// removeEmptyLines удаляет пустые строки из SQL
func removeEmptyLines(sql string) string {
	lines := strings.Split(sql, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			cleanedLines = append(cleanedLines, trimmedLine)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// splitSQLStatements разделяет SQL-файл на отдельные запросы
func splitSQLStatements(sql string) []string {
	statements := strings.Split(sql, ";")
	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}
	return result
}

// parseCreateTable парсит CREATE TABLE
func parseCreateTable(stmt string) (models.Table, error) {
	table := models.Table{
		Name:        "",
		Columns:     make(map[string]models.Column),
		Indexes:     make(map[string]models.Index),
		Constraints: make(map[string]models.Constraint),
	}

	// Извлечение имени таблицы
	tableNameRegex := regexp.MustCompile(`(?i)CREATE TABLE (\w+)`)
	tableNameMatch := tableNameRegex.FindStringSubmatch(stmt)
	if len(tableNameMatch) < 2 {
		return table, fmt.Errorf("не удалось извлечь имя таблицы")
	}
	table.Name = tableNameMatch[1]

	// Извлечение тела таблицы
	bodyRegex := regexp.MustCompile(`(?i)CREATE TABLE \w+ \((.*)\)`)
	bodyMatch := bodyRegex.FindStringSubmatch(stmt)
	if len(bodyMatch) < 2 {
		return table, fmt.Errorf("не удалось извлечь тело таблицы")
	}
	body := bodyMatch[1]

	// Разделение строк тела таблицы
	lines := strings.Split(body, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Парсинг колонок
		if strings.HasPrefix(line, "PRIMARY KEY") || strings.HasPrefix(line, "FOREIGN KEY") || strings.HasPrefix(line, "UNIQUE") || strings.HasPrefix(line, "CHECK") {
			// Парсинг констрейнов
			constraint, err := parseConstraint(line)
			if err != nil {
				log.Printf("Ошибка парсинга констрейна: %v", err)
				continue
			}
			table.Constraints[constraint.Name] = constraint
		} else if strings.HasPrefix(line, "INDEX") || strings.HasPrefix(line, "UNIQUE INDEX") {
			// Парсинг индексов
			index, err := parseIndex(line)
			if err != nil {
				log.Printf("Ошибка парсинга индекса: %v", err)
				continue
			}
			table.Indexes[index.Name] = index
		} else {
			// Парсинг колонок
			column, err := parseColumn(line)
			if err != nil {
				log.Printf("Ошибка парсинга колонки: %v", err)
				continue
			}
			table.Columns[column.Name] = column
		}
	}

	return table, nil
}

// parseColumn парсит колонку
func parseColumn(line string) (models.Column, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return models.Column{}, fmt.Errorf("неверный формат колонки: %s", line)
	}

	column := models.Column{
		Name: parts[0],
		Type: parts[1],
	}

	// Обработка атрибутов колонки
	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "NOT":
			if i+1 < len(parts) && parts[i+1] == "NULL" {
				column.NotNull = true
				i++
			}
		case "DEFAULT":
			if i+1 < len(parts) {
				column.Default = strings.Trim(parts[i+1], "'")
				i++
			}
		case "COMMENT":
			if i+1 < len(parts) {
				column.Comment = strings.Trim(parts[i+1], "'")
				i++
			}
		}
	}

	return column, nil
}

// parseIndex парсит индекс
func parseIndex(line string) (models.Index, error) {
	re := regexp.MustCompile(`(?i)(UNIQUE )?INDEX (\w+) \((.+)\)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 4 {
		return models.Index{}, fmt.Errorf("неверный формат индекса: %s", line)
	}

	index := models.Index{
		Name:    matches[2],
		Columns: strings.Split(matches[3], ","),
		Unique:  matches[1] != "",
	}

	return index, nil
}

// parseConstraint парсит констрейн
func parseConstraint(line string) (models.Constraint, error) {
	re := regexp.MustCompile(`(?i)(PRIMARY KEY|FOREIGN KEY|UNIQUE|CHECK) (\w+)? \((.+)\)( REFERENCES (\w+)\((\w+)\))?`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 4 {
		return models.Constraint{}, fmt.Errorf("неверный формат констрейна: %s", line)
	}

	constraint := models.Constraint{
		Type:    matches[1],
		Columns: strings.Split(matches[3], ","),
	}

	if matches[1] == "FOREIGN KEY" && len(matches) >= 7 {
		constraint.References = matches[5] + "(" + matches[6] + ")"
	}

	return constraint, nil
}
