
package main
 
import (
	"testing"
	"net/http"
	"fmt"
	"net/http/httptest"
	"log"
	"bytes"
	"encoding/json"
)

type Response struct {
	Message string  `json:"message"`
	Status  string `json:"Status"`
}

//Setting up a GET test
func setupTestGet(tb testing.TB, param string, value string) (func(tb testing.TB), *httptest.ResponseRecorder) {
	log.Println("Setup test for GET endpoint")
	
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		tb.Fatal(err)
	}

	queryParam := req.URL.Query()
	queryParam.Add(param, fmt.Sprint(value))
	req.URL.RawQuery = queryParam.Encode()
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(GetUrl)
	h.ServeHTTP(rr, req)

	return func(tb testing.TB) {
		log.Println("Teardown test for GET endpoint")
	}, rr
}

//Setting up a POST test
func setupTestPost(tb testing.TB, value string) (func(tb testing.TB), *httptest.ResponseRecorder) {
	log.Println("Setup test for POST endpoint")

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	if err != nil {
		tb.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(ShortenUrl)
	h.ServeHTTP(rr, req)

	return func(tb testing.TB) {
		log.Println("Teardown test for POST endpoint")
	}, rr
}

/* 
	This test invokes a first call to the GET endpoint with URL that doesn't exist
	Expected result is status code 400
*/
func TestUrlGet(t *testing.T) {

	tearDownTest, rr := setupTestGet(t, "url", "NOTEXISTS")
	defer tearDownTest(t)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Get endpoint returned wrong status code: got %v but expected %v",
			status, http.StatusBadRequest)
	}
}

/* 
	This test invokes a first call to the POST endpoint with a new URL, expected result is status code 200
	After that, it calls to the GET endpoint with the short URL, expected result is 200
*/
func TestUrlPostGet(t *testing.T) {

	tearDownTest, rr := setupTestPost(t, `{"Url":"https://example.com"}`)
	defer tearDownTest(t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Post endpoint returned wrong status code: got %v but expected %v",
			status, http.StatusOK)
	}
	actual := Response{}
	json.Unmarshal([]byte(rr.Body.String()), &actual)

	tearDownTest, rr = setupTestGet(t, "url", actual.Message)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Get endpoint returned wrong status code: got %v but expected %v",
			status, http.StatusOK)
	}
}

/* 
	This test invokes a call to the POST endpoint without Url field (Url1 by mistake)
	Expected result is 400
*/
func TestUrlPostNoUrl(t *testing.T) {

	tearDownTest, rr := setupTestPost(t, `{"Url1":"https://example.com"}`)
	defer tearDownTest(t)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Post endpoint returned wrong status code: got %v but expected %v",
			status, http.StatusOK)
	}

}

/* 
	This test invokes a call to the GET endpoint without url
	Expected result is 400
*/
func TestUrlGetNoUrl(t *testing.T) {

	tearDownTest, rr := setupTestGet(t, "x", "NOTEXISTS")
	defer tearDownTest(t)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Get endpoint returned wrong status code: got %v but expected %v",
			status, http.StatusBadRequest)
	}
}
