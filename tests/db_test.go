package tests

import (
	"database/sql"
	"log"
	"testing"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"url-shortener/db"
)

var (
	urlTest = "https://www.test.com/12222212/33423423423"
	shortTest = "https://tic.tac/7PQTO2E"
)


func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error '%s' occured when opening a stub database connection", err)
	}
	return db, mock
}

func TestFindByPkey(t *testing.T) {
	_db, mock := NewMock()
	client := &db.DBClient{_db}
	defer client.Close()

	rows := sqlmock.NewRows([]string{"short"}).AddRow(shortTest)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short FROM url_to_short WHERE real_url = $1`)).WithArgs(urlTest).WillReturnRows(rows)
	short, err := client.FindByPkey("short", "url_to_short", "real_url", urlTest)
	
	assert.Equal(t, short, shortTest)
	assert.NotNil(t, short)
	assert.NoError(t, err)
}

func TestFindByPkeyEmpty(t *testing.T) {
	_db, mock := NewMock()
	client := &db.DBClient{_db}
	defer client.Close()

	rows := sqlmock.NewRows([]string{"short"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short FROM url_to_short WHERE real_url = $1`)).WithArgs(urlTest).WillReturnRows(rows)
	short, err := client.FindByPkey("short", "url_to_short", "real_url", urlTest)
	
	assert.Empty(t, short)
	assert.NoError(t, err)
}

func TestFindByPkeyError(t *testing.T) {
	_db, mock := NewMock()
	client := &db.DBClient{_db}
	defer client.Close()

	rows := sqlmock.NewRows([]string{"short"}).AddRow(shortTest)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short FROM url_to_short WHERE real_url = $1`)).WithArgs(urlTest).WillReturnRows(rows)
	_, err := client.FindByPkey("short", "url_to_short", "real_url", "url")
	
	assert.Error(t, err)
}

func TestInsert(t *testing.T) {
	_db, mock := NewMock()
	client := &db.DBClient{_db}
	defer client.Close()

	query := regexp.QuoteMeta(`INSERT INTO url_to_short (real_url, short) VALUES ($1, $2)`)
	mock.ExpectExec(query).WithArgs(urlTest, shortTest).WillReturnResult(sqlmock.NewResult(0, 1))
	err := client.Insert("url_to_short", [2]string{"real_url", "short"}, [2]string{urlTest, shortTest})
	assert.NoError(t, err)
}

func TestInsertError(t *testing.T) {
	_db, mock := NewMock()
	client := &db.DBClient{_db}
	defer client.Close()

	query := regexp.QuoteMeta(`INSERT INTO url_to_short (real_url, short) VALUES ($1, $2)`)
	mock.ExpectExec(query).WithArgs(urlTest, shortTest).WillReturnResult(sqlmock.NewResult(0, 1))
	err := client.Insert("url_to_short", [2]string{"real_url", "short"}, [2]string{"url", shortTest})
	assert.Error(t, err)
}
