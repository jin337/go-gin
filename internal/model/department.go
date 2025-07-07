package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID        uint           `json:"id" gorm:"primarykey;comment:唯一标识符"`
	Name      string         `json:"name" gorm:"comment:部门名称" `
	Type      int            `json:"type" gorm:"comment:部门类型(1机构，2部门)"`
	Pid       uint           `json:"pid" gorm:"comment:父级ID" `
	HeadUpId  uint           `json:"head_up_id" gorm:"comment:部门负责人id"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// 响应字段
type DepartmentRes struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Type      int              `json:"type"`
	Pid       uint             `json:"pid"`
	HeadUpId  uint             `json:"head_up_id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Children  *[]DepartmentRes `json:"children"`
}

// 请求字段-新增
type DepartmentReq struct {
	Name     string `json:"name"  validate:"required"`
	Type     int    `json:"type"  validate:"required,oneof=1 2"`
	Pid      *uint  `json:"pid" validate:"required"`
	HeadUpId uint   `json:"head_up_id"`
}

// 请求字段-更新
type DepartmentUpdateReq struct {
	ID       uint    `json:"id" validate:"required"`
	Name     *string `json:"name" validate:"omitempty"`
	Type     *int    `json:"type"  validate:"omitempty,oneof=1 2"`
	Pid      *uint   `json:"pid" validate:"omitempty"`
	HeadUpId *uint   `json:"head_up_id"`
}
