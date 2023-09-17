package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
)

type EthAccountMap struct {
	Id      int     `json:"id" gorm:"primaryKey;autoIncrement"` // 唯一编码
	Address string  `json:"address" gorm:"size:255;"`           // 地址
	Balance float64 `json:"balance" gorm:"decimal(24,6);"`
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (*EthAccountMap) TableName() string {
	return "eth_account_map"
}

func (e *EthAccountMap) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EthAccountMap) GetId() interface{} {
	return e.Id
}

func (e *EthAccountMap) SetCreateBy(createBy int) {
	//e.CreateBy = createBy
}

func (e *EthAccountMap) SetUpdateBy(updateBy int) {
	//e.UpdateBy = updateBy
}

func (e *EthAccountMap) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Find(list).Error
}

func (e *EthAccountMap) GetListExist(tx *gorm.DB, addresses interface{}, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("address in (?)", addresses).Find(list).Error
}

// GetListBalanceLimit 获取未更新余额的以太地址记录
func (e *EthAccountMap) GetListBalanceLimit(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("balance is null").Limit(5).Find(list).Error
}

// Update 更新EthAccountMap
func (e *EthAccountMap) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *EthAccountMap) RemoveAll(tx *gorm.DB) (err error) {
	tx.Exec("delete from " + e.TableName())
	return
}

func (e *EthAccountMap) CreateInBatches(tx *gorm.DB, list interface{}) (err error) {
	tx.CreateInBatches(list, 500)
	return
}
