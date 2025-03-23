package main

import (
	"fmt"
	"log"
	"os"

	"schemadiff/internal/comparator"
	"schemadiff/internal/parser"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Использование: schemadiff <формат1> <файл1> <формат2> <файл2>")
		os.Exit(1)
	}

	format1 := os.Args[1]
	file1 := os.Args[2]
	format2 := os.Args[3]
	file2 := os.Args[4]

	// Парсинг первой схемы
	parser1, err := parser.NewParser(format1)
	if err != nil {
		log.Fatalf("Ошибка создания парсера для %s: %v", format1, err)
	}
	schema1, err := parser1.Parse(file1)
	if err != nil {
		log.Fatalf("Ошибка парсинга %s: %v", file1, err)
	}

	// Парсинг второй схемы
	parser2, err := parser.NewParser(format2)
	if err != nil {
		log.Fatalf("Ошибка создания парсера для %s: %v", format2, err)
	}
	schema2, err := parser2.Parse(file2)
	if err != nil {
		log.Fatalf("Ошибка парсинга %s: %v", file2, err)
	}

	// Сравнение схем
	areEqual, differences := comparator.CompareSchemas(schema1, schema2)
	if areEqual {
		fmt.Println("Схемы идентичны.")
	} else {
		fmt.Println("Схемы отличаются:")
		for key, diff := range differences {
			fmt.Printf("%s: %s\n", key, diff)
		}
	}
}
