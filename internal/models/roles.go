package models

type Roles struct {
	Id   int64  `gorm:"primary key"`
	Role string `gorm:"not null;uniqueIndex"`
}
