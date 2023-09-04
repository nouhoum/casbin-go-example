package model

import (
	"time"
)

type User struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	Email     string     `gorm:"index" json:"email"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Password  string     `json:"-"`
	RoleID    int        `json:"role_id"`
	CreatedAt *time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	Role Role `json:"-" gorm:"foreignKey:RoleID"`
}

type TodoItem struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`

	CompletedAt *time.Time `gorm:"index" json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Owner *User `json:"-" gorm:"foreignKey:OwnerID"`
}

func (item TodoItem) IsComplete() bool {
	return item.CompletedAt != nil
}

type Role struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	Slug      string     `json:"slug" gorm:"size:512;uniqueIndex:unique_index"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CasbinRule struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:512;uniqueIndex:unique_index"`
	V0    string `gorm:"size:512;uniqueIndex:unique_index"`
	V1    string `gorm:"size:512;uniqueIndex:unique_index"`
	V2    string `gorm:"size:512;uniqueIndex:unique_index"`
	V3    string `gorm:"size:512;uniqueIndex:unique_index"`
	V4    string `gorm:"size:512;uniqueIndex:unique_index"`
	V5    string `gorm:"size:512;uniqueIndex:unique_index"`
	V6    string `gorm:"size:512;uniqueIndex:unique_index"`
	V7    string `gorm:"size:512;uniqueIndex:unique_index"`
}

func (cr *CasbinRule) TableName() string {
	return "casbin_rules"
}
