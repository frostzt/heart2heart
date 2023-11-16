package main

import (
	misc_v1 "apps/hippo/controllers/v1/misc"
	"apps/hippo/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestMiscVersion(t *testing.T) {
	r := SetUpRouter()
	miscController := misc_v1.NewMiscController(utils.GetLogger())
	r.GET("/version", miscController.GetVersion)
	req, _ := http.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
