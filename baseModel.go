package persistentstorage

import "time"

type RecNo uint

type BaseModel struct {
	RecNo       RecNo      `gorm:"column:recNo;primary_key" json:"recNo"`
	CreatedAt   time.Time  `gorm:"column:createDate" json:"createDate"`
	DeletedAt   *time.Time `gorm:"column:cancelDate" json:"cancelDate"`
	UpdateStamp int64      `gorm:"column:updateStamp" json:"updateStamp"`
}
