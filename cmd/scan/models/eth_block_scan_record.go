package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
)

type EthBlockScanRecord struct {
	Id         int `json:"id" gorm:"primaryKey;autoIncrement"` // 编码
	StartBlock int `json:"start_block" gorm:"int;"`            // 名称
	CurBlock   int `json:"cur_block" gorm:"int;"`              // 名称
	EndBlock   int `json:"end_block" gorm:"int;"`              // 名称
	IsFinished int `json:"is_finished" gorm:"int;"`            // 是否结束
	Status     int `json:"status" gorm:"int;"`                 // 状态

	DataScope string `json:"dataScope" gorm:"-"`
}

func (e *EthBlockScanRecord) TableName() string {
	return "eth_block_scan_record"
}

func (e *EthBlockScanRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *EthBlockScanRecord) GetId() interface{} {
	return e.Id
}

func (e *EthBlockScanRecord) SetCreateBy(createBy int) {
	//e.CreateBy = createBy
}

func (e *EthBlockScanRecord) SetUpdateBy(updateBy int) {
	//e.UpdateBy = updateBy
}

func (e *EthBlockScanRecord) GetToDoList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("status = 1 and is_finished = 0").Find(list).Error
}

func (e *EthBlockScanRecord) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Find(list).Error
}

// Update 更新EthBlockScanRecord
func (e *EthBlockScanRecord) Update(tx *gorm.DB) (err error) {
	return tx.Table(e.TableName()).Updates(&e).Error
}

func (e *EthBlockScanRecord) RemoveAll(tx *gorm.DB) (err error) {
	tx.Exec("delete from " + e.TableName())
	return
}

func (e *EthBlockScanRecord) Create(tx *gorm.DB, list interface{}) (err error) {
	tx.CreateInBatches(list, 500)
	return
}
