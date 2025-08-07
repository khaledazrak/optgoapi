package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/helloworld", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.JSON(200, gin.H{"message": "Hello World"})
		} else {
			c.JSON(200, gin.H{"message": "Hello " + FormatName(name)})
		}
	})
	r.GET("/versionz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"project":  "go-api-project",
			"git_hash": GetGitHash(),
		})
	})
	return r
}

func TestHelloWorld(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/helloworld", nil)
	setupRouter().ServeHTTP(w, req)

	if w.Code != 200 || w.Body.String() != "{\"message\":\"Hello World\"}" {
		t.Errorf("Expected Hello World, got %s", w.Body.String())
	}
}

func TestHelloName(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/helloworld?name=hello_world", nil)
	setupRouter().ServeHTTP(w, req)

	if w.Code != 200 || w.Body.String() != "{\"message\":\"Hello Hello World\"}" {
		t.Errorf("Expected formatted name, got %s", w.Body.String())
	}
}

func TestVersion(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/versionz", nil)
	setupRouter().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

