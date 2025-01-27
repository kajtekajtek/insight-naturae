// internal/api/test.go - common variables and functions for api testing
package api

import (
	"os"
	"database/sql"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
)

const (
	DBPath = "test.db"
	JWTSecret = "testjwtscretkey"
	username = "testuser"
	password = "testpassword"
)

func SetupTestDB(t *testing.T) *sql.DB {
	// delete the test database if it exists
	if _, err := os.Stat(DBPath); err == nil {
		err := os.Remove(DBPath)
		assert.NoError(t, err)
	}

	// create the test database
	db, err := dbutils.CreateDatabase(DBPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	return db
}

func TearDownTestDB(db *sql.DB) {
	db.Close()
	os.Remove(DBPath)
}

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/register", RegisterHandler(db))
	r.POST("/login", LoginHandler(db, []byte(JWTSecret)))

	protected := r.Group("/user", middleware.AuthMiddleware([]byte(JWTSecret)))
	{
		protected.POST("/sensors", AddSensorHandler(db))
		protected.GET("/sensors", GetSensorsHandler(db))
		protected.DELETE("/sensors/:id", RemoveSensorHandler(db))
	}

	return r
}