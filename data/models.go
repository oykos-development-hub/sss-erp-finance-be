package data

import (
	"fmt"

	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"

	"database/sql"
	"os"
)

//nolint:all
var db *sql.DB

//nolint:all
var upper db2.Session

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		upper, _ = mysql.New(databasePool)
	case "postgres", "postgresql":
		upper, _ = postgresql.New(databasePool)
	default:
		// do nothing
	}

	return Models{}
}

//nolint:all
func getInsertId(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
