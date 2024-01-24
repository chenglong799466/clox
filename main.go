package main

import (
	"colx/code_generator"
	"colx/config"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	tableName = flag.String("table", "111", "input generate table")
	//dbUser    = "jacklcheng"
	//dbPass    = "Jacklcheng@2024"
	//dbHost    = "43.153.70.226"
	//dbPort    = "3306"
	//dbName    = "qcloud_market"
)

// CreateTable represents the result of "show create table" statement.
type CreateTable struct {
	TableName   string `gorm:"column:Table"`
	CreateTable string `gorm:"column:Create Table"`
	CreateView  string `gorm:"column:Create View"` // View support
	ViewName    string `gorm:"column:View"`        // View support
}

func main() {
	flag.Parse()
	if len(*tableName) == 0 {
		log.Fatalf("table name is empty")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbConfig.User, config.DbConfig.PWD, config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.DataBase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	showCreateTable := fmt.Sprintf("show create table %s;", *tableName)
	var res CreateTable
	if err = db.Raw(showCreateTable).Find(&res).Error; err != nil {
		log.Fatalf("failed to retrieve create table statement: %v", err)
	}

	if res.CreateView != "" || res.ViewName != "" {
		log.Fatal("views are not supported")
	}

	astNode, err := parse(res.CreateTable)
	if err != nil {
		log.Fatalf("failed to parse create table statement: %v", err)
	}

	fmt.Printf("%v\n", *astNode)

	info := &code_generator.TemplateData{}
	(*astNode).Accept(info)

	templateFile := "./code_generator/generate.tmpl"
	outputFile := fmt.Sprintf("./%s.go", info.TableName)

	if err = generateCodeFile(info, templateFile, outputFile); err != nil {
		log.Fatalf("failed to generate code file: %v", err)
	}

	fmt.Printf("Code file generated: %s\n", outputFile)
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.ParseSQL(sql)
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func generateCodeFile(info *code_generator.TemplateData, templateFile, outputFile string) error {

	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	var funcMap = template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"group":        code_generator.Group,
		"sortMap":      code_generator.SortMap,
		"getMaxLength": code_generator.GetMaxLength,
		"fillBlank":    code_generator.FillBlank,
	}

	generator, err := code_generator.NewCodeGenerator(info, templateContent, funcMap)
	if err != nil {
		return fmt.Errorf("failed to generate code file: %w", err)
	}
	if err = generator.Generate(outputFile); err != nil {
		return fmt.Errorf("failed to generate code file: %w", err)
	}
	return nil
}
