package code_generator

import (
	"fmt"
	"reflect"
	"strconv"
)

// ColumnOp 列属性
type ColumnOp struct {
	Comment    string
	ColumnName string
	FiledName  string
	Type       string
	GormTag    string
}

var Group = func(columns []*ColumnOp, subGroupLength int64) [][]*ColumnOp {

	max := int64(len(columns))

	var segments = make([][]*ColumnOp, 0)
	quantity := max / subGroupLength
	remainder := max % subGroupLength
	for i := int64(0); i < quantity; i++ {
		segments = append(segments, columns[i*subGroupLength:(i+1)*subGroupLength])
	}
	if quantity == 0 || remainder != 0 {
		segments = append(segments, columns[quantity*subGroupLength:quantity*subGroupLength+remainder])
	}
	return segments
}

var GetMaxLength = func(columns []*ColumnOp, filed string) int {
	maxLength := 0
	for _, column := range columns {
		v := reflect.ValueOf(column).Elem()
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			name := t.Field(i).Name
			if name == filed {
				val := v.Field(i).Interface().(string)
				maxLength = len(val)
			}
		}
	}
	// 默认值最大长度加1
	return maxLength + 1
}
var SortMap = func(index int) string {
	m := map[int]string{
		0: "One",
		1: "Two",
		2: "Three",
		3: "Four",
		4: "Five",
		5: "Six",
		6: "Seven",
		7: "Eight",
		8: "Nine",
		9: "Ten",
	}
	indexEn, ok := m[index]
	if ok {
		return indexEn
	}
	return strconv.Itoa(index)
}

var FillBlank = func(field string, max int) string {
	if len(field) >= max {
		return field
	}
	return fmt.Sprintf("%-*s", max, field)
}
