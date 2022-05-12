package sqli

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"

	// database import for sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// RegisterRoutes is called in framework init to register routes in this package.
func RegisterRoutes( /* framework - unused */ string) {
	common.Register(common.Route{
		Name:     "SQL Injection",
		Link:     "https://www.owasp.org/index.php/SQL_Injection",
		Base:     "sqlInjection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"body", "query", "headers-json"},
		Sinks: []common.Sink{
			{
				Name:    "sqlite3.exec",
				Method:  "GET",
				Handler: sqliteInj{}.execHandler,
			},
		},
		Payload: "Robert'; DROP TABLE Students;--",
	})
}

type sqliteInj struct {
	path string
	db   *sql.DB
}

func (si sqliteInj) execHandler(mode common.Safety, in string, _ interface{}) template.HTML {
	log.Println("sqlite exec handler")
	var err error
	var res sql.Result

	if err = si.initDB(); err != nil {
		return template.HTML(err.Error())
	}
	defer si.cleanupDB()

	switch mode {
	case common.Unsafe:
		query := fmt.Sprintf("SELECT '%s' as '%s'", in, "test")
		res, err = si.db.Exec(query)
	case common.Safe:
		// Safe uses a parameterized query which is built by exec from
		// parameters which are passed in along with a static query string
		query := "SELECT '?' as '?'"
		res, err = si.db.Exec(query, in, "test")
	default: // mode is no-op or invalid
		return template.HTML("NOOP")
	}

	if err != nil {
		return template.HTML(err.Error())
	}
	r := fmt.Sprintf("Result: %#v\n", res)
	log.Println("Result: ", r)
	return template.HTML(r)
}

func (si *sqliteInj) initDB() error {
	// setting up a database to execute the built query
	si.path = "tempDatabase.db"
	_ = os.Remove(si.path)
	log.Printf("Creating %s...", si.path)
	file, err := os.Create(si.path)
	if err != nil {
		log.Println(err)
		return err
	}
	file.Close()
	db, err := sql.Open("sqlite3", si.path)
	if err != nil {
		return err
	}
	si.db = db
	return nil
}

func (si *sqliteInj) cleanupDB() {
	_ = si.db.Close()
	_ = os.Remove(si.path)
}
