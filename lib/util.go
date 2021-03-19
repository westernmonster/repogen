package lib

import (
	"github.com/westernmonster/sqalx"
	"github.com/ztrue/tracerr"
)

type TableColumnDef struct {
	TableName              string `json:"table_name" db:"TABLE_NAME"`
	ColumnName             string `json:"column_name" db:"COLUMN_NAME"`
	ColumnComment          string `json:"column_comment" db:"COLUMN_COMMENT"`
	DataType               string `json:"data_type" db:"DATA_TYPE"`
	IsNullable             string `json:"is_nullable" db:"IS_NULLABLE"`
	CharacterMaximumLength *int   `json:"character_maximum_length" db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   *int   `json:"character_octet_length" db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       *int   `json:"numeric_precision" db:"NUMERIC_PRECISION"`
	NumericScale           *int   `json:"numeric_scale" db:"NUMERIC_SCALE"`
	ColumnType             string `json:"column_type" db:"COLUMN_TYPE"`
	OrdinalPosition        int    `json:"ordinal_position" db:"ORDINAL_POSITION"`
}

type TableName struct {
	TableName    string `json:"table_name" db:"TABLE_NAME"`
	TableComment string `json:"table_comment" db:"TABLE_COMMENT"`
}

func GetTableColumnDefs(node sqalx.Node, tableName, databaseName string) (items []*TableColumnDef, err error) {
	stmt, err := node.PrepareNamed(SqlGetTableColumns)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}
	items = make([]*TableColumnDef, 0)
	if e := stmt.Select(&items, map[string]interface{}{"tablename": tableName, "schema": databaseName}); e != nil {
		err = tracerr.Wrap(e)
		return
	}
	return
}

func GetTableNames(node sqalx.Node, databaseName string) (items []*TableName, err error) {
	stmt, err := node.PrepareNamed(SqlGetTableNames)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}
	items = make([]*TableName, 0)
	if e := stmt.Select(&items, map[string]interface{}{"database_name": databaseName}); e != nil {
		err = tracerr.Wrap(e)
		return
	}
	return
}
