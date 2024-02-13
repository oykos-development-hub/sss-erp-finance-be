package data

import (
	"fmt"

	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"

	"database/sql"
	"os"

	up "github.com/upper/db/v4"
)

//nolint:all
var db *sql.DB

//nolint:all
var upper up.Session

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application
	Invoice                Invoice
	Article                Article
	Budget                 Budget
	FinancialBudget        FinancialBudget
	FinancialBudgetLimit   FinancialBudgetLimit
	NonFinancialBudget     NonFinancialBudget
	NonFinancialBudgetGoal NonFinancialBudgetGoal
	Program                Program
	Activity               Activity
	GoalIndicator          GoalIndicator
	FilledFinancialBudget  FilledFinancialBudget
	BudgetRequest          BudgetRequest
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

	return Models{
		Invoice:                Invoice{},
		Article:                Article{},
		Budget:                 Budget{},
		FinancialBudget:        FinancialBudget{},
		FinancialBudgetLimit:   FinancialBudgetLimit{},
		NonFinancialBudget:     NonFinancialBudget{},
		NonFinancialBudgetGoal: NonFinancialBudgetGoal{},
		Program:                Program{},
		Activity:               Activity{},
		GoalIndicator:          GoalIndicator{},
		FilledFinancialBudget:  FilledFinancialBudget{},
		BudgetRequest:          BudgetRequest{},
	}
}

//nolint:all
func getInsertId(i up.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}

func paginateResult(res up.Result, page int, pageSize int) up.Result {
	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize

	// Apply pagination to the query
	res = res.Offset(offset).Limit(pageSize)

	return res
}
