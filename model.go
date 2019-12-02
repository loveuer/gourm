package gourm

import (
	"reflect"
)

func getTableName(v reflect.Value) string {
	return v.Type().Field(0).Tag.Get("table")
}
func getPrimaryKey(v reflect.Value) string {
	return v.Type().Field(0).Tag.Get("primary_key")
}
func getPrimaryVal(v reflect.Value) interface{} {
	primary := getPrimaryKey(v)
	for i := 1; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("col") == primary {
			return v.Field(i).Interface()
		}
	}

	return nil
}
func getPrimaryAddr(v reflect.Value) interface{} {
	primary := getPrimaryKey(v)
	for i := 1; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("col") == primary {
			return v.Field(i).Addr().Interface()
		}
	}
	return nil
}
func getColumns(v reflect.Value) []string {
	var cols []string
	for i := 1; i < v.NumField(); i++ {
		cols = append(cols, v.Type().Field(i).Tag.Get("col"))
	}
	return cols
}
func getColumnVals(v reflect.Value) map[string]interface{} {
	r := make(map[string]interface{})
	for i := 1; i < v.NumField(); i++ {
		r[v.Type().Field(i).Tag.Get("col")] = v.Field(i).Interface()
	}

	return r
}
func getColumnVal(v reflect.Value, col string) interface{} {
	for i := 1; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("col") == col {
			return v.Field(i).Interface()
		}
	}
	return nil
}

func (q *queryCondition) Where(query interface{}, values ...interface{}) *queryCondition {
	q.where = append(q.where, map[string]interface{}{"query": query, "args": values})
	return q
}

func (q *queryCondition) Or(query interface{}, values ...interface{}) *queryCondition {
	q.or = append(q.or, map[string]interface{}{"query": query, "args": values})
	return q
}

func (q *queryCondition) And(query interface{}, values ...interface{}) *queryCondition {
	q.and = append(q.and, map[string]interface{}{"query": query, "args": values})
	return q
}

func (q *queryCondition) Select(query ...interface{}) *queryCondition {
	if len(query) == 1 && reflect.TypeOf(query[0]).Kind() == reflect.Slice {
		if reflect.TypeOf(query[0]).Elem().Kind() == reflect.String {
			for _, v := range query[0].([]string) {
				q.selects = append(q.selects, v)
			}
		}

		return q
	}

	for _, v := range query {
		switch v.(type) {
		case string:
			q.selects = append(q.selects, v.(string))
		}
	}
	return q
}

func (q *queryCondition) Order(query string) *queryCondition {
	q.order = query
	return q
}

func (q *queryCondition) Offset(query int) *queryCondition {
	q.offset = query
	return q
}

func (q *queryCondition) Limit(query int) *queryCondition {
	q.limit = query
	return q
}

func (d *DB) Where(query interface{}, values ...interface{}) *queryCondition {
	d.query = new(queryCondition)
	return d.query.Where(query, values...)
}
func (d *DB) Select(query ...interface{}) *queryCondition {
	d.query = new(queryCondition)
	return d.query.Select(query...)
}
