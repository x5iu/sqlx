package sqlx

import (
	"strconv"
	"strings"
	"sync"

	defc "github.com/x5iu/defc/runtime"
)

// Bindvar types supported by Rebind, BindMap and BindStruct.
const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

var defaultBinds = map[int][]string{
	DOLLAR:   []string{"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"},
	QUESTION: []string{"mysql", "sqlite3", "nrmysql", "nrsqlite3"},
	NAMED:    []string{"oci8", "ora", "goracle", "godror"},
	AT:       []string{"sqlserver"},
}

var binds sync.Map

func init() {
	for bind, drivers := range defaultBinds {
		for _, driver := range drivers {
			BindDriver(driver, bind)
		}
	}

}

// BindType returns the bind type for a given database given a driver name.
func BindType(driverName string) int {
	itype, ok := binds.Load(driverName)
	if !ok {
		return UNKNOWN
	}
	return itype.(int)
}

// BindDriver sets the BindType for driverName to bindType.
func BindDriver(driverName string, bindType int) {
	binds.Store(driverName, bindType)
}

// Rebind a query from the default bind type (QUESTION) to the target bind type.
func Rebind(bindType int, query string) string {
	switch bindType {
	case QUESTION, UNKNOWN:
		return query
	default:
	}
	tokens := defc.SplitTokens(query)
	targetQuery := make([]string, 0, len(tokens))
	var j int
	for _, token := range tokens {
		j++
		switch token {
		case "?":
			switch bindType {
			case DOLLAR:
				targetQuery = append(targetQuery, "$"+strconv.Itoa(j))
			case NAMED:
				targetQuery = append(targetQuery, ":arg"+strconv.Itoa(j))
			case AT:
				targetQuery = append(targetQuery, "@p"+strconv.Itoa(j))
			default:
				panic("unknown bind type")
			}
		default:
			targetQuery = append(targetQuery, token)
		}
	}
	return strings.Join(targetQuery, " ")
}

// In expands slice values in args, returning the modified query string
// and a new arg list that can be executed by a database. The `query` should
// use the `?` bindVar.  The return value uses the `?` bindVar.
var In = defc.In[[]any]
