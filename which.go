package gourm

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (d *DB) Which(out interface{}, outPrimaryVal ...interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("out not a struct pointer")
	}

	rflctVal := reflect.ValueOf(out).Elem()
	columns := getColumns(rflctVal)
	primaryKey := getPrimaryKey(rflctVal)
	var condition string

	if len(outPrimaryVal) == 0 {
		primaryVal := getColumnVal(rflctVal, primaryKey)
		condition = fmt.Sprintf("%s = '%v'", primaryKey, primaryVal)

	} else if len(outPrimaryVal) == 1 {
		condition = fmt.Sprintf("%s = '%v'", primaryKey, outPrimaryVal[0])

	} else if len(outPrimaryVal) == 2 {
		if reflect.TypeOf(outPrimaryVal[0]).Kind() != reflect.String {
			return fmt.Errorf("gourm which: col val must be string but => %v", reflect.TypeOf(outPrimaryVal[0]))
		}
		condition = fmt.Sprintf("%s = '%v'", outPrimaryVal[0].(string), outPrimaryVal[1])

	} else {
		return fmt.Errorf("primary key length not 0 or 1 but %d", len(outPrimaryVal))
	}

	var sliceout []interface{}

	sentence := fmt.Sprintf("select %s from %s where %s limit 1", strings.Join(columns, ", "), getTableName(rflctVal), condition)
	for i := 1; i < rflctVal.NumField(); i++ {
		sliceout = append(sliceout, rflctVal.Field(i).Addr().Interface())
	}

	row := dbConn.QueryRow(sentence)
	err := row.Scan(sliceout...)
	if err != nil {
		log.Printf("<model><model><Which>scan with sentence: %s err: %v\n", sentence, err)
		return err
	}

	return nil
}
