package servegin

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	//import sqlite3 so it can be found by sql package
	_ "github.com/mattn/go-sqlite3"
)

func addSQLi(r *gin.Engine, dbSrc *os.File) {
	db := initDB(dbSrc)
	sqli := r.Group("/sqlInjection")
	sqli.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sqlInjection.gohtml", templateData("sqlInjection"))
	})

	sqli.GET("/:source/sqlite3Exec/:type", sqliHandlerFunc(db))
	sqli.POST("/:source/sqlite3Exec/:type", sqliHandlerFunc(db))
}

func initDB(dbSrc *os.File) *sql.DB {
	db, err := sql.Open("sqlite3", dbSrc.Name())
	if err != nil {
		panic(err)
	}

	return db
}

func sqliHandlerFunc(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		source := c.Param("source")
		payload := extractInput(c, source)

		var (
			res sql.Result
			err error
		)
		switch c.Param("type") {
		case "noop":
			c.String(http.StatusOK, "noop")
		case "safe":
			query := "SELECT '?' as '?'"
			res, err = db.Exec(query, payload, "test")
		case "unsafe":
			query := fmt.Sprintf("SELECT '%s' as '%s'", payload, "test")
			res, err = db.Exec(query)
		}
		if err != nil {
			c.Error(err)
			return
		}

		var n int64
		if res != nil {
			n, _ = res.RowsAffected()
		}
		c.String(http.StatusOK, fmt.Sprintf("Rows affected: %d", n))
	}
}
