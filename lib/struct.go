package lib

import (
	"fmt"
	"strings"

	"github.com/codemodus/kace"
	. "github.com/dave/jennifer/jen"
	"github.com/spf13/viper"
	"github.com/westernmonster/sqalx"
	"github.com/ztrue/tracerr"
)

func GenStruct(node sqalx.Node, tableName, databaseName string) (st *Statement, err error) {
	columns, err := GetTableColumnDefs(node, tableName, databaseName)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	// Generate Struct
	st, err = GenerateStruct(tableName, columns)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	return
}

func GenerateStruct(tableName string, columns []*TableColumnDef) (st *Statement, err error) {
	prefix := viper.GetString("development.prefix")
	tableName = strings.TrimPrefix(tableName, prefix)
	structName := kace.Pascal(tableName)
	fields := make([]Code, 0)
	for _, v := range columns {
		columnName := kace.Pascal(v.ColumnName)
		ptr := Op("*")
		comment := fmt.Sprintf("%s %s", columnName, v.ColumnComment)

		switch v.DataType {
		case "varchar":
			fallthrough
		case "text":
			fallthrough
		case "longtext":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).String().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).String().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
		case "datetime":
			fallthrough
		case "timestamp":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Qual("time", "Time").Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Qual("time", "Time").Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
		case "tinyint":
			fallthrough
		case "int":
			if v.IsNullable == "NO" {
				if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
					fields = append(fields, Id(columnName).Int32().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				} else {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				}
				break
			}
			if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
				fields = append(fields, Id(columnName).Add(ptr).Int32().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
			} else {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
			}
		case "bigint":
			if v.IsNullable == "NO" {
				if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName + ",string", "db": v.ColumnName}).Comment(comment))
				} else {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				}
				break
			}
			if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
			} else {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
			}

		case "decimal":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Float64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Float64().Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
		case "bit":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Qual("github.com/jmoiron/sqlx/types", "BitBool").Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Qual("github.com/jmoiron/sqlx/types", "BitBool").Tag(map[string]string{"json": v.ColumnName, "db": v.ColumnName}).Comment(comment))
		default:
			err = tracerr.Errorf("convert column %s, data type %s failed", v.ColumnName, v.DataType)
			return
		}

	}
	st = Type().Id(structName).Struct(
		fields...,
	)
	return
}
