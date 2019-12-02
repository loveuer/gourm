package urm

import "reflect"

import "fmt"

func (d *DB) Delete(val ...interface{}) error {
	one := val[0]
	if reflect.TypeOf(one).Kind() != reflect.Ptr {
		return fmt.Errorf("delete target not a ptr")
	}
	if reflect.TypeOf(one).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("delete target not a struct")
	}

	rflctVal := reflect.ValueOf(one).Elem()
	tablename := getTableName(rflctVal)

	if len(val) == 1 {
		primaryKey := getPrimaryKey(rflctVal)
		if primaryKey == "" {
			return fmt.Errorf("no target imf")
		}
		primaryVal := getPrimaryVal(rflctVal)
		sentence := fmt.Sprintf("delete from %s where %s = '%v'", tablename, primaryKey, primaryVal)
		_, err = dbConn.Exec(sentence)
		return err
	} else if len(val) == 3 {
		sentence := fmt.Sprintf("delete from %s where %s = '%v'", tablename, val[1], val[2])
		_, err = dbConn.Exec(sentence)
		return err
	}

	return fmt.Errorf("length of params err")
}
