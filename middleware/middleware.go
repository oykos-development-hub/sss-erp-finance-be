package middleware

import (
	"gitlab.sudovi.me/erp/finance-api/data"

	"github.com/oykos-development-hub/celeritas"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
