package sqlInjection

import (
	"net/http"
	"html/template"
	"bytes"
	"strings"
	"log"
	"os"
	"fmt"
	"net/url"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml","./views/pages/sqlInjection.gohtml", "./views/partials/ruleInfo.gohtml"))

func defaultHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "sqlInjection", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true 

}

func sqlite3Handler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, splitUrl []string) (template.HTML, bool) {
	query := fmt.Sprintf("SELECT '%s' as '%s'",r.URL.Query().Get("input"), "test")
	os.Remove("tempDatabase.db")
	log.Println("Creating tempDatabase.db...")
	file, err := os.Create("tempDatabase.db")
	if err != nil {
		log.Println(err)
		return template.HTML(err.Error()), false
	}
	file.Close()
	log.Println("tempDatabase.db created")
	sqlite3Database, err := sql.Open("sqlite3", "./tempDatabase.db")
	defer sqlite3Database.Close()
	switch splitUrl[4]{
	case "unsafe":
		res, err := sqlite3Database.Exec(query)
		log.Println("Result: ", res, " Error: ", err)
	case "safe":
		query = url.QueryEscape(query)
		res, err := sqlite3Database.Exec(query)
		log.Println("Result: ", res, " Error: ", err)
	default:
		log.Println(splitUrl[4])
	}
	os.Remove("tempDatabase.db")
	return template.HTML(query), false //change to out desired output
}

func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitUrl := strings.Split(r.URL.Path, "/")
	if len(splitUrl) < 4 {
		return defaultHandler(w, r, pd.Rulebar[pd.Name])
	}
	if splitUrl[4] == "noop"{
		return template.HTML("NOOP"), false
	}
	switch splitUrl[3] {
		case "sqlite3Exec":
			return sqlite3Handler(w,r, pd.Rulebar[pd.Name], splitUrl)
		case "":
			return defaultHandler(w, r, pd.Rulebar[pd.Name])
		default:
			log.Fatal("sqlInjection Handler reached incorrectly")
			return "", false
		}
}
