package models

import (
	"eth-scan/common/models"
	"gorm.io/gorm"
	"strconv"
	"time"
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

// GetListNotDeleted 获取标记删除的记录
func (e *EthAccountMap) GetListNotDeleted(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("deleted_at is null").Limit(500).Find(list).Error
}

// GetListBalanceLimit 获取未更新余额的以太地址记录
//
//	@Description:
//	@receiver e
//	@param tx
//	@param list
//	@return err
func (e *EthAccountMap) GetListBalanceLimit(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("balance is null").Limit(5).Find(list).Error
}

// Update 更新EthAccountMap
func (e *EthAccountMap) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *EthAccountMap) LogicDelete(tx *gorm.DB) (err error) {
	e.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return tx.Save(e).Error
}

func (e *EthAccountMap) Delete(tx *gorm.DB) (err error) {
	//tx.Table(e.TableName()).Where("id = ?", e.Id).Delete(e)
	tx.Exec("delete from " + e.TableName() + " where id = " + strconv.Itoa(e.Id))
	return
}

func (e *EthAccountMap) CreateInBatches(tx *gorm.DB, list interface{}) (err error) {
	tx.CreateInBatches(list, 500)
	return
}
