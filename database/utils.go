package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FetchColumnNames(db *sqlx.DB, tableName string) ([]string, error) {
	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = $1")
	rows, err := db.Queryx(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

// snakeToCamel converts snake_case column names to CamelCase struct field names.
func snakeToCamel(s string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}

// ParseRows parses SQL rows into a slice of structs.
func ParseRows(rows *sql.Rows, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	sliceValue := destValue.Elem()
	elemType := sliceValue.Type().Elem()

	// Ensure elemType is a struct
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return errors.New("destination slice must contain struct elements")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Create a new instance of the struct (or a pointer to the struct)
		elem := reflect.New(elemType).Elem()

		for i, column := range columns {
			fieldName := snakeToCamel(column)
			field := elem.FieldByName(fieldName)

			if field.IsValid() && field.CanSet() {
				val := reflect.ValueOf(columnValues[i])

				// Handle null for slice fields (e.g., []string)
				if field.Type() == reflect.TypeOf([]string{}) {
					if val.IsNil() {
						field.Set(reflect.ValueOf([]string{}))
					} else {
						var slice []string
						if err := json.Unmarshal(columnValues[i].([]byte), &slice); err == nil {
							field.Set(reflect.ValueOf(slice))
						} else {
							log.Println("Failed to unmarshal JSON for field", fieldName, ":", err)
						}
					}
				} else if val.Type() == reflect.TypeOf(time.Time{}) {
					if timeVal, ok := val.Interface().(time.Time); ok {
						field.Set(reflect.ValueOf(timeVal))
					}
				} else {
					// Standard conversion for other types
					if val.Kind() == reflect.Ptr && !val.IsNil() {
						val = val.Elem()
					}
					if val.Type().ConvertibleTo(field.Type()) {
						field.Set(val.Convert(field.Type()))
					}
				}
			}
		}

		// Append the new struct pointer to the slice
		sliceValue.Set(reflect.Append(sliceValue, elem.Addr()))
	}

	return nil
}

// InsertStruct inserts a struct's fields into the specified table.
func InsertStruct(db *sqlx.DB, tableName string, data interface{}) error {
	// Prepare columns and values for the SQL query
	var columns []string
	var values []interface{}

	val := reflect.ValueOf(data).Elem() // Get the value pointed to by data
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		column := field.Tag.Get("json") // Get column name from `json` tag
		dbTag := field.Tag.Get("db")

		// Skip fields not mapped to a database column (e.g., with `db:"-"` tag)
		if dbTag == "-" || column == "" {
			continue
		}

		// Handle JSON marshaling for fields marked as JSON
		if dbTag == "json" {
			jsonValue, err := json.Marshal(val.Field(i).Interface())
			if err != nil {
				log.Println("Failed to marshal JSON field:", column, err)
				return err
			}
			values = append(values, jsonValue)
		} else {
			// Use the actual value with correct type
			values = append(values, val.Field(i).Interface())
		}

		columns = append(columns, column)
	}

	// Construct the SQL query with dynamic placeholders
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(columns, ", "), placeholders(len(values)))

	// Execute the query
	_, err := db.Exec(query, values...)
	if err != nil {
		log.Println("Failed to insert data into table:", err)
		return err
	}

	return nil
}

// UpdateStruct updates fields in the specified table for a given struct based on a condition.
func UpdateStruct(db *sqlx.DB, tableName string, data interface{}, conditionField string, conditionValue interface{}) error {
	var columns []string
	var values []interface{}

	val := reflect.ValueOf(data).Elem() // Get the value pointed to by data
	typ := val.Type()

	placeholderIdx := 1 // Start placeholder index at 1

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		column := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")

		// Skip fields not mapped to a database column (e.g., with `db:"-"` tag)
		if dbTag == "-" || column == "" {
			continue
		}

		// Skip the condition field to avoid updating it
		if column == conditionField {
			continue
		}

		// Handle JSON marshalling for fields marked as JSON
		if dbTag == "json" {
			jsonValue, err := json.Marshal(val.Field(i).Interface())
			if err != nil {
				log.Println("Failed to marshal JSON field:", column, err)
				return err
			}
			values = append(values, jsonValue)
		} else {
			// Convert values to the correct type based on the struct's field type
			values = append(values, val.Field(i).Interface())
		}

		// Add the column update statement with the placeholder
		columns = append(columns, fmt.Sprintf("%s = $%d", column, placeholderIdx))
		placeholderIdx++
	}

	// Add the condition field and its value as the last placeholder
	values = append(values, conditionValue)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $%d`, tableName, strings.Join(columns, ", "), conditionField, placeholderIdx)

	// Log the final query for debugging
	fmt.Println("Generated SQL query:", query)

	// Execute the query
	_, err := db.Exec(query, values...)
	if err != nil {
		log.Println("Failed to update record in table:", err)
		return err
	}

	return nil
}

// placeholders generates a string of placeholders for SQL based on the number of fields.
func placeholders(n int) string {
	ph := make([]string, n)
	for i := range ph {
		ph[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(ph, ", ")
}
