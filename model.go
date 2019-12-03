package gourm

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Model struct{}

type queryCondition struct {
	where   []map[string]interface{}
	or      []map[string]interface{}
	and     []map[string]interface{}
	selects []string
	order   string
	offset  int
	limit   int
}

type DB struct {
	Model
	query *queryCondition
}

var (
	dbConn *sql.DB
	err    error
)

func New(sqlType, config string, ping bool) (*DB, error) {
	dbConn, err = sql.Open(sqlType, config)
	if err != nil {
		return new(DB), err
	}

	if ping {
		ec := make(chan error)
		go checkdb(ec)
		timer := time.NewTimer(3 * time.Second)
		for {
			select {
			case err = <-ec:
				if err != nil {
					log.Fatalf("ping db(%s) err: %v\n", config, err)
				}
				log.Println("connect db success")
				return new(DB), nil

			case <-timer.C:
				return new(DB), fmt.Errorf("ping db(%s) time out", config)
			}
		}
	}

	return new(DB), nil
}
func checkdb(ec chan error) {
	err := dbConn.Ping()
	if err != nil {
		ec <- err
		return
	}

	ec <- nil
}
