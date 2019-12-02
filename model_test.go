package gourm

import (
	"testing"
	"time"
)

type TestStruct struct {
	Model     `table:"sw_test" primary_key:"id"`
	ID        int    `col:"id"`
	Name      string `col:"name"`
	Age       int    `col:"age"`
	LastLogin string `col:"lastlogin"`
}
type Test2Struct struct {
	Model `table:"sw_test"`
	ID    string `col:"id"`
	Name  string `col:"name"`
	Admin int    `col:"admin"`
	Class int    `col:"class"`
}

// func TestUpdate(t *testing.T) {
// 	zyp := &TestStruct{ID: 9119, Age: 18}

// 	err := new(DB).Update(zyp)
// 	err = new(DB).Select("name", "id").Update(zyp)
// 	err = new(DB).Where("lastlogin = ?", time.Now().Format("2006/01/02 15:04:05")).Select("name").Update(zyp)

// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestInsert(t *testing.T) {
	zyp := &TestStruct{ID: 9, Name: "赵育鹏", Age: 28, LastLogin: time.Now().Format("2006/01/02 15:04:05")}
	err := new(DB).Insert(zyp)
	if err != nil {
		t.Error(err)
	}

	sjp := &Test2Struct{ID: "b23-37-01", Name: "sjp", Admin: 9, Class: 23}
	err = new(DB).Insert(sjp)
	if err != nil {
		t.Error(err)
	}
}
