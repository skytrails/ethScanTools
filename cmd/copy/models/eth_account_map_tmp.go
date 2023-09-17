package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
)

type EthAccountMapTmp struct {
	Id      int     `json:"id" gorm:"primaryKey;autoIncrement"` // 唯一编码
	Address string  `json:"address" gorm:"size:255;"`           // 地址
	Balance float64 `json:"balance" gorm:"decimal(24,6);"`
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (*EthAccountMapTmp) TableName() string {
	return "eth_account_map_tmp"
}

func (e *EthAccountMapTmp) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EthAccountMapTmp) GetId() interface{} {
	return e.Id
}

func (e *EthAccountMapTmp) SetCreateBy(createBy int) {
	//e.CreateBy = createBy
}

func (e *EthAccountMapTmp) SetUpdateBy(updateBy int) {
	//e.UpdateBy = updateBy
}

func (e *EthAccountMapTmp) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("deleted is null").Limit(50).Find(list).Error
}

// GetListBalanceLimit 获取未更新余额的以太地址记录
func (e *EthAccountMapTmp) GetListBalanceLimit(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("balance is null").Limit(5).Find(list).Error
}

// Update 更新EthAccountMap
func (e *EthAccountMapTmp) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *EthAccountMapTmp) RemoveAll(tx *gorm.DB) (err error) {
	tx.Exec("delete from " + e.TableName())
	return
}

func (e *EthAccountMapTmp) CreateInBatches(tx *gorm.DB, list interface{}) (err error) {
	tx.CreateInBatches(list, 500)
	return
}

func (e *EthAccountMapTmp) Create(tx *gorm.DB, value interface{}) (err error) {
	tx.Table(e.TableName()).Create(value)
	return
}
