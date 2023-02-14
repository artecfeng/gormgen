package admin

import "time"

// Admin 管理员表
//go:generate gormgen -structs Admin -input .
type Admin struct {
	Id          int32     `gorm:"id"`           // 主键
	Username    string    `gorm:"username"`     // 用户名
	Password    string    `gorm:"password"`     // 密码
	Nickname    string    `gorm:"nickname"`     // 昵称
	Mobile      string    `gorm:"mobile"`       // 手机号
	IsUsed      int32     `gorm:"is_used"`      // 是否启用 1:是  -1:否
	IsDeleted   int32     `gorm:"is_deleted"`   // 是否删除 1:是  -1:否
	CreatedAt   time.Time `gorm:"created_at"`   // 创建时间
	CreatedUser string    `gorm:"created_user"` // 创建人
	UpdatedAt   time.Time `gorm:"updated_at"`   // 更新时间
	UpdatedUser string    `gorm:"updated_user"` // 更新人
}
