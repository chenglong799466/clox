package code_generator

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/pingcap/tidb/pkg/parser/types"
)

// TemplateData 自定义visitor访问Node,
type TemplateData struct {
	PackageName string
	Imports     []string
	TableName   string // 从建表语句获取
	StructName  string // 通过表名转换而来
	Comment     string // 表的comment
	Columns     []*ColumnOp
}

func (v *TemplateData) Enter(in ast.Node) (ast.Node, bool) {

	createStmt, ok := in.(*ast.CreateTableStmt)
	if !ok {
		fmt.Println("templateInfo enter is not ok")
		return nil, false
	}

	// TableName
	v.TableName = createStmt.Table.Name.O
	// StructName
	var fieldStyle string
	splits := strings.Split(v.TableName, "_")
	for _, split := range splits {
		fieldStyle = fieldStyle + strings.ToUpper(split[:1]) + split[1:]
	}
	v.StructName = fieldStyle
	// Comment
	for _, option := range createStmt.Options {
		switch option.Tp {
		case ast.TableOptionComment:
			v.Comment = option.StrValue
		}
	}
	// Columns
	m := make(map[string]string)
	for _, cons := range createStmt.Constraints {
		switch cons.Tp {
		case ast.ConstraintPrimaryKey:
			if len(cons.Keys) > 0 {
				m[cons.Keys[0].Column.Name.O] = "primaryKey"
			}
		}
	}
	var columnOps []*ColumnOp
	for _, col := range createStmt.Cols {
		colName := col.Name.OrigColName()
		// comment,tag
		var comment string
		tag := GormTag{}
		for _, option := range col.Options {
			switch option.Tp {
			case ast.ColumnOptionNotNull:
				tag.NotNull = true
			case ast.ColumnOptionComment:
				if val, ok := option.Expr.(*test_driver.ValueExpr); ok {
					comment = val.Datum.GetString()
				}
			case ast.ColumnOptionAutoIncrement:
				tag.AutoIncrement = true
			case ast.ColumnOptionDefaultValue:
				if val, ok := option.Expr.(*test_driver.ValueExpr); ok {
					tag.Default = val.Datum.GetString()
				}

			}
		}
		// primaryKey
		_, ok = m[tag.Column]
		if ok {
			tag.PrimaryKey = true
		}

		// gormTagStr
		tag.Type = col.Tp.InfoSchemaStr()
		tag.Column = colName
		gormTagStr := tag.fillGormTag()
		//
		columnOp := ColumnOp{
			Comment:    comment,
			FiledName:  TransferDbToFieldStyle(colName),
			ColumnName: colName,
			Type:       transferMysqlType2GoType(col.Tp, tag.NotNull),
			GormTag:    gormTagStr,
		}
		fmt.Println(columnOp)
		columnOps = append(columnOps, &columnOp)
	}
	v.Columns = columnOps
	// 跳过子节点
	return in, true
}

func (v *TemplateData) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

// mysql类型转成go类型
func transferMysqlType2GoType(col *types.FieldType, isNotNull bool) string {
	switch col.GetType() {
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong,
		mysql.TypeBit, mysql.TypeYear:
		if !isNotNull {
			return "sql.NullInt64"
		}
		if mysql.HasUnsignedFlag(col.GetFlag()) {
			return "uint64"
		}
		return "int64"
	case mysql.TypeFloat, mysql.TypeDouble:
		if !isNotNull {
			return "sql.NullFloat64"
		}
		return "float64"
	case mysql.TypeNewDecimal:
		if !isNotNull {
			return "sql.NullFloat64"
		}
		return "float64"
	case mysql.TypeDate, mysql.TypeDatetime:
		if !isNotNull {
			return "sql.NullTime"
		}
		return "time.Time"
	case mysql.TypeTimestamp:
		if !isNotNull {
			return "sql.NullTime"
		}
		return "time.Time"
	case mysql.TypeDuration:
		if !isNotNull {
			return "sql.NullTime"
		}
		return "time.Time"
	}
	// 如果not null字段为空
	if !isNotNull {
		return "sql.NullString"
	}
	return "string"
}

/*
声明 model 时，tag 是可选的，GORM 支持以下 tag： tag 名大小写不敏感，但建议使用 camelCase 风格

标签名	说明
column	指定 db 列名
type	列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：not null、size, autoIncrement… 像 varbinary(8) 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：MEDIUMINT UNSIGNED not NULL AUTO_INSTREMENT
size	指定列大小，例如：size:256
primaryKey	指定列为主键
unique	指定列为唯一
default	指定列的默认值
precision	指定列的精度
scale	指定列大小
not null	指定列为 NOT NULL
autoIncrement	指定列为自动增长
embedded	嵌套字段
embeddedPrefix	嵌入字段的列名前缀
autoCreateTime	创建时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoCreateTime:nano
autoUpdateTime	创建 / 更新时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoUpdateTime:milli
index	根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 索引 获取详情
uniqueIndex	与 index 相同，但创建的是唯一索引
check	创建检查约束，例如 check:age > 13，查看 约束 获取详情
<-	设置字段写入的权限， <-:create 只创建、<-:update 只更新、<-:false 无写入权限、<- 创建和更新权限
->	设置字段读的权限，->:false 无读权限
-	忽略该字段，- 无读写权限

————————————————
https://learnku.com/docs/gorm/v2/models/9729#15a96d
*/

// GormTag 生成的gorm的tag结构体
type GormTag struct {
	PrimaryKey    bool
	AutoIncrement bool
	NotNull       bool
	Column        string
	Type          string
	Default       string
}

func (g *GormTag) fillGormTag() string {
	if g == nil {
		return ""
	}

	var tagParts []string
	if g.PrimaryKey {
		tagParts = append(tagParts, "primaryKey")
	}
	if g.AutoIncrement {
		tagParts = append(tagParts, "autoIncrement")
	}
	if len(g.Column) > 0 {
		tagParts = append(tagParts, fmt.Sprintf("column:%s", g.Column))
	}
	if len(g.Type) > 0 {
		tagParts = append(tagParts, fmt.Sprintf("type:%s", g.Type))
	}
	if len(g.Default) > 0 {
		tagParts = append(tagParts, fmt.Sprintf("default:%s", g.Default))
	}
	if g.NotNull {
		tagParts = append(tagParts, "not null")
	}

	return strings.Join(tagParts, ";")
}

func TransferDbToFieldStyle(dbStyle string) string {
	var fieldStyle string
	splits := strings.Split(dbStyle, "_")
	for _, split := range splits {
		fieldStyle = fieldStyle + strings.ToUpper(split[:1]) + split[1:]
	}
	return fieldStyle
}
