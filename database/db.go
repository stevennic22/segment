package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"segment/core"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	// Get DB Starter check
	dbStartLimit, _ := strconv.Atoi(core.CurrentEnv.RetrieveValue("db_starter"))
	dbStartAttempts := 1

	var lastErr error

	maximumDelayTime := 120 * time.Second

	// Open DB connection
	for dbStartAttempts <= dbStartLimit {
		log.Println("db Start")
		db, dbStartErr := sql.Open(
			"mysql",
			fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s",
				core.CurrentEnv.RetrieveValue("db_user"),
				core.CurrentEnv.RetrieveValue("db_pass"),
				core.CurrentEnv.RetrieveValue("db_ip"),
				core.CurrentEnv.RetrieveValue("db_port"),
				core.CurrentEnv.RetrieveValue("db_to_use"),
			),
		)

		if core.CheckError(dbStartErr, false) {
			lastErr = dbStartErr
			dbStartAttempts += 1
			time.Sleep(core.DelayDuration(dbStartAttempts, rand.Float64(), maximumDelayTime))
			continue
		}

		// Open doesn't open a connection. Validate DSN data:
		log.Println("db.Ping")
		dbPingErr := db.Ping()
		if core.CheckError(dbPingErr, false) {
			lastErr = dbPingErr
			dbStartAttempts += 1
			time.Sleep(core.DelayDuration(dbStartAttempts, rand.Float64(), maximumDelayTime))
			continue
		}

		// Connect and check the server version
		var version string
		db.QueryRow("SELECT VERSION()").Scan(&version)
		log.Printf("Connected to: %s\n", version)

		return db, nil
	}

	return nil, lastErr
}

// DBHandler will hold everything that controller needs
type DBHandler struct {
	db *sql.DB
}

// NewDBHandler returns a new DBHandler
func NewDBHandler(db *sql.DB) *DBHandler {
	return &DBHandler{
		db: db,
	}
}

// HelloWorld returns Hello, World
func (h *DBHandler) HelloWorld(c *gin.Context) {
	if err := h.db.Ping(); err != nil {
		log.Println("DB Error")
	}

	c.String(http.StatusOK, "DB ping successful.")
}

func (h *DBHandler) prepQuery(table string) (*sql.Stmt, error) {
	formattedQuery := fmt.Sprintf("Select val FROM %s where name = ?", table)

	query, err := h.db.Prepare(formattedQuery)

	if core.CheckError(err, false) {
		query.Close()
		return nil, fmt.Errorf("DBQueryPrep: Failed to prep query")
	}

	return query, nil
}

func (h *DBHandler) DBBoolFromKey(table, key string) (bool, error) {
	var val bool

	preppedQuery, preppedErr := h.prepQuery(table)

	if core.CheckError(preppedErr, false) {
		return val, fmt.Errorf("DBBoolFromKey %s: Failed to prep query", key)
	}

	defer preppedQuery.Close()

	row := preppedQuery.QueryRow(key)
	if scanErr := row.Scan(&val); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return val, fmt.Errorf("DBBoolFromKey %s: no such key", key)
		}
		return val, fmt.Errorf("DBBoolFromKey %s: %v", key, scanErr)
	}
	return val, nil
}

func (h *DBHandler) DBStringFromKey(table, key string) (string, error) {
	var val string

	preppedQuery, preppedErr := h.prepQuery(table)

	if core.CheckError(preppedErr, false) {
		return val, fmt.Errorf("DBStringFromKey %s: Failed to prep query", key)
	}

	defer preppedQuery.Close()

	row := preppedQuery.QueryRow(key)
	if err := row.Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return val, fmt.Errorf("DBStringFromKey %s: no such key", key)
		}
		return val, fmt.Errorf("DBStringFromKey %s: %v", key, err)
	}
	return val, nil
}

func (h *DBHandler) AllowedOriginCheck(origin string) bool {
	originAllowed, oAErr := h.DBBoolFromKey("AOSubs", origin)

	if core.CheckError(oAErr, false) {
		return false
	}

	if originAllowed {
		return true
	} else {
		return false
	}
}
