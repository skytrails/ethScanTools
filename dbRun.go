package main

import (
	"database/sql"
	"fmt"
)

func dbRun() {
	dbLocal, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		panic(err.Error())
	}
	defer dbLocal.Close()

	dbRemote, err := sql.Open("mysql", "root:Iwantfly@2022@tcp(120.27.135.165:3306)/go_db")
	if err != nil {
		panic(err.Error())
	}
	defer dbRemote.Close()
	rows := getDataLocal(dbLocal)

	var id int
	var address string
	for rows.Next() {
		err := rows.Scan(&id, &address)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("id: %d, address: %s\n", id, address)
	}
}

func getDataLocal(db *sql.DB) (rows *sql.Rows) {
	rows, err := db.Query("select id, address from eth_account_map")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	return rows
}

func insertDataRemote(db *sql.DB, addresses []string) {
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	for _, address := range addresses {
		stmt, err := tx.Prepare("insert into eth_account_map (address, created_time) values (?, now())")
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(address)

		if err != nil {
			tx.Rollback()
			panic(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	fmt.Println("create committed!")

}
