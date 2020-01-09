package gourm

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *Model) Insert(insertRow interface{}) error {

	// 检验一下参数是不是struct pointer
	t := reflect.TypeOf(insertRow)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("insert value not a pointer")
	}
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("insert value not a struct")
	}

	v := reflect.ValueOf(insertRow).Elem()

	tablename := getTableName(v)
	// 检查一下 table 名称, 这个是一个必须的参数
	if tablename == "" {
		return fmt.Errorf("no table name")
	}

	var (
		primary string
		cols    []string
		vals    []string
	)

	primary = getPrimaryKey(v)

	colVals := getColumnVals(v)

	for i, k := range colVals {
		if i == primary {
			continue
		}
		cols = append(cols, i)
		vals = append(vals, fmt.Sprintf("'%v'", k))
	}

	sentence := fmt.Sprintf("insert into %s (%s) values (%s) ", tablename, strings.Join(cols, ", "), strings.Join(vals, ", "))

	if primary != "" {
		sentence += fmt.Sprintf("returning %s", primary)
		row := dbConn.QueryRow(sentence)
		primaryAddr := getPrimaryAddr(v)
		err = row.Scan(primaryAddr)
		if err != nil {
			return fmt.Errorf("gourm insert: query sql err => %v", err)
		}
		return nil
	} else {
		_, err := dbConn.Exec(sentence)
		if err != nil {
			return fmt.Errorf("gourm insert: exec sql err => %v", err)
		}
		return nil
	}
}
