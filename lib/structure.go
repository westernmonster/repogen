package lib

import (
	"fmt"
	"strings"

	"github.com/codemodus/kace"
	. "github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/spf13/viper"
)

func GenInterface(tableName string) (st *Statement, err error) {
	prefix := viper.GetString("development.prefix")
	tableName = strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(tableName)

	st = Type().Id("d").Interface(
		Id("Get"+pascalName+"ByCond").Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("cond").Map(String()).Interface(),
		).Params(
			Id("items").Index().Op("*").Id("model."+pascalName),
			Id("err").Error(),
		),
		Id("Get"+pascalName).Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
		).Params(
			Id("items").Index().Op("*").Id("model."+pascalName),
			Id("err").Error(),
		),
		Id("Get"+pascalName+"ByID").Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("id").Int64(),
		).Params(
			Id("item").Op("*").Id("model."+pascalName),
			Id("err").Error(),
		),
		Id("Get"+pascalName+"ByCond").Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("cond").Map(String()).Interface(),
		).Params(
			Id("item").Op("*").Id("model."+pascalName),
			Id("err").Error(),
		),
		Id("Add"+pascalName).Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("item").Op("*").Id("model."+pascalName),
		).Params(
			Id("err").Error(),
		),
		Id("Update"+pascalName).Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("item").Op("*").Id("model."+pascalName),
		).Params(
			Id("err").Error(),
		),
		Id("Del"+pascalName).Params(
			Id("c").Qual("context", "Context"),
			Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
			Id("id").Int64(),
		).Params(
			Id("err").Error(),
		),
	)

	return st, nil
}

func GenerateGetAll(tableName string, columns []*TableColumnDef) (st *Statement, sqlSelect string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	cols := make([]string, 0)
	for _, v := range columns {
		cols = append(cols, "a."+v.ColumnName)
	}
	sqlSelect = `SELECT %s FROM %s WHERE 1=1 ORDER BY a.id DESC `
	sqlSelect = fmt.Sprintf(sqlSelect, strings.Join(cols, ","), tableName+" a")

	st = Func().Params(Id("p").Op("*").Id("Dao")).Id("Get"+kace.Pascal(inflection.Plural(goName))).Params(
		Id("c").Qual("context", "Context"),
		Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
	).Params(
		Id("items").Index().Op("*").Id("model."+pascalName),
		Id("err").Error(),
	).Block(
		Id("items").Op("=").Make(Index().Op("*").Id("model."+pascalName), Lit(0)),
		Id("sqlSelect").Op(":=").Lit(sqlSelect),
		Line(),
		If(
			Err().Op("=").Id("node").Dot("SelectContext").Call(Id("c"), Op("&").Id("items"), Id("sqlSelect")), Err().Op("!=").Nil()).Block(
			Id("log").Dot("For").Call(Id("c")).Dot("Errorf").Call(Lit("dao.Get"+pascalName+" err(%+v)"), Id("err")),
			Return(),
		),
		Return(),
	)

	return
}

func GenerateGetByID(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	cols := make([]string, 0)
	for _, v := range columns {
		cols = append(cols, "a."+v.ColumnName)
	}
	sql = `SELECT %s FROM %s WHERE a.id=?`
	sql = fmt.Sprintf(sql, strings.Join(cols, ","), tableName+" a")

	st = Func().Params(Id("p").Op("*").Id("Dao")).Id("Get"+pascalName+"ByID").Params(
		Id("c").Qual("context", "Context"),
		Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
		Id("id").Int64(),
	).Params(
		Id("item").Op("*").Id("model."+pascalName),
		Id("err").Error(),
	).Block(
		Id("item").Op("=").New(Id("model."+pascalName)),
		Id("sqlSelect").Op(":=").Lit(sql),
		Line(),
		If(
			Id("err").Op("=").Id("node").Dot("GetContext").
				Call(Id("c"), Id("item"), Id("sqlSelect"), Id("id")), Id("err").Op("!=").Nil()).Block(
			If(Id("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				Id("item").Op("=").Nil(),
				Id("err").Op("=").Nil(),
				Return(),
			),
			Id("log").Dot("For").Call(Id("c")).Dot("Errorf").Call(Lit("dao.Get"+pascalName+"ByID err(%+v), id(%+v)"), Id("err"), Id("id")),
		),
		Line(),
		Return(),
	)

	return
}

func GenerateInsert(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	sql = `INSERT INTO %s( %s) VALUES ( %s)`

	columnStr := ""
	valueStr := ""

	codes := make([]Code, 0)
	codes = append(codes, Id("c"))
	codes = append(codes, Id("sqlInsert"))
	lenCols := len(columns)
	for idx, v := range columns {
		codes = append(codes, Id("item").Dot(kace.Pascal(v.ColumnName)))
		columnStr += fmt.Sprintf("%s,", v.ColumnName)
		if idx == lenCols-1 {
			columnStr = strings.TrimRight(columnStr, ",")
		}

		valueStr += "?,"
		if idx == lenCols-1 {
			valueStr = strings.TrimRight(valueStr, ",")
		}
	}

	sql = fmt.Sprintf(sql, tableName, columnStr, valueStr)

	st = Func().Params(Id("p").Op("*").Id("Dao")).Id("Add"+pascalName).Params(
		Id("c").Qual("context", "Context"),
		Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
		Id("item").Op("*").Id("model."+pascalName),
	).Params(
		Id("err").Error(),
	).Block(
		Id("sqlInsert").Op(":=").Lit(sql),
		Line(),
		If(List(Id("_"), Err()).Op("=").Id("node").Dot("ExecContext").Call(codes...), Err().Op("!=").Nil()).Block(
			Id("log").Dot("For").Call(Id("c")).Dot("Errorf").Call(Lit("dao.Add"+pascalName+" err(%+v), item(%+v)"), Id("err"), Id("item")),
			Return(),
		),
		Line(),
		Return(),
	)

	return
}

func GenerateUpdate(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	sql = `UPDATE %s SET %s WHERE id=?`
	columnStr := ""
	lenCols := len(columns)

	codes := make([]Code, 0)
	codes = append(codes, Id("c"))
	codes = append(codes, Id("sqlUpdate"))
	for idx, v := range columns {
		if v.ColumnName == "id" || v.ColumnName == "created_at" || v.ColumnName == "deleted" {
			continue
		}

		codes = append(codes, Id("item").Dot(kace.Pascal(v.ColumnName)))
		columnStr += fmt.Sprintf("%s=?,", v.ColumnName)
		if idx == lenCols-1 {
			columnStr = strings.TrimRight(columnStr, ",")
		}

	}
	codes = append(codes, Id("item").Dot("ID"))

	sql = fmt.Sprintf(sql, tableName, columnStr)

	st = Func().Params(Id("p").Op("*").Id("Dao")).Id("Update"+pascalName).Params(
		Id("c").Qual("context", "Context"),
		Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
		Id("item").Op("*").Id("model."+pascalName),
	).Params(
		Id("err").Error(),
	).Block(
		Id("sqlUpdate").Op(":=").Lit(sql),
		Line(),

		List(Id("_"), Err()).Op("=").Id("node").Dot("ExecContext").Call(codes...),
		If(Err().Op("!=").Nil()).Block(
			Id("log").Dot("For").Call(Id("c")).Dot("Errorf").Call(Lit("dao.Update"+pascalName+" err(%+v), item(%+v)"), Id("err"), Id("item")),
			Return(),
		),
		Line(),
		Return(),
	)

	return
}

func GenerateDelete(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	sql = `DELETE FROM %s WHERE id=? `
	sql = fmt.Sprintf(sql, tableName)

	st = Func().Params(Id("p").Op("*").Id("Dao")).Id("Del"+pascalName).Params(
		Id("c").Qual("context", "Context"),
		Id("node").Qual("github.com/westernmonster/sqalx", "Node"),
		Id("id").Int64(),
	).Params(
		Id("err").Error(),
	).Block(
		Id("sqlDelete").Op(":=").Lit(sql),
		Line(),

		If(List(Id("_"), Err()).Op("=").Id("node").Dot("ExecContext").Call(Id("c"), Id("sqlDelete"), Id("id")),
			Err().Op("!=").Nil()).Block(
			Id("log").Dot("For").Call(Id("c")).Dot("Errorf").Call(Lit("dao.Del"+pascalName+" err(%+v), item(%+v)"), Id("err"), Id("id")),
			Return(),
		),
		Line(),
		Return(),
	)

	return
}
