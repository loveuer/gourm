package gourm

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (m *Model) Update(uv interface{}, cds ...interface{}) error {
	var q = &queryCondition{}
	if len(cds) == 0 {
		rv := reflect.ValueOf(uv).Elem()
		primaryKey := getPrimaryKey(rv)
		primaryVal := getPrimaryVal(rv)
		q.Where(fmt.Sprintf("%s = ?", primaryKey), primaryVal)
	} else if len(cds) == 2 {
		if reflect.TypeOf(cds[0]).Kind() != reflect.String {
			return fmt.Errorf("gourm update: column must be string but => %v", reflect.TypeOf(cds[0]))
		}
		q.Where(fmt.Sprintf("%s = '%v'", cds[0], cds[1]))
	} else {
		return fmt.Errorf("gourm update: update where condition must be 0(use primary key) or 2(column and val)")
	}

	return q.DoUpdate(uv)
}

func (q *queryCondition) Update(uv interface{}) error {
	rv := reflect.ValueOf(uv).Elem()
	primaryKey := getPrimaryKey(rv)
	primaryVal := getPrimaryVal(rv)
	if primaryKey != "" {
		switch primaryVal.(type) {
		case int:
			if primaryVal.(int) != 0 {
				q.Where(fmt.Sprintf("%s = ?", primaryKey), primaryVal)
			}
		case string:
			if primaryVal.(string) != "" {
				q.Where(fmt.Sprintf("%s = ?", primaryKey), primaryVal)
			}
		}
	}

	return q.DoUpdate(uv)
}

func (q *queryCondition) DoUpdate(uv interface{}) error {
	rv := reflect.ValueOf(uv).Elem()
	tablename := getTableName(rv)
	if tablename == "" {
		return fmt.Errorf("no table name")
	}

	var (
		chgstr string
		chgs   []string
	)
	if len(q.selects) == 0 {
		colvals := getColumnVals(rv)
		for i, k := range colvals {
			chgs = append(chgs, fmt.Sprintf("%s = '%v'", i, k))
		}
	} else {
		for _, k := range q.selects {
			chgs = append(chgs, fmt.Sprintf("%s = '%v'", k, getColumnVal(rv, k)))
		}
	}

	chgstr = strings.Join(chgs, ", ")

	whereConditions := conditions2sentence(q.where)
	sentence := fmt.Sprintf("update %s set %s where %s", tablename, chgstr, strings.Join(whereConditions, ", "))

	_, err := dbConn.Exec(sentence)
	if err != nil {
		log.Printf("<urm><update> with sentence: (%s) err: %v\n", sentence, err)
	}

	return err
}
