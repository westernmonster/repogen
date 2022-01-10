package lib

import (
	"strings"

	"github.com/codemodus/kace"
	. "github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/spf13/viper"
)

func HTTPGetPaged(tableName string, columns []*TableColumnDef) (st *Statement, sqlCount, sqlSelect string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	// pascalName := kace.Pascal(goName)
	pluralName := kace.Pascal(inflection.Plural(goName))
	st = Func().Id("get"+pluralName+"ListPaged").Params(
		Id("c").Op("*").Qual("vin", "Context"),
	).Block(
		Id("filters").Op(":=").Id("c.Query").Call(Lit("filters")),
		Line(),
		Id("cond").Op(":=").Make(Map(String()).Interface()),
		If(Qual("strings", "TrimSpace").Call(Id("filters")).Op("!=").Lit("")).Block(
			If(Err().Op(":=").Qual("jsoniter", "Unmarshal").Call(Index().Byte().Call(Id("filters")), Op("&").Id("cond")),
				Err().Op("!=").Nil()).Block(
				Id("c.JSON").Call(Nil(), Id("ecode.RequestErr")),
				Return(),
			),
		),
		Line(),
		Id("page").Op(":=").Id("c.QueryIntDefault").Call(Lit("page"), Lit(1)),
		Id("pageSize").Op(":=").Id("c.QueryIntDefault").Call(Lit("page_size"), Lit(10)),
		Line(),
		Id("c.JSON").Call(Id("srv.Get"+pluralName+"ListPaged").Call(Id("c"), Id("cond"), Id("page"), Id("pageSize"))),
	)

	return
}

func HTTPGetAll(tableName string, columns []*TableColumnDef) (st *Statement, sqlSelect string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	// pascalName := kace.Pascal(goName)
	pluralName := kace.Pascal(inflection.Plural(goName))

	st = Func().Id("getAll" + pluralName).Params(
		Id("c").Op("*").Qual("vin", "Context"),
	).Block(
		Id("c.JSON").Call(Id("srv.GetAll" + pluralName).Call(Id("c"))),
	)

	return
}

func HTTPGetByID(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	st = Func().Id("get"+pascalName+"ByID").Params(
		Id("c").Op("*").Qual("vin", "Context"),
	).Block(
		List(Id("id"), Err()).Op(":=").Id("c.QueryInt64").Call(Lit("id")),
		If(Id("err").Op("!=").Nil()).Block(
			Id("c.JSON").Call(Nil(), Id("ecode.RequestErr")),
			Return(),
		),
		Line(),
		Id("c.JSON").Call(Id("srv.Get"+pascalName+"ByID").Call(Id("c"), Id("id"))),
	)

	return
}
