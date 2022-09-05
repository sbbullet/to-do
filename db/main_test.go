package db

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/sbbullet/to-do/util"
)

var testStore *Store
var testDB *sql.DB

func TestMain(m *testing.M) {
	config := util.LoadConfig("app", "env", "..")
	config.DBSource = "../test_todo.db"

	if _, err := os.Stat(config.DBSource); !errors.Is(err, os.ErrNotExist) {
		// file exists
		os.Remove(config.DBSource)
	}

	testDB = NewDB(config)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
