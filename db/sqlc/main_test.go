package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/arun6783/go-postgress-k8s/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	var config utils.Config
	var err error

	config, err = utils.LoadConfig("../../.")

	if err != nil {
		panic(err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
