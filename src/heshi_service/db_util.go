package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"util"

	"github.com/fatih/color"
)

var sqlPrinter = color.New(color.FgBlue)

func traceSQL(query string, args ...interface{}) {
	if util.ShouldTrace() {
		util.Traceln(sqlPrinter.Sprintf("\t-> %s %s", query, strings.Trim(fmt.Sprint(args), "[]")))
	}
}

func dbExec(query string, args ...interface{}) (sql.Result, error) {
	traceSQL(query, args)
	r, err := db.Exec(query, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func dbExecContext(query string, args ...interface{}) (sql.Result, error) {
	traceSQL(query, args)
	r, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func dbQuery(query string, args ...interface{}) (*sql.Rows, error) {
	traceSQL(query, args)
	r, err := db.Query(query, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func dbQueryContext(query string, args ...interface{}) (*sql.Rows, error) {
	traceSQL(query, args)
	r, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return r, err
	}
	return r, nil
}

func dbQueryRow(query string, args ...interface{}) *sql.Row {
	traceSQL(query, args)
	return db.QueryRow(query, args...)
}

func dbQueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	traceSQL(query, args)
	return db.QueryRowContext(ctx, query, args...)
}

func dbTransact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
