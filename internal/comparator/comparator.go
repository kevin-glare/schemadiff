package comparator

import (
	"reflect"

	"schemadiff/internal/models"
)

// CompareSchemas сравнивает две схемы и возвращает различия
func CompareSchemas(schema1, schema2 *models.Schema) (bool, map[string]string) {
	differences := make(map[string]string)

	// Сравнение таблиц
	for tableName, table1 := range schema1.Tables {
		table2, exists := schema2.Tables[tableName]
		if !exists {
			differences[tableName] = "Таблица отсутствует во второй схеме"
			continue
		}

		// Сравнение колонок
		for colName, col1 := range table1.Columns {
			col2, exists := table2.Columns[colName]
			if !exists {
				differences[tableName+"."+colName] = "Колонка отсутствует во второй схеме"
				continue
			}
			if !reflect.DeepEqual(col1, col2) {
				differences[tableName+"."+colName] = "Колонки отличаются"
			}
		}

		// Сравнение индексов
		for idxName, idx1 := range table1.Indexes {
			idx2, exists := table2.Indexes[idxName]
			if !exists {
				differences[tableName+"."+idxName] = "Индекс отсутствует во второй схеме"
				continue
			}
			if !reflect.DeepEqual(idx1, idx2) {
				differences[tableName+"."+idxName] = "Индексы отличаются"
			}
		}

		// Сравнение констрейнов
		for constrName, constr1 := range table1.Constraints {
			constr2, exists := table2.Constraints[constrName]
			if !exists {
				differences[tableName+"."+constrName] = "Констрейн отсутствует во второй схеме"
				continue
			}
			if !reflect.DeepEqual(constr1, constr2) {
				differences[tableName+"."+constrName] = "Констрейны отличаются"
			}
		}
	}

	// Проверка таблиц, которые есть во второй схеме, но отсутствуют в первой
	for tableName := range schema2.Tables {
		if _, exists := schema1.Tables[tableName]; !exists {
			differences[tableName] = "Таблица отсутствует в первой схеме"
		}
	}

	return len(differences) == 0, differences
}
