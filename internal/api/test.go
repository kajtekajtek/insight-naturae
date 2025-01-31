// internal/api/test.go - common variables and functions for api testing
package api

import (
	"os"
	"database/sql"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"bytes"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/internal/middleware"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/internal/wsutils"
)

const (
	DBPath = "test.db"
	JWTSecret = "testjwtscretkey"
	username = "testuser"
	password = "testpassword"
	sensorID = "testsensor"
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

func SetupTestRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	cm := wsutils.NewWSClientManager()

	r.POST("/register", RegisterHandler(db))
	r.POST("/login", LoginHandler(db, []byte(JWTSecret)))

	protected := r.Group("/user", middleware.AuthMiddleware([]byte(JWTSecret)))
	{
		protected.POST("/sensors", SubscribeSensorHandler(db, cm))
		protected.GET("/sensors", GetSensorsHandler(db))
		protected.DELETE("/sensors/:id", UnsubscribeSensorHandler(db))
	}

	return r
}

// register, login and return the token string
func SetupUser(t *testing.T, db *sql.DB, r *gin.Engine) string {
	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}	
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// log in
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	tokenString := strings.Split(w.Body.String(), "\"")[3]
	assert.Len(t, tokenString, 135)

	// return the token string
	return tokenString
}