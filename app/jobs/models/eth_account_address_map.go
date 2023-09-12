package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
)

type EthAccountAddressMap struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 编码
	PrivKey string `json:"privKey" gorm:"size:255;"`           // 名称
	Address string `json:"address" gorm:"size:255;"`           // 名称
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (*EthAccountAddressMap) TableName() string {
	return "eth_account_address_map"
}

func (e *EthAccountAddressMap) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EthAccountAddressMap) GetId() interface{} {
	return e.Id
}

func (e *EthAccountAddressMap) SetCreateBy(createBy int) {
	//e.CreateBy = createBy
}

func (e *EthAccountAddressMap) SetUpdateBy(updateBy int) {
	//e.UpdateBy = updateBy
}

func (e *EthAccountAddressMap) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Find(list).Error
}

// Update 更新EthAccountAddressMap
func (e *EthAccountAddressMap) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *EthAccountAddressMap) RemoveAll(tx *gorm.DB) (err error) {
	tx.Exec("delete from " + e.TableName())
	return
}

func (e *EthAccountAddressMap) Create(tx *gorm.DB, list interface{}) (err error) {
	tx.CreateInBatches(list, 500)
	return
}

func (e *EthAccountAddressMap) Check(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Select("eth_account_address_map.*").Joins("inner join eth_account_map on eth_account_map.address = eth_account_address_map.address").Find(list).Error
}
