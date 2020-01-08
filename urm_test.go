package gourm

import (
	"fmt"
	"testing"
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

type TestUser struct {
	Model    `table:"sw_user" primary_key:"id"`
	ID       int    `col:"id"`
	Username string `col:"username"`
	Realname string `col:"realname"`
	Password string `col:"password"`
}

func TestWhich(t *testing.T) {
	db, err := New("postgres", dbsetting, true)
	if err != nil {
		t.Errorf("new db err => %v\n", err)
	}

	nu := &TestUser{
		ID:       1,
		Username: "zhaoyupeng",
		Realname: "赵育鹏",
		Password: "1314lalala",
	}

	db.Update(nu, "username", nu.Username)

	fmt.Println(nu)
}
