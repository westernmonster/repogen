package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/InVisionApp/tabular"
	"github.com/codemodus/kace"
	. "github.com/dave/jennifer/jen"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/inflection"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/westernmonster/sqalx"
	"github.com/ztrue/tracerr"

	"github.com/westernmonster/repogen/lib"
)

var tab tabular.Table

func init() {
	tab = tabular.New()
	tab.Col("tn", "Table Name", 20)
	tab.Col("cn", "Column Name", 20)
	tab.Col("ct", "Column Type", 14)
	tab.Col("null", "Is Nullable", 12)
	tab.Col("op", "OrdinalPosition", 15)
	tab.Col("cm", "Comment", 20)

}

var rootCmd = &cobra.Command{
	Use:   "repogen",
	Short: "Fly Wiki repository generator",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// 1. With empty args: List all table names in database
// 2. With args: list specified table's columns
var listTableCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all tables in database",
	Long:  `Print all tables in database`,
	Run: func(cmd *cobra.Command, args []string) {
		_, node, err := initDatabase()
		if err != nil {
			log.Fatal(err)
		}
		database := viper.GetString("development.database")
		if len(args) == 0 {

			tables, err := lib.GetTableNames(node, database)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range tables {
				fmt.Printf("%s \t \t %s\n", v.TableName, v.TableComment)
			}

			return
		}

		for _, v := range args {
			columns, err := lib.GetTableColumnDefs(node, v, database)
			if err != nil {
				log.Fatal(err)
			}

			format := tab.Print("*")
			for _, v := range columns {
				fmt.Printf(format, v.TableName, v.ColumnName, v.ColumnType, v.IsNullable, v.OrdinalPosition, v.ColumnComment)
			}

			fmt.Println()
		}
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Code",
	Long:  "Generate Code",
	Args:  cobra.MinimumNArgs(1),
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Generate Repository",
	Long:  "Generate Repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, node, err := initDatabase()
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range args {
			err := generateRepo(node, v)
			if err != nil {
				log.Fatal(err)
			}

		}
	},
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Generate Struct",
	Long:  "Generate Struct",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database := viper.GetString("development.database")

		_, node, err := initDatabase()
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range args {
			code, err := lib.GenStruct(node, v, database)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("\n%#v\n", code)
		}
	},
}

func main() {
	cobra.OnInitialize(initConfig)

	genCmd.AddCommand(repoCmd)
	genCmd.AddCommand(structCmd)

	rootCmd.AddCommand(listTableCmd)
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	viper.SetConfigName("dbconfig")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: :%s \n", err))
	}
}

func initDatabase() (db *sqlx.DB, node sqalx.Node, err error) {
	connStr := viper.GetString("development.datasource")
	dialect := viper.GetString("development.dialect")

	db, err = sqlx.Open(dialect, connStr)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	node, err = sqalx.New(db)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	return
}

func generateRepo(node sqalx.Node, tableName string) (err error) {

	database := viper.GetString("development.database")
	columns, err := lib.GetTableColumnDefs(node, tableName, database)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	if len(columns) == 0 {
		log.Fatal("Could not find column info of table " + tableName)
	}

	f := NewFile("repo")

	// Generate Struct
	structCode, err := lib.GenerateStruct(tableName, columns)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	f.Add(structCode)
	f.Add(Line())

	interfaceCode, err := lib.GenInterface(tableName)
	if err != nil {
		err = tracerr.Wrap(err)
		return
	}

	f.Add(interfaceCode)
	f.Add(Line())

	singularName := inflection.Singular(tableName)
	pascalName := kace.Pascal(singularName)
	f.Add(Type().Id(pascalName + "Repository").Struct())
	f.Add(Line())

	codeQuestAll, _ := lib.GenerateGetAll(tableName, columns)
	f.Add(Comment("GetAll get all records"))
	f.Add(codeQuestAll)
	f.Add(Line())

	codeGetByID, _ := lib.GenerateGetByID(tableName, columns)
	f.Add(Comment("GetByID get a record by ID"))
	f.Add(codeGetByID)
	f.Add(Line())

	codeInsert, _ := lib.GenerateInsert(tableName, columns)
	f.Add(Comment("Insert insert a new record"))
	f.Add(codeInsert)
	f.Add(Line())

	codeUpdate, _ := lib.GenerateUpdate(tableName, columns)
	f.Add(Comment("Update update a exist record"))
	f.Add(codeUpdate)
	f.Add(Line())

	codeDelete, _ := lib.GenerateDelete(tableName, columns)
	f.Add(Comment("Delete logic delete a exist record"))
	f.Add(codeDelete)
	f.Add(Line())

	// codeBatchDelete := lib.GenerateBatchDelete(tableName, columns)
	// f.Add(Comment("BatchDelete logic batch delete records"))
	// f.Add(codeBatchDelete)
	// f.Add(Line())

	os.MkdirAll(path.Join("./repos", singularName), os.ModePerm)
	os.MkdirAll(path.Join("./repos", singularName, "sql", singularName), os.ModePerm)

	codeFileName := path.Join("./repos", singularName, fmt.Sprintf("%s_repo.go", singularName))
	err = f.Save(codeFileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}
