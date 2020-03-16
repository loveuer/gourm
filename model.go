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
	DEBUG  bool
)

func New(sqlType, config string, ping bool, debug bool) (*DB, error) {
	dbConn, err = sql.Open(sqlType, config)
	if err != nil {
		return new(DB), err
	}

	DEBUG = debug

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
