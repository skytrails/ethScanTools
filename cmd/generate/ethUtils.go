package generate

import (
	"crypto/ecdsa"
	"encoding/hex"
	"eth-scan/cmd/generate/models"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/go-admin-team/go-admin-core/logger"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"strings"
)

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

/**
 * 通过Random生成私钥
 */
func generateByRandom(db *gorm.DB) error {
	var data []models.EthAccountAddressMap
	var addresses []string
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

		data = append(data, models.EthAccountAddressMap{
			PrivKey: priv,
			Address: address,
		})

		addresses = append(addresses, address)
	}

	ethAccountMap := models.EthAccountMap{}
	ethAccountMaps := make([]models.EthAccountMap, 0)
	err := ethAccountMap.GetListExist(db, addresses, &ethAccountMaps)
	if err != nil {
		return err
	}
	if len(ethAccountMaps) > 0 {
		for _, eth := range ethAccountMaps {
			for _, d := range data {
				if strings.ToLower(eth.Address) == strings.ToLower(d.Address) {
					log.Debug("private_key", d.PrivKey, " Address:", d.Address)
					ethAccountAvalid := models.EthAccountAvaild{PrivKey: d.PrivKey, Address: d.Address}
					err := ethAccountAvalid.Create(db)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}

/*
 * 通过GEth生成key
 */
func generateByEth(db *gorm.DB) error {
	var data []models.EthAccountAddressMap
	var addresses []string
	for i := 0; i < 1000; i++ {
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

		data = append(data, models.EthAccountAddressMap{
			PrivKey: priv,
			Address: address,
		})
		addresses = append(addresses, address)
	}
	ethAccountMap := models.EthAccountMap{}
	ethAccountMaps := make([]models.EthAccountMap, 0)
	err := ethAccountMap.GetListExist(db, addresses, &ethAccountMaps)
	if err != nil {
		return err
	}
	if len(ethAccountMaps) > 0 {
		for _, eth := range ethAccountMaps {
			for _, d := range data {
				fmt.Println("d.Address:", d.Address, " eth.Address:", eth.Address)
				if strings.ToLower(eth.Address) == strings.ToLower(d.Address) {
					log.Debug("private_key", d.PrivKey, " Address:", d.Address)
					ethAccountAvalid := models.EthAccountAvaild{PrivKey: d.PrivKey, Address: d.Address}
					err := ethAccountAvalid.Create(db)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}
