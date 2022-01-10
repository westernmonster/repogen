package lib

import (
	"strings"

	"github.com/codemodus/kace"
	. "github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/spf13/viper"
	"github.com/ztrue/tracerr"
)

func ServiceListRespStruct(tableName string, columns []*TableColumnDef) (st *Statement, err error) {
	prefix := viper.GetString("development.prefix")
	tableName = strings.TrimPrefix(tableName, prefix)
	structName := kace.Pascal(tableName) + "ListResp"
	st = Type().Id(structName).Struct(
		Id("Items").Index().Op("*").Id(kace.Pascal(tableName)+"Item").Tag(map[string]string{"json": "items"}),
		Id("Total").Int().Tag(map[string]string{"json": "total"}),
	)
	return
}

func GenerateItemStruct(tableName string, columns []*TableColumnDef) (st *Statement, err error) {
	prefix := viper.GetString("development.prefix")
	tableName = strings.TrimPrefix(tableName, prefix)
	structName := kace.Pascal(tableName) + "Item"
	fields := make([]Code, 0)
	for _, v := range columns {
		columnName := kace.Pascal(v.ColumnName)
		ptr := Op("*")
		// comment := fmt.Sprintf("%s %s", columnName, v.ColumnComment)

		switch v.DataType {
		case "varchar":
			fallthrough
		case "text":
			fallthrough
		case "longtext":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).String().Tag(map[string]string{"json": v.ColumnName}))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).String().Tag(map[string]string{"json": v.ColumnName}))
		case "datetime":
			fallthrough
		case "timestamp":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Qual("time", "Time").Tag(map[string]string{"json": v.ColumnName}))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Qual("time", "Time").Tag(map[string]string{"json": v.ColumnName}))
		case "tinyint":
			fallthrough
		case "int":
			if v.IsNullable == "NO" {
				if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
					fields = append(fields, Id(columnName).Int32().Tag(map[string]string{"json": v.ColumnName}))
				} else {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName}))
				}
				break
			}
			if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
				fields = append(fields, Id(columnName).Add(ptr).Int32().Tag(map[string]string{"json": v.ColumnName}))
			} else {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName}))
			}
		case "bigint":
			if v.IsNullable == "NO" {
				if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName + ",string"}))
				} else {
					fields = append(fields, Id(columnName).Int64().Tag(map[string]string{"json": v.ColumnName}))
				}
				break
			}
			if v.ColumnName != "created_at" && v.ColumnName != "updated_at" {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName}))
			} else {
				fields = append(fields, Id(columnName).Add(ptr).Int64().Tag(map[string]string{"json": v.ColumnName}))
			}

		case "decimal":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Float64().Tag(map[string]string{"json": v.ColumnName}))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Float64().Tag(map[string]string{"json": v.ColumnName}))
		case "bit":
			if v.IsNullable == "NO" {
				fields = append(fields, Id(columnName).Bool().Tag(map[string]string{"json": v.ColumnName}))
				break
			}
			fields = append(fields, Id(columnName).Add(ptr).Bool().Tag(map[string]string{"json": v.ColumnName}))
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

func ServiceGetPaged(tableName string, columns []*TableColumnDef) (st *Statement, sqlCount, sqlSelect string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)
	pluralName := kace.Pascal(inflection.Plural(goName))

	fields := make([]Code, 0)
	for _, v := range columns {
		columnName := kace.Pascal(v.ColumnName)
		if v.DataType == "bit" {
			fields = append(fields, Id(columnName).Op(":").Id("bool").Call(Id("v.").Id(columnName)).Op(","))
		} else {
			fields = append(fields, Id(columnName).Op(":").Id("v.").Id(columnName).Op(","))
		}
	}

	st = Func().Params(Id("p").Op("*").Id("Service")).Id("Get"+pluralName+"ListPaged").Params(
		Id("c").Qual("context", "Context"),
		Id("cond").Map(String()).Interface(),
		List(Id("page"), Id("pageSize")).Int(),
	).Params(
		Id("resp").Op("*").Id("model."+pascalName+"ListResp"),
		Id("err").Error(),
	).Block(
		Id("offset").Op(":=").Call(Id("page").Op("-").Lit(1)).Op("*").Id("pageSize"),
		Line(),
		Var().Id("data").Index().Op("*").Id("model."+pascalName),
		Var().Id("total").Int(),
		If(
			List(Id("total"), Id("data"), Err()).Op("=").Id("p").Dot("d").Dot("Get"+pluralName+"Paged").Call(Id("c"), Id("p").Dot("d").Dot("DB").Call(), Id("cond"), Id("pageSize"), Id("offset")),
			Err().Op("!=").Nil().Block(Return()),
		),
		Line(),
		Id("resp").Op("=").Op("&").Id("model."+pascalName+"ListResp").Block(
			Id("Items").Op(":").Make(Index().Op("*").Id("model."+pascalName+"Item"), Len(Id("data"))).Op(","),
			Id("Total").Op(":").Id("total").Op(","),
		),
		Line(),
		For(List(Id("i"), Id("v")).Op(":=").Range().Id("data")).Block(
			Id("resp.Items").Index(Id("i")).Op("=").Op("&").Id("model."+pascalName+"Item").Block(
				fields...,
			),
		),
		Return(),
	)

	return
}

func ServiceGetAll(tableName string, columns []*TableColumnDef) (st *Statement, sqlSelect string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)
	pluralName := kace.Pascal(inflection.Plural(goName))

	fields := make([]Code, 0)
	for _, v := range columns {
		columnName := kace.Pascal(v.ColumnName)
		if v.DataType == "bit" {
			fields = append(fields, Id(columnName).Op(":").Id("bool").Call(Id("v.").Id(columnName)).Op(","))
		} else {
			fields = append(fields, Id(columnName).Op(":").Id("v.").Id(columnName).Op(","))
		}
	}

	st = Func().Params(Id("p").Op("*").Id("Service")).Id("GetAll"+pluralName).Params(
		Id("c").Qual("context", "Context"),
	).Params(
		Id("items").Index().Op("*").Id("model."+pascalName+"Item"),
		Id("err").Error(),
	).Block(
		Var().Id("data").Index().Op("*").Id("model."+pascalName),
		If(
			List(Id("data"), Err()).Op("=").Id("p").Dot("d").Dot("GetAll"+pluralName).Call(Id("c"), Id("p").Dot("d").Dot("DB").Call()),
			Err().Op("!=").Nil().Block(Return()),
		),
		Line(),
		Id("items").Op("=").Make(Index().Op("*").Id("model."+pascalName+"Item"), Len(Id("data"))),
		For(List(Id("i"), Id("v")).Op(":=").Range().Id("data")).Block(
			Id("items").Index(Id("i")).Op("=").Op("&").Id("model."+pascalName+"Item").Block(
				fields...,
			),
		),
		Return(),
	)

	return
}

func ServiceGetByID(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	st = Func().Params(Id("p").Op("*").Id("Service")).Id("Get"+pascalName+"ByID").Params(
		Id("c").Qual("context", "Context"),
		Id("id").Int64(),
	).Params(
		Id("item").Op("*").Id("model."+pascalName),
		Id("err").Error(),
	).Block(
		Return().Id("p.get"+pascalName+"ByID").Call(Id("c"), Id("id")),
	)

	return
}

func ServicegetByID(tableName string, columns []*TableColumnDef) (st *Statement, sql string) {
	prefix := viper.GetString("development.prefix")
	goName := strings.TrimPrefix(tableName, prefix)
	pascalName := kace.Pascal(goName)

	st = Func().Params(Id("p").Op("*").Id("Service")).Id("get"+pascalName+"ByID").Params(
		Id("c").Qual("context", "Context"),
		Id("id").Int64(),
	).Params(
		Id("item").Op("*").Id("model."+pascalName),
		Id("err").Error(),
	).Block(
		If(
			List(Id("item"), Err()).Op("=").Id("p.d.Get"+pascalName+"ByID").Call(
				Id("c"),
				Id("p").Dot("d").Dot("DB").Call(),
				Id("id"),
			),
			Err().Op("!=").Nil().Block(Return()),
		).Else().If(Id("item").Op("==").Nil()).Block(
			Id("err").Op("=").Id("ecode.").Id("NothingFound"),
			Return(),
		),
		Line(),
		Return(),
	)

	return
}
