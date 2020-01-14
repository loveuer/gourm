package gourm

import (
	"fmt"
	"strings"
)

func (q *queryCondition) All(out interface{}) error {
	q.Where("1 = ?", 1)
	sentence, err := q.querySlice(out)
	if err != nil {
		return fmt.Errorf("gourm all: condition to setence err => %v", err)
	}

	sentence = fmt.Sprintf("select %s %s", strings.Join(q.selects, ", "), sentence)
	return setVals(out, sentence)
}

func (d *DB) All(out interface{}) error {
	d.query = new(queryCondition)
	return d.query.All(out)
}
