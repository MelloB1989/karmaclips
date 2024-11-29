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
	"unicode"

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

func ParseRows(rows *sql.Rows, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	sliceValue := destValue.Elem()
	elemType := sliceValue.Type().Elem()

	// Ensure elemType is a struct
	isPtr := false
	if elemType.Kind() == reflect.Ptr {
		isPtr = true
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return errors.New("destination slice must contain struct elements")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Build a mapping from column names to struct fields
	columnToFieldMap := make(map[string]reflect.StructField)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnToFieldMap[dbTag] = field
		} else {
			columnToFieldMap[camelToSnake(field.Name)] = field
		}
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
			fieldInfo, ok := columnToFieldMap[column]
			if !ok {
				// Try snakeToCamel conversion
				fieldInfo, ok = columnToFieldMap[snakeToCamel(column)]
				if !ok {
					continue // Skip columns without corresponding struct fields
				}
			}

			field := elem.FieldByName(fieldInfo.Name)
			if field.IsValid() && field.CanSet() {
				val := reflect.ValueOf(columnValues[i])

				// Handle fields with 'db' tag
				if fieldInfo.Tag.Get("db") != "" {
					// Obtain the data as []byte
					var data []byte
					switch v := val.Interface().(type) {
					case []byte:
						data = v
					case string:
						data = []byte(v)
					default:
						log.Println("Unsupported type for field with 'db' tag:", fieldInfo.Name)
						continue
					}

					// Unmarshal the JSON into the field
					if err := json.Unmarshal(data, field.Addr().Interface()); err != nil {
						log.Println("Failed to unmarshal JSON for field", fieldInfo.Name, ":", err)
						continue
					}
				} else {
					// Standard processing for other fields
					if val.Kind() == reflect.Ptr && !val.IsNil() {
						val = val.Elem()
					}
					if val.Kind() == reflect.Interface && !val.IsNil() {
						val = val.Elem()
					}
					if val.Type().ConvertibleTo(field.Type()) {
						field.Set(val.Convert(field.Type()))
					} else if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8 && field.Kind() == reflect.String {
						// Convert []byte to string
						field.SetString(string(val.Interface().([]byte)))
					} else {
						log.Println("Cannot set field", fieldInfo.Name, "with value of type", val.Type())
					}
				}
			}
		}

		// Append the new struct (or pointer) to the slice
		if isPtr {
			sliceValue.Set(reflect.Append(sliceValue, elem.Addr()))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, elem))
		}
	}

	return nil
}

func snakeToCamel(s string) string {
	caser := cases.Title(language.English)
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, "")
}

func camelToSnake(s string) string {
	var snake string
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 {
				snake += "_"
			}
			snake += string(unicode.ToLower(c))
		} else {
			snake += string(c)
		}
	}
	return snake
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
