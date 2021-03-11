package sqli

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Contrast-Security-OSS/go-test-bench/utils"
	// database import for sqlite3
	_ "github.com/mattn/go-sqlite3"

	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var templates = template.Must(template.ParseFiles(
	"./views/partials/safeButtons.gohtml",
	"./views/pages/sqlInjection.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

func sqliTemplate(w http.ResponseWriter, r *http.Request, params utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "sqlInjection", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true

}

func headersHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, splitURL []string) (template.HTML, bool) {
	// Currently we pass only json credentials in headers for SQL Injection
	if splitURL[3] != "json" {
		return template.HTML("INVALID URL"), false
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal([]byte(r.Header.Get("credentials")), &credentials)
	if err != nil {
		log.Printf("Could not parse headers to json, err = %s", err)
	}

	// setting up a database to execute the built query
	var sqlite3Database *sql.DB
	sqlite3Database, err = setupSqlite3()
	if err != nil {
		return template.HTML(err.Error()), false
	}
	defer func() {
		_ = sqlite3Database.Close()
	}()

	query := getSqliteSafetyQuery(sqlite3Database, credentials.Username, splitURL[4])

	_ = os.Remove("tempDatabase.db")
	return template.HTML(query), false //change to out desired output
}

// setupSqlite3 helper function to initialize the database
// closing db connection is executed in the context where it's used
func setupSqlite3() (*sql.DB, error) {
	_ = os.Remove("tempDatabase.db")
	log.Println("Creating tempDatabase.db...")
	file, err := os.Create("tempDatabase.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_ = file.Close()
	log.Println("tempDatabase.db created")
	sqlite3Database, _ := sql.Open("sqlite3", "./tempDatabase.db")

	return sqlite3Database, nil
}

func getSqliteSafetyQuery(db *sql.DB, userInput, safety string) string{
	var query string
	switch safety {
	case "unsafe":
		query = fmt.Sprintf("SELECT '%s' as '%s'", userInput, "test")
		res, err := db.Exec(query)
		log.Println("Result: ", res, " Error: ", err)
	case "safe":
		// Safe uses a parameterized query which is built by exec from
		// parameters which are passed in along with a static query string
		query = "SELECT '?' as '?'"
		res, err := db.Exec(query, userInput, "test")
		log.Println("Result: ", res, " Error: ", err)
	default:
		log.Println(safety)
	}

	return query
}

func sqlite3Handler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, splitURL []string) (template.HTML, bool) {
	sqlite3Database, err := setupSqlite3()
	if err != nil {
		return template.HTML(err.Error()), false
	}
	defer func() {
		_ = sqlite3Database.Close()
	}()

	userInput := utils.GetUserInput(r)
	query := getSqliteSafetyQuery(sqlite3Database, userInput, splitURL[4])

	_ = os.Remove("tempDatabase.db")
	return template.HTML(query), false //change to out desired output
}

// Handler is the API handler for sql injection
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 4 {
		return sqliTemplate(w, r, pd)
	}
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}
	if splitURL[2] == "headers" {
		return headersHandler(w, r, pd.Rulebar[pd.Name], splitURL)
	}

	switch splitURL[3] {
	case "sqlite3Exec":
		return sqlite3Handler(w, r, pd.Rulebar[pd.Name], splitURL)
	case "":
		return sqliTemplate(w, r, pd)
	default:
		log.Fatal("sqlInjection Handler reached incorrectly")
		return "", false
	}
}
