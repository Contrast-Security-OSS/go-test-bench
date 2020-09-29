package sqli

import (
	"bytes"
	"database/sql"
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

func sqliTemplate(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "sqlInjection", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true

}

func sqlite3Handler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, splitURL []string) (template.HTML, bool) {
	query := fmt.Sprintf("SELECT '%s' as '%s'", r.URL.Query().Get("input"), "test")
	_ = os.Remove("tempDatabase.db")
	log.Println("Creating tempDatabase.db...")
	file, err := os.Create("tempDatabase.db")
	if err != nil {
		log.Println(err)
		return template.HTML(err.Error()), false
	}
	_ = file.Close()
	log.Println("tempDatabase.db created")
	sqlite3Database, _ := sql.Open("sqlite3", "./tempDatabase.db")

	defer func() {
		_ = sqlite3Database.Close()
	}()

	switch splitURL[4] {
	case "unsafe":
		res, err := sqlite3Database.Exec(query)
		log.Println("Result: ", res, " Error: ", err)
	case "safe":
		// Safe uses a parameterized query which is built by exec from
		// parameters which are passed in along with a static query string
		query := "SELECT '?' as '?'"
		res, err := sqlite3Database.Exec(query, r.URL.Query().Get("input"), "test")
		log.Println("Result: ", res, " Error: ", err)
	default:
		log.Println(splitURL[4])
	}
	_ = os.Remove("tempDatabase.db")
	return template.HTML(query), false //change to out desired output
}

// Handler is the API handler for sql injection
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 4 {
		return sqliTemplate(w, r, pd.Rulebar[pd.Name])
	}
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}
	switch splitURL[3] {
	case "sqlite3Exec":
		return sqlite3Handler(w, r, pd.Rulebar[pd.Name], splitURL)
	case "":
		return sqliTemplate(w, r, pd.Rulebar[pd.Name])
	default:
		log.Fatal("sqlInjection Handler reached incorrectly")
		return "", false
	}
}
