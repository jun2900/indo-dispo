package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type AdminService interface {
	GetAdminByEmail(email string) (result *model.Admin, err error)
}

func NewAdminService(mysqlConnection *gorm.DB) AdminService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) GetAdminByEmail(email string) (result *model.Admin, err error) {
	if err = r.mysql.Model(&model.Admin{}).Where("admin_email = ?", email).First(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
