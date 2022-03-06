package baseHeila

import (
	. "bitbucket.org/HeilaSystems/persistentstorage"
	"gorm.io/gorm"
	"time"
)

var UpdStampGetter UpdateStampGetter

type BaseHeilaEntity struct {
	RecNo       uint           `gorm:"column:recNo;primary_key" json:"recNo"`
	CreatedAt   time.Time      `gorm:"column:createDate" json:"createDate"`
	DeletedAt   gorm.DeletedAt `gorm:"column:cancelDate" json:"cancelDate"`
	UpdateStamp int64          `gorm:"column:updateStamp" json:"updateStamp"`
}

func (be *BaseHeilaEntity) setUpdateStamp(tx *gorm.DB) error {
	var up int64  = -1
	err := tx.Table(tx.Statement.Table).Select("IFNULL(MAX(updateStamp), 0) AS updateStamp").Pluck("updateStamp", &up).Error
	if err != nil{
		return err
	}
	be.UpdateStamp = up+1
	return nil
}

func (m *BaseHeilaEntity) BeforeUpdate(tx *gorm.DB) error {
	return m.setUpdateStamp(tx)
}

func (m *BaseHeilaEntity) BeforeCreate(tx *gorm.DB) error {
	return m.setUpdateStamp(tx)
}

func (m *BaseHeilaEntity) BeforeDelete(tx *gorm.DB) error {
	return m.setUpdateStamp(tx)
}
