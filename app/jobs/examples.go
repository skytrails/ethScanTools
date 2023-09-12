package jobs

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	models2 "go-admin/app/jobs/models"
	"go-admin/common/models"
	"math/rand"
	"time"
)

// InitJob
// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func InitJob() {
	jobList = map[string]JobExec{
		"ExamplesOne":      ExamplesOne{},      // 根据Eth生成地址
		"GenerateByRandom": GenerateByRandom{}, // 根据Random生成地址
		"CheckEthValue":    CheckEthValue{},    // 校验值
		// ...
	}
}

/*
 * 通过GEth生成key
 */
func generateByEth() error {
	var data = []models2.EthAccountAddressMap{}
	for i := 0; i < 100; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Error(err)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)
		priv := hexutil.Encode(privateKeyBytes)[2:]
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		data = append(data, models2.EthAccountAddressMap{
			PrivKey: priv,
			Address: address,
		})
	}

	// 插入数据库
	dbList := sdk.Runtime.GetDb()
	ethAccountAddressMap := models2.EthAccountAddressMap{}
	for _, d := range dbList {
		tx := d.Begin()
		err := ethAccountAddressMap.Create(tx, data)
		if err != nil {
			fmt.Println(time.Now().Format(timeFormat), " [ERROR] generateByEth.JobCore create error", err)
			d.Rollback()
		} else {
			tx.Commit()
		}
	}

	return nil
}

func generate32() (str string) {
	var b [16]byte
	n, err := rand.Read(b[:])
	if n <= 0 {
		log.Fatal(err)
		return
	}
	buf := make([]byte, 32)
	hex.Encode(buf, b[:])
	return string(buf)
}

func generate16() (str string) {
	var b [16]byte
	n, err := rand.Read(b[:])
	if n <= 0 {
		log.Fatal(err)
		return
	}
	buf := make([]byte, 16)
	hex.Encode(buf, b[:])
	return string(buf)
}

/**
 * 通过Random生成私钥
 */
func generateByRandom() error {
	var data = []models2.EthAccountAddressMap{}
	for i := 0; i < 1000; i++ {
		buf0 := generate32()
		buf1 := generate32()
		privateKey, err := crypto.HexToECDSA(buf0 + buf1)
		if err != nil {
			log.Fatal(err)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)
		priv := hexutil.Encode(privateKeyBytes)[2:]
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		data = append(data, models2.EthAccountAddressMap{
			PrivKey: priv,
			Address: address,
		})
	}

	// 插入数据库
	dbList := sdk.Runtime.GetDb()
	ethAccountAddressMap := models2.EthAccountAddressMap{}
	for _, d := range dbList {
		tx := d.Begin()
		err := ethAccountAddressMap.Create(tx, data)
		if err != nil {
			log.Debug(time.Now().Format(timeFormat), " [ERROR] generateByRandom.JobCore init create error", err)
			tx.Rollback()
		} else {
			//log.Debug(jobList)
			tx.Commit()
		}
	}

	return nil
}

func checkEthValue() {
	dbList := sdk.Runtime.GetDb()
	ethAccountAddressMap := models2.EthAccountAddressMap{}
	ethAccountAvaild := models2.EthAccountAvaild{}
	jobList := make([]models2.EthAccountAddressMap, 0)
	for _, d := range dbList {

		log.Debug("--------------------")
		tx := d.Debug().Begin()
		err := ethAccountAddressMap.Check(tx, &jobList)
		if err != nil {
			log.Debug(time.Now().Format(timeFormat), " [ERROR] checkEthValue.JobCore check error", err)
			return
		} else {
			for _, data := range jobList {
				// insert
				err := ethAccountAvaild.Create(tx, &models2.EthAccountAvaild{
					PrivKey:   data.PrivKey,
					Address:   data.Address,
					ModelTime: models.ModelTime{},
					DataScope: "",
				})
				if err != nil {
					break
				}
			}
			err := ethAccountAddressMap.RemoveAll(tx)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		tx.Debug().Commit()
	}

}

// ExamplesOne
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct {
}

func (t ExamplesOne) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesOne exec success"
	// TODO: 这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			log.Debug(str, arg.(string))
			err := generateByEth()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}

// GenerateByRandom
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type GenerateByRandom struct {
}

func (t GenerateByRandom) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore GenerateByRandom exec success"
	// TODO: 这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			log.Debug(str, arg.(string))
			err := generateByRandom()
			if err != nil {
				return err
			}

		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}
	return nil
}

// CheckEthValue
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type CheckEthValue struct {
}

func (t CheckEthValue) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore CheckEthValue exec success"
	// TODO: 这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			log.Debug(str, arg.(string))
			checkEthValue()

		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}
