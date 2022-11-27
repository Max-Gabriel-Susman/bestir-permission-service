package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/gocraft/dbr/v2"
)

type (
	Column  = string
	Columns []Column
)

type Table struct {
	Name    string
	Columns Columns
}

// NewTable takes an obj and looks for all fields that have the
// `db` tag set. If the field has the `table` tag set it will be ignored.
func NewTable(name string, obj interface{}) Table {
	columns := Columns{}
	forEachDBColumn(obj, func(column, _ string, hasTag bool) {
		if !hasTag {
			columns = append(columns, column)
		}
	})
	return Table{
		Name:    name,
		Columns: columns,
	}
}

// PrefixColumnsWithName turns all the column names into
// $table_name.$column_name as $column_name.
func (t Table) PrefixColumnsWithName() Table {
	columns := make(Columns, 0, len(t.Columns))
	for _, c := range t.Columns {
		columns = append(columns, prefixColumnWithTable(t.Name, c))
	}
	return Table{
		Name:    t.Name,
		Columns: columns,
	}
}

// UpdateFrom looks over the fields on the obj and genereates
// a map of the fields db tag and their values. If the field
// is a nil pointer then it is ignored.
func (t Table) UpdateFrom(obj interface{}, ignoreColumns ...string) (map[string]interface{}, error) {
	updateMap := map[string]interface{}{}
	value := reflect.ValueOf(obj)
	typ := value.Type()

	if typ.Kind() != reflect.Struct {
		return updateMap, errors.New("expected a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		columnName := getFieldDBName(typ.Field(i))
		if columnName == "" {
			continue
		}

		ignore := false
		for _, name := range ignoreColumns {
			if name == columnName {
				ignore = true
			}
		}
		if ignore {
			continue
		}

		fieldValue := value.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			v, ok := fieldValue.Interface().(driver.Valuer)
			if ok {
				value, err := v.Value()
				if err != nil {
					return updateMap, err
				}
				if value != nil {
					updateMap[columnName] = value
				}
				continue
			}
		}

		if fieldValue.Kind() == reflect.Ptr {
			// only grab non nil pointer values
			if !fieldValue.IsNil() {
				value := fieldValue.Elem().Interface()
				if value != nil {
					updateMap[columnName] = value
				}
			}
			continue
		}

		// take the value directly
		if value := fieldValue.Interface(); value != nil {
			updateMap[columnName] = value
		}
	}
	return updateMap, nil
}

func prefixColumnWithTable(table, column string) string {
	return fmt.Sprintf("%s.%s as %s", table, column, column)
}

type Query struct {
	Columns Columns
}

// NewQuery takes an obj and looks for all fields that have the
// `db` and `table` tag set. If the field only has the `db` tag
// it will be ignored.
func NewQuery(obj interface{}) Query {
	columns := Columns{}
	forEachDBColumn(obj, func(column, tableTag string, hasTag bool) {
		if hasTag {
			columns = append(columns, getTableTagColumnName(column, tableTag))
		}
	})
	return Query{Columns: columns}
}

// NewQueryWithDefaultTable takes an obj and name, and looks for all fields that have
// the `db` and/or `table` tags set. If the field only has the `db` tag
// set and not `table`, the supplied name will be used instead.
func NewQueryWithDefaultTable(obj interface{}, table string) Query {
	columns := Columns{}
	forEachDBColumn(obj, func(column, tableTag string, hasTag bool) {
		if hasTag {
			columns = append(columns, getTableTagColumnName(column, tableTag))
			return
		}

		tableTag = table
		columns = append(columns, prefixColumnWithTable(tableTag, column))
	})
	return Query{Columns: columns}
}

func getTableTagColumnName(dbTag, tableTag string) string {
	// if the table tag specifies a table and column, render the column as
	// "table.column as dbTag"
	if vals := strings.Split(tableTag, "."); len(vals) == 2 {
		return fmt.Sprintf("%s as %s", tableTag, dbTag)
	}
	return prefixColumnWithTable(tableTag, dbTag)
}

// forEachDBColumn iterates over all of the objects fields
// and calls fn on each field that has `db` tag set.
func forEachDBColumn(obj interface{}, fn func(column, tableTag string, hasTag bool)) {
	typ := reflect.TypeOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		columnName := getFieldDBName(field)
		if columnName == "" {
			continue
		}

		table, hasTableTag := field.Tag.Lookup("table")
		fn(columnName, table, hasTableTag)
	}
}

func getFieldDBName(field reflect.StructField) string {
	columnName := field.Tag.Get("db")
	switch columnName {
	case "-":
		return ""
	case "":
		// ignore non public fields without a tag
		if unicode.IsLower(rune(field.Name[0])) {
			return ""
		}
		columnName = dbr.NameMapping(field.Name)
	}
	return columnName
}
