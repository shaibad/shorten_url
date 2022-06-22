package db

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

type DB interface {
	Close()
	FindByPkey(string, string, string, string) (string, error)
	Insert(string, [2]string, [2]string) error
}

type DBClient struct {
	Db *sql.DB
}

func NewDB(kind, info string) (DB, error) {
	db, err := sql.Open(kind, info)
	if err != nil {
		return nil, err
	}

	return &DBClient{db}, nil
}

func (r *DBClient) FindByPkey(valueToSelect, tableName, pkey, value string) (string, error) {
	var result string
	s := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, valueToSelect, tableName, pkey)
	err := r.Db.QueryRow(s, value).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return "", err
	}
	return result, nil
}

func (r *DBClient) Insert(tableName string, valueKeys, values [2]string) error {
	statement := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2)`, tableName, valueKeys[0], valueKeys[1])
	_, err := r.Db.Exec(statement, values[0], values[1])
	if err != nil {
		return err
	}
	return nil
}

func (r *DBClient) Close() {
	r.Db.Close()
}
