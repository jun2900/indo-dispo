// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAdmin = "admin"

// Admin mapped from table <admin>
type Admin struct {
	AdminID       int32  `gorm:"column:admin_id;primaryKey;autoIncrement:true" json:"admin_id"`
	RoleID        int32  `gorm:"column:role_id;not null" json:"role_id"`
	AdminName     string `gorm:"column:admin_name;not null" json:"admin_name"`
	AdminEmail    string `gorm:"column:admin_email;not null" json:"admin_email"`
	AdminPassword string `gorm:"column:admin_password;not null" json:"admin_password"`
}

// TableName Admin's table name
func (*Admin) TableName() string {
	return TableNameAdmin
}
