package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
)

type EthAccountAvaild struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 编码
	PrivKey string `json:"privKey" gorm:"size:255;"`           // 名称
	Address string `json:"address" gorm:"size:255;"`           // 名称
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (*EthAccountAvaild) TableName() string {
	return "eth_account_avalid"
}

func (e *EthAccountAvaild) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EthAccountAvaild) GetId() interface{} {
	return e.Id
}

func (e *EthAccountAvaild) SetCreateBy(createBy int) {
	//e.CreateBy = createBy
}

func (e *EthAccountAvaild) SetUpdateBy(updateBy int) {
	//e.UpdateBy = updateBy
}

func (e *EthAccountAvaild) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Find(list).Error
}

// Update 更新EthAccountAddressMap
func (e *EthAccountAvaild) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *EthAccountAvaild) RemoveAllEntryID(tx *gorm.DB) (update EthAccountAddressMap, err error) {
	if err = tx.Table(e.TableName()).Where("entry_id > ?", 0).Update("entry_id", 0).Error; err != nil {
		return
	}
	return
}

//func (e *EthAccountAvaild) Create(tx *gorm.DB, value interface{}) (err error) {
//	tx.Table(e.TableName()).Create(value)
//	return
//}

func (e *EthAccountAvaild) Create(tx *gorm.DB) (err error) {
	tx.Table(e.TableName()).Create(e)
	return
}

func (e *EthAccountAvaild) Check(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Select("eth_account_address_map.*").Joins("right join eth_account_map on eth_account_map.address = eth_account_address_map.address").Find(list).Error
}
