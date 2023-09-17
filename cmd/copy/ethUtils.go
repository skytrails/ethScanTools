package copy

import (
	"context"
	"eth-scan/cmd/copy/models"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

func process(db *gorm.DB) {
	ethAccountMap := models.EthAccountMap{}
	var data []models.EthAccountMap
	var i int
	for {
		start := time.Now().UnixMilli()
		i++
		err := ethAccountMap.GetListNotDeleted(db, &data)
		if err != nil {
			log.Fatal("GetListNotDeleted is failed!")
			return
		}
		tx := db.Begin()
		fmt.Printf(" 批次: %d, 处理量: %d\n", i, len(data))
		for _, d := range data {
			create(tx, d)
		}
		tx.Commit()
		end := time.Now().UnixMilli()
		fmt.Printf(" 处理量: %d, 总处理量: %d, 消耗时间: %d\n", len(data), len(data)*i, end-start)
	}

}

func create(db *gorm.DB, accountMap models.EthAccountMap) {
	tmp := models.EthAccountMapTmp{}
	tmp.Address = strings.ToLower(accountMap.Address)
	err := tmp.Create(db, &models.EthAccountMapTmp{Address: strings.ToLower(accountMap.Address)})
	if err != nil {
		log.Println(err)
		err := accountMap.Delete(db)
		if err != nil {
			return
		}
		return
	}
	err = accountMap.Delete(db)
}

// 更新地址余额
func scanBalance(db *gorm.DB) {
	for {
		ethAccountMap := models.EthAccountMap{}
		ethAccountMaps := make([]models.EthAccountMap, 0)
		err := ethAccountMap.GetListBalanceLimit(db, &ethAccountMaps)
		if err != nil {
			return
		}
		fmt.Println(ethAccountMaps)
		if len(ethAccountMaps) == 0 {
			fmt.Println("All Eth Balance updated!")
			return
		}

		for _, row := range ethAccountMaps {
			ethValue := getEthValue(row.Address)
			if ethValue != nil {
				row.Balance, _ = ethValue.Float64()
			} else {
				row.Balance = -1.0
			}
			err := row.Update(db, row.Id)
			if err != nil {
				return
			}
		}
	}
}

// 遍历块
func iteratorBlock(db *gorm.DB, blockNumber *big.Int) {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		fmt.Println(err)
		return
	}

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("block.Number = ", block.Number().Uint64())         // 5671744
	fmt.Println("block.Time = ", block.Time())                      // 1527211625
	fmt.Println("block.Transactions = ", len(block.Transactions())) // 144
	count := len(block.Transactions())
	if count == 0 {
		fmt.Println("block.Transactions 数量为 0, 退出")
		return
	}
	addresses := make([]string, 0)
	fmt.Printf("------------------- blockId: %d, time: %s, count: %d------------------\n", block.Number().Uint64(), time.Unix(int64(block.Time()), 0), count)
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		addresses = append(addresses, tx.To().Hex())
	}

	addresses = SliceRemoveDuplicates(addresses)
	ethAccountMaps := make([]models.EthAccountMap, 0)
	for _, address := range addresses {
		ethAccountMaps = append(ethAccountMaps, models.EthAccountMap{Address: address})
	}
	count = len(addresses)
	ethAccountMap := models.EthAccountMap{}
	err = ethAccountMap.CreateInBatches(db, ethAccountMaps)
	if err != nil {
		return
	}

	fmt.Println("count: " + strconv.Itoa(count)) // 144
}

// SliceRemoveDuplicates 去重
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

func updateScanIsFinished(db *gorm.DB, id *int, status int) {
	ethBlockScanRecord := models.EthBlockScanRecord{Id: *id, Status: 0, IsFinished: status}
	err := ethBlockScanRecord.Update(db)
	if err != nil {
		return
	}
}

// 根据地址获取余额
func getEthValue(address string) (ethValue *big.Float) {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress(address)
	fmt.Println("address:", address, " account:", account)

	blockNumber := big.NewInt(18115209)
	bytecode, err := client.CodeAt(context.Background(), account, blockNumber) // nil is latest block
	if err != nil {
		return nil
	}

	isContract := len(bytecode) > 0

	fmt.Printf("is contract: %v\n", isContract)
	if isContract == true {
		fmt.Println("这是个合约地址")
		return
	}
	balance, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())

	ethValue = new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue
}
