package dbml

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"schemadiff/internal/models"
)

// DBMLParser реализует парсер для DBML
type DBMLParser struct{}

// Parse парсит DBML-файл и возвращает схему
func (p *DBMLParser) Parse(file string) (*models.Schema, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	content := string(data)
	schema := &models.Schema{
		Tables: make(map[string]models.Table),
	}

	// Регулярное выражение для поиска таблиц
	tableRegex := regexp.MustCompile(`Table "(\w+)" \{([^}]+)\}`)
	tableMatches := tableRegex.FindAllStringSubmatch(content, -1)

	for _, match := range tableMatches {
		if len(match) < 3 {
			continue
		}

		tableName := match[1]
		tableBody := match[2]

		table := models.Table{
			Name:        tableName,
			Columns:     make(map[string]models.Column),
			Indexes:     make(map[string]models.Index),
			Constraints: make(map[string]models.Constraint),
		}

		// Парсинг колонок
		columnRegex := regexp.MustCompile(`"(\w+)"\s+([\w\(\)]+)\s*\[([^\]]+)\]`)
		columnMatches := columnRegex.FindAllStringSubmatch(tableBody, -1)

		for _, colMatch := range columnMatches {
			if len(colMatch) < 4 {
				continue
			}

			colName := colMatch[1]
			colType := colMatch[2]
			colAttributes := colMatch[3]

			column := models.Column{
				Name: colName,
				Type: colType,
			}

			// Парсинг атрибутов колонки
			attrs := strings.Split(colAttributes, ",")
			for _, attr := range attrs {
				attr = strings.TrimSpace(attr)
				switch {
				case attr == "pk":
					table.Constraints[colName+"_pk"] = models.Constraint{
						Name:    colName + "_pk",
						Type:    "PRIMARY KEY",
						Columns: []string{colName},
					}
				case attr == "not null":
					column.NotNull = true
				case strings.HasPrefix(attr, "default:"):
					column.Default = strings.Trim(strings.TrimPrefix(attr, "default:"), "`")
				case strings.HasPrefix(attr, "ref:"):
					refParts := strings.Split(strings.Trim(strings.TrimPrefix(attr, "ref:"), " >"), ".")
					if len(refParts) == 2 {
						table.Constraints[colName+"_fk"] = models.Constraint{
							Name:       colName + "_fk",
							Type:       "FOREIGN KEY",
							Columns:    []string{colName},
							References: refParts[0] + "(" + refParts[1] + ")",
						}
					}
				case strings.HasPrefix(attr, "note:"):
					column.Comment = strings.Trim(strings.TrimPrefix(attr, "note:"), "'")
				}
			}

			table.Columns[colName] = column
		}

		// Парсинг индексов
		indexRegex := regexp.MustCompile(`Indexes \{([^}]+)\}`)
		indexMatch := indexRegex.FindStringSubmatch(tableBody)
		if len(indexMatch) > 1 {
			indexBody := indexMatch[1]
			indexLines := strings.Split(indexBody, "\n")
			for _, line := range indexLines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// Парсинг индекса
				indexParts := strings.Split(line, "[")
				if len(indexParts) < 2 {
					continue
				}

				indexDef := indexParts[0]
				indexAttrs := indexParts[1]

				index := models.Index{
					Name:    "",
					Columns: []string{},
					Unique:  strings.Contains(indexAttrs, "unique"),
				}

				// Извлечение имени индекса
				nameRegex := regexp.MustCompile(`name:\s*"([^"]+)"`)
				nameMatch := nameRegex.FindStringSubmatch(indexAttrs)
				if len(nameMatch) > 1 {
					index.Name = nameMatch[1]
				}

				// Извлечение колонок
				columnsRegex := regexp.MustCompile(`\(([^)]+)\)`)
				columnsMatch := columnsRegex.FindStringSubmatch(indexDef)
				if len(columnsMatch) > 1 {
					index.Columns = strings.Split(columnsMatch[1], ",")
				} else {
					index.Columns = []string{strings.TrimSpace(indexDef)}
				}

				table.Indexes[index.Name] = index
			}
		}

		// Парсинг комментария к таблице
		noteRegex := regexp.MustCompile(`Note:\s*'([^']+)'`)
		noteMatch := noteRegex.FindStringSubmatch(tableBody)
		if len(noteMatch) > 1 {
			table.Comment = noteMatch[1]
		}

		schema.Tables[tableName] = table
	}

	return schema, nil
}
