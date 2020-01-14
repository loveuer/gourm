package gourm

import (
	"fmt"
	"log"
	"testing"
)

type TestStruct struct {
	Model     `table:"sw_user" primary_key:"id"`
	ID        int    `col:"id"`
	Username  string `col:"username"`
	Realname  string `col:"realname"`
	LastLogin string `col:"lastlogin"`
}

func TestWhich(t *testing.T) {
	db, err := New("postgres", dbsetting, true)
	if err != nil {
		log.Fatal("ping db err: ", err)
	}

	// var users []TestStruct

	// err = db.Select("id", "realname", "lastlogin").All(&users)
	// if err != nil {
	// 	t.Error(err)
	// }

	// for _, v := range users {
	// 	fmt.Println("one user: ", v)
	// }

	var mu []TestStruct

	err = db.Where("1 = ?", 1).Select("id", "realname", "lastlogin").Order("lastlogin desc").Find(&mu)
	if err != nil {
		t.Error("find err ", err)
	}

	for _, mv := range mu {
		fmt.Println("mone: ", mv)
	}
}
