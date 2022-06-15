package sqli

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"

	// database import for sqlite3
	_ "github.com/mattn/go-sqlite3"
)

const mime = "text/plain"

// RegisterRoutes is called in framework init to register routes in this package.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := []*common.Sink{
		{
			Name:                 "sqlite3.exec",
			Handler:              sqliteInj{}.execHandler,
			ExpectedUnsafeStatus: http.StatusBadRequest,
		},
	}
	sinks = append(sinks, frameworkSinks...)
	common.Register(common.Route{
		Name:     "SQL Injection",
		Link:     "https://www.owasp.org/index.php/SQL_Injection",
		Base:     "sqlInjection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"body", "query", "headers-json"},
		Sinks:    sinks,
		Payload:  "Robert'; DROP TABLE Students;--",
	})
}

type sqliteInj struct {
	path string
	db   *sql.DB
}

func (si sqliteInj) execHandler(mode common.Safety, in string, _ interface{}) (string, string, int) {
	log.Println("sqlite exec handler")
	var err error
	var res sql.Result

	if err = si.initDB(); err != nil {
		return err.Error(), mime, http.StatusBadRequest
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
		return "NOOP", mime, http.StatusOK
	}

	if err != nil {
		return err.Error(), mime, http.StatusBadRequest
	}
	r := fmt.Sprintf("Result: %#v\n", res)
	return r, mime, http.StatusOK
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
