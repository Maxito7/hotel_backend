package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Conexión a la base de datos
	connStr := "host=hotel-db-reservas.c5uw44cqcsx6.us-east-1.rds.amazonaws.com port=5432 user=postgres password=vamosB0ys!7 dbname=postgres sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== TABLAS EN LA BASE DE DATOS ===\n")

	// Listar todas las tablas
	tablesQuery := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name;
	`
	rows, err := db.Query(tablesQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, tableName)
		fmt.Printf("✓ %s\n", tableName)
	}

	fmt.Println("\n=== ESTRUCTURA DE CADA TABLA ===\n")

	// Para cada tabla, mostrar sus columnas
	for _, table := range tables {
		fmt.Printf("\n📋 Tabla: %s\n", table)
		fmt.Println("-----------------------------------")

		columnsQuery := `
			SELECT 
				column_name, 
				data_type, 
				is_nullable,
				column_default
			FROM information_schema.columns 
			WHERE table_name = $1 
			ORDER BY ordinal_position;
		`
		colRows, err := db.Query(columnsQuery, table)
		if err != nil {
			log.Printf("Error getting columns for %s: %v", table, err)
			continue
		}

		for colRows.Next() {
			var colName, dataType, isNullable string
			var colDefault sql.NullString
			if err := colRows.Scan(&colName, &dataType, &isNullable, &colDefault); err != nil {
				log.Fatal(err)
			}

			nullable := ""
			if isNullable == "NO" {
				nullable = " NOT NULL"
			}

			defaultVal := ""
			if colDefault.Valid {
				defaultVal = fmt.Sprintf(" DEFAULT %s", colDefault.String)
			}

			fmt.Printf("  • %-25s %-20s%s%s\n", colName, dataType, nullable, defaultVal)
		}
		colRows.Close()

		// Mostrar claves primarias
		pkQuery := `
			SELECT a.attname
			FROM pg_index i
			JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
			WHERE i.indrelid = $1::regclass AND i.indisprimary;
		`
		pkRows, err := db.Query(pkQuery, table)
		if err == nil {
			var pks []string
			for pkRows.Next() {
				var pkName string
				if err := pkRows.Scan(&pkName); err == nil {
					pks = append(pks, pkName)
				}
			}
			if len(pks) > 0 {
				fmt.Printf("  🔑 PRIMARY KEY: %v\n", pks)
			}
			pkRows.Close()
		}

		// Mostrar claves foráneas
		fkQuery := `
			SELECT
				kcu.column_name,
				ccu.table_name AS foreign_table_name,
				ccu.column_name AS foreign_column_name
			FROM information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
				ON tc.constraint_name = kcu.constraint_name
				AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage AS ccu
				ON ccu.constraint_name = tc.constraint_name
				AND ccu.table_schema = tc.table_schema
			WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1;
		`
		fkRows, err := db.Query(fkQuery, table)
		if err == nil {
			hasFk := false
			for fkRows.Next() {
				if !hasFk {
					fmt.Println("  🔗 FOREIGN KEYS:")
					hasFk = true
				}
				var colName, foreignTable, foreignCol string
				if err := fkRows.Scan(&colName, &foreignTable, &foreignCol); err == nil {
					fmt.Printf("     %s -> %s(%s)\n", colName, foreignTable, foreignCol)
				}
			}
			fkRows.Close()
		}
	}

	fmt.Println("\n=== VERIFICACIÓN COMPLETADA ===")
}
