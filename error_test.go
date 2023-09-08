package err

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var BadRequestErr = fmt.Errorf("bad request error")

func TestErrToStatusCode(t *testing.T) {
	r := gin.Default()
	r.Use(Error(NewErrMap(BadRequestErr).StatusCode(http.StatusBadRequest)))
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(BadRequestErr)
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	if recorder.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("invalid the status code %+v", recorder.Result().StatusCode)
	}
}

// TestWrappedErrToStatusCode ensures that the middleware also works with wrapped errors.
func TestWrappedErrToStatusCode(t *testing.T) {
	r := gin.Default()
	r.Use(Error(NewErrMap(BadRequestErr).StatusCode(http.StatusBadRequest)))
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("%w: this is a wrapped error", BadRequestErr))
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	if recorder.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("invalid the status code %+v", recorder.Result().StatusCode)
	}
}

func TestErrToResponse(t *testing.T) {
	r := gin.Default()
	r.Use(Error(
		NewErrMap(BadRequestErr).Response(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": BadRequestErr.Error()})
		})))
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(BadRequestErr)
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	if recorder.Result().StatusCode != http.StatusNotFound {
		t.Errorf("invalid the status code %+v", recorder.Result().StatusCode)
		return
	}

	buf, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("failed to read body err: %+v", err)
		return
	}

	var rsp = struct {
		Error string `json:"error"`
	}{}
	err = json.Unmarshal(buf, &rsp)
	if err != nil {
		t.Errorf("failed to unmarshal body err: %+v", err)
		return
	}

	if rsp.Error != BadRequestErr.Error() {
		t.Errorf("invalid the response body %+v", rsp.Error)
	}
}

func TestErrorMap_matchError(t *testing.T) {
	t.Run("single error", func(t *testing.T) {
		em := NewErrMap(BadRequestErr)
		err := BadRequestErr
		if em.matchError(err) == false {
			t.Errorf("error should match")
		}
	})

	t.Run("wrapped error", func(t *testing.T) {
		em := NewErrMap(BadRequestErr)
		err := fmt.Errorf("%w: Details about this error", BadRequestErr)
		if em.matchError(err) == false {
			t.Errorf("error should match")
		}
	})
}
