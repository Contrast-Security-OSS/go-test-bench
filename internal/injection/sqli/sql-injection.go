package sqli

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"

	// database import for sqlite3
	_ "github.com/mattn/go-sqlite3"
)

func sqliTemplate(w http.ResponseWriter, r *http.Request, params common.Parameters) (template.HTML, bool) {
	return "sqlInjection.gohtml", true
}

func jsonHeadersHandler(w http.ResponseWriter, r *http.Request, splitURL []string) (template.HTML, bool) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	rawCredentials := r.Header.Get("credentials")
	credBytes := []byte(rawCredentials)
	err := json.Unmarshal(credBytes, &credentials)
	if err != nil {
		log.Printf("Could not parse headers to json, err = %s", err)
	}

	var query string
	query, err = getSqliteQuery(credentials.Username, splitURL[4])
	if err != nil {
		return template.HTML(err.Error()), false
	}

	return template.HTML(query), false //change to out desired output
}

func getSqliteSafetyQuery(db *sql.DB, userInput, safety string) string {
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

func sqlite3Handler(w http.ResponseWriter, r *http.Request, splitURL []string) (template.HTML, bool) {
	// Split by source
	var query string
	switch splitURL[2] {
	case "headers-json":
		return jsonHeadersHandler(w, r, splitURL)
	case "body":
		return bodyHandler(w, r, splitURL)
	case "query":
		var err error
		userInput := common.GetUserInput(r)
		query, err = getSqliteQuery(userInput, splitURL[4])
		if err != nil {
			return template.HTML(err.Error()), false
		}
	default:
		log.Printf("Invalid source type '%s' for sql injection", splitURL[2])
	}

	return template.HTML(query), false //change to out desired output
}

func bodyHandler(w http.ResponseWriter, r *http.Request, splitURL []string) (template.HTML, bool) {
	var (
		err       error
		n         int
		ret       []byte
		userInput string
		bCnt      = 10
		bRet      = make([]byte, bCnt)
	)
	for err != io.EOF {
		n, err = r.Body.Read(bRet)
		ret = append(ret, bRet[0:n]...)
	}
	userInput, err = url.QueryUnescape(string(ret)) // POST body comes encoded but we need it in raw format
	if err != nil {
		log.Printf("Could not escape body: %s", err)
		return template.HTML(err.Error()), false
	}
	userInput = strings.TrimPrefix(userInput, "input=")

	query, err := getSqliteQuery(userInput, splitURL[4])
	if err != nil {
		return template.HTML(err.Error()), false
	}
	return template.HTML(query), false //change to out desired output
}

func getSqliteQuery(userInput, safety string) (string, error) {
	// setting up a database to execute the built query
	_ = os.Remove("tempDatabase.db")
	log.Println("Creating tempDatabase.db...")
	file, err := os.Create("tempDatabase.db")
	if err != nil {
		log.Println(err)
		return "", err
	}
	_ = file.Close()
	log.Println("tempDatabase.db created")
	sqlite3Database, _ := sql.Open("sqlite3", "./tempDatabase.db")

	defer func() {
		_ = sqlite3Database.Close()
	}()

	query := getSqliteSafetyQuery(sqlite3Database, userInput, safety)

	_ = os.Remove("tempDatabase.db")

	return query, nil
}

// Handler is the API handler for sql injection
func Handler(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 4 {
		return sqliTemplate(w, r, pd)
	}
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}
	// Split by sink
	switch splitURL[3] {
	case "sqlite3Exec":
		return sqlite3Handler(w, r, splitURL)
	case "":
		return sqliTemplate(w, r, pd)
	default:
		log.Fatal("sqlInjection Handler reached incorrectly")
		return "", false
	}
}
