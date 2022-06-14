
package main
 
import (
	"testing"
	"net/http"
	"fmt"
	"net/http/httptest"
)

// Make a get request and verify we get an error
func TestUrlGetError(t *testing.T) {

	req := httptest.NewRequest("GET", "http://localhost:10000/get_url", nil)
	w := httptest.NewRecorder()

	//TODO continue implementation
}

// Make a valid post request and make sure we get 200ok
// Make anohter request of the same url and make sure we get an error
func TestUrlPost(t *testing.T) {
}

// Make a valid get request for an existing url and make sure we get an error 200ok
func TestUrlGet(t *testing.T) {
}