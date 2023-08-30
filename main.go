package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"math/big"
	"sort"
	"strconv"
	"time"
)

// 错误处理
func main() {
	//var wg sync.WaitGroup
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows := getScanBlock1(db)
	var id int
	var start int
	var cur int
	var end int
	for rows.Next() {
		err := rows.Scan(&id, &start, &cur, &end)
		//wg.Add(id)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("id: %d, start: %d, cur: %d, end: %d\n", id, start, cur, end)
		go run(db, id, cur, end)
	}

	select {}
	fmt.Println("All goroutines completed.")

}

func run(db *sql.DB, id int, cur int, end int) {
	for i := cur; i >= end; i-- {
		iteratorBlock(db, big.NewInt(int64(i)))
		if i%10 == 0 {
			updateScanCurBlock(db, &id, &i)
		}
	}
	updateScanIsFinished(db, &id, 1)
}

// 查询数据
func queryAccountMap(db *sql.DB) (rows *sql.Rows) {
	rows, err := db.Query("select id, address from eth_account_map")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var id int
	var address string
	for rows.Next() {
		err := rows.Scan(&id, &address)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("id: %d, address: %s\n", id, address)
	}
	return rows
}

// 插入account_map数据
func createAccountMap(db *sql.DB, addresses []string) {
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

// 遍历块
func iteratorBlock(db *sql.DB, blockNumber *big.Int) {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(header.Number.String()) // 5671744

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144
	addresses := make([]string, 0)
	fmt.Printf("------------------- blockId: %d, time: %s, count: %d------------------\n", block.Number().Uint64(), time.Unix(int64(block.Time()), 0), len(block.Transactions()))
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		fmt.Println("to address: " + tx.To().Hex()) // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
		addresses = append(addresses, tx.To().Hex())
	}

	addresses = SliceRemoveDuplicates(addresses)
	count := len(addresses)
	createAccountMap(db, addresses)

	fmt.Println("count: " + strconv.Itoa(count)) // 144
}

// 去重
func SliceRemoveDuplicates(slice []string) []string {
	sort.Strings(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice
}

func getScanBlock1(db *sql.DB) (rows *sql.Rows) {
	rows1, err := db.Query("select id, start_block, cur_block, end_block from eth_block_scan_record where status = 1 and is_finished = 0")
	if err != nil {
		panic(err.Error())
	}

	return rows1
}

// 获取扫描块范围
func getScanBlock(db *sql.DB) (id *int, startBlock *int, curBlock *int, endBlock *int) {
	rows, err := db.Query("select id, start_block, cur_block, end_block from eth_block_scan_record where status = 1 and is_finished = 0")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &startBlock, &curBlock, &endBlock)
		if err != nil {
			panic(err.Error())
		}
	}
	return id, startBlock, curBlock, endBlock
}

func updateScanCurBlock(db *sql.DB, id *int, curBlock *int) {
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	stmt, err := tx.Prepare("update eth_block_scan_record set cur_block = ? where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(curBlock, id)

	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}
}

func updateScanIsFinished(db *sql.DB, id *int, status int) {
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	stmt, err := tx.Prepare("update eth_block_scan_record set is_finished = ? where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)

	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}
}

// 根据地址获取余额
func getEthValue(address string) (ethValue *big.Float) {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25893180161173005034
	fbalance := new(big.Float)
	fbalance.SetString(fbalance.String())

	ethValue = new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue
}

func getAccountMapList(db *sql.DB) (rows *sql.Rows) {
	rows, err := db.Query("select id, address from eth_account_map")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var id int
	var address string
	for rows.Next() {
		err := rows.Scan(&id, &address)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("id: %d, address: %s\n", id, address)
	}
	return rows
}

func updateAccountMapValue(db *sql.DB, id *int, balance *float32) {
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}

	stmt, err := tx.Prepare("update eth_account_map set balance = ? where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(balance, id)

	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err.Error())
	}
}
