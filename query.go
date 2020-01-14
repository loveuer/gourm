package gourm

import (
	"fmt"
	"reflect"
	"strings"
)

// 把where接收到的转换成语句的slice
func conditions2sentence(where []map[string]interface{}) []string {
	var ws []string

	for _, onew := range where {
		var ones string
		for i := 0; i < reflect.ValueOf(onew["args"]).Len(); i++ {
			onev := reflect.ValueOf(onew["args"]).Index(i).Interface()
			ones = onew["query"].(string)
			switch onev.(type) {
			case int:
				ri := fmt.Sprintf("%d", onev)
				ones = strings.Replace(ones, "?", ri, 1)
			case string:
				ro := fmt.Sprintf("'%s'", onev)
				ones = strings.Replace(ones, "?", ro, 1)
			default:
				if reflect.TypeOf(onev).Kind() == reflect.Slice {
					var vv []string
					for j := 0; j < reflect.ValueOf(onev).Len(); j++ {
						vv = append(vv, fmt.Sprintf("'%v'", reflect.ValueOf(onev).Index(j).Interface()))
					}
					ones = strings.Replace(ones, "?", strings.Join(vv, ", "), 1)
				} else {
					ones = strings.Replace(ones, "?", fmt.Sprintf("'%v'", onev), 1)
				}
			}
		}
		ws = append(ws, ones)
	}
	return ws
}

func (q *queryCondition) query(out interface{}) (string, error) {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return "", fmt.Errorf("out not a pointer")
	}

	if reflect.TypeOf(out).Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("out not a struct")
	}

	wheres := conditions2sentence(q.where)
	if len(wheres) == 0 {
		return "", fmt.Errorf("no Where condition")
	}

	rflctVal := reflect.ValueOf(out).Elem()

	tablename := getTableName(rflctVal)
	if tablename == "" {
		return "", fmt.Errorf("no table name")
	}

	ands := conditions2sentence(q.and)
	ors := conditions2sentence(q.or)

	sentence := fmt.Sprintf("from %s where %s", tablename, strings.Join(wheres, " and "))

	if len(ands) != 0 {
		sentence = fmt.Sprintf("%s and %s", sentence, strings.Join(ands, " and "))
	}
	if len(ors) != 0 {
		sentence = fmt.Sprintf("%s or %s", sentence, strings.Join(ors, " or "))
	}

	return sentence, nil
}

func (q *queryCondition) querySlice(out interface{}) (string, error) {
	// check if out is slice
	if reflect.TypeOf(out).Elem().Kind() != reflect.Slice {
		return "", fmt.Errorf("out not a slice")
	}
	// check if out is ptr
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return "", fmt.Errorf("out not a ptr")
	}
	// check if out is struct
	if reflect.TypeOf(out).Elem().Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("out[] not a struct")
	}

	wheres := conditions2sentence(q.where)
	if len(wheres) == 0 {
		return "", fmt.Errorf("no Where condition")
	}

	ands := conditions2sentence(q.and)
	ors := conditions2sentence(q.or)

	onetype := reflect.TypeOf(out).Elem().Elem()
	newone := reflect.New(onetype)
	if len(q.selects) == 0 {
		for _, c := range getColumns(newone.Elem()) {
			q.selects = append(q.selects, c)
		}
	}

	tablename := getTableName(newone.Elem())
	if tablename == "" {
		return "", fmt.Errorf("no table name")
	}

	sentence := fmt.Sprintf("from %s where %s", tablename, strings.Join(wheres, " and "))

	if len(ands) != 0 {
		sentence = fmt.Sprintf("%s and %s", sentence, strings.Join(ands, " and "))
	}
	if len(ors) != 0 {
		sentence = fmt.Sprintf("%s or %s", sentence, strings.Join(ors, " or "))
	}
	if q.order != "" {
		sentence = fmt.Sprintf("%s order by %s", sentence, q.order)
	}
	if q.offset != 0 {
		sentence = fmt.Sprintf("%s offset %d", sentence, q.offset)
	}
	if q.limit != 0 {
		sentence = fmt.Sprintf("%s limit %d", sentence, q.limit)
	}

	return sentence, nil
}

func (q *queryCondition) Count(out interface{}, col ...string) (int, error) {
	sentence, err := q.query(out)
	if err != nil {
		return -1, err
	}

	var countcol string
	if len(col) == 0 {
		countcol = getPrimaryKey(reflect.ValueOf(out).Elem())
		if countcol == "" {
			return -1, fmt.Errorf("no primary key")
		}
	} else {
		countcol = col[0]
	}

	sentence = fmt.Sprintf("select count(%s) %s", countcol, sentence)

	row := dbConn.QueryRow(sentence)
	var count int
	err = row.Scan(&count)
	return count, err
}

func (q *queryCondition) Find(out interface{}) error {

	sentence, err := q.querySlice(out)
	if err != nil {
		return err
	}

	sentence = fmt.Sprintf("select %s %s", strings.Join(q.selects, ", "), sentence)

	return setVals(out, sentence)
}

func setVals(out interface{}, sentence string) error {
	rows, err := dbConn.Query(sentence)
	if err != nil {
		return fmt.Errorf("gourm query: query sql err => %v", err)
	}

	onetype := reflect.TypeOf(out).Elem().Elem()
	for rows.Next() {
		one := reflect.New(onetype)
		var oneFields []interface{}
		selectedCols, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("gourm setvals: get selected columns err => %v", err)
		}

		fieldnums := one.Elem().NumField()
		for _, col := range selectedCols {
			for i := 1; i < fieldnums; i++ {
				if one.Elem().Type().Field(i).Tag.Get("col") == col {
					oneFields = append(oneFields, one.Elem().Field(i).Addr().Interface())
				}
			}
		}

		err = rows.Scan(oneFields...)
		if err != nil {
			return fmt.Errorf("gourm setvals rows scan err => %v", err)
		}

		outv := reflect.ValueOf(out).Elem()
		outv.Set(reflect.Append(outv, one.Elem()))
	}

	return nil
}
