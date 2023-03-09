package services

import (
	"anara/model"

	"gorm.io/gorm"
)

type AttachmentService interface {
	GetAttachmentByBillId(billId int32) (results []model.Attachment, totalRows int64, err error)
	CreateAttachments(record []model.Attachment) (result []model.Attachment, RowsAffected int64, err error)
	DeleteAttachmentByBillId(billId int32) (err error)
}

func NewAttachmentService(mysqlConnection *gorm.DB) AttachmentService {
	return &mysqlDBRepository{
		mysql: mysqlConnection,
	}
}

func (r *mysqlDBRepository) CreateAttachments(record []model.Attachment) (result []model.Attachment, RowsAffected int64, err error) {
	db := r.mysql.Save(&record)
	if err = db.Error; err != nil {
		return nil, -1, err
	}

	return record, db.RowsAffected, nil
}

func (r *mysqlDBRepository) GetAttachmentByBillId(billId int32) (results []model.Attachment, totalRows int64, err error) {
	if err = r.mysql.Model(&model.Attachment{}).Where("bill_id = ?", billId).Count(&totalRows).Find(&results).Error; err != nil {
		return nil, -1, err
	}

	return results, totalRows, nil
}

func (r *mysqlDBRepository) DeleteAttachmentByBillId(billId int32) (err error) {
	if err := r.mysql.Where("bill_id = ?", billId).Delete(&model.Attachment{}).Error; err != nil {
		return err
	}
	return nil
}
