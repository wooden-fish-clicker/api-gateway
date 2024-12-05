package db

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	dbClient *gorm.DB

	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	ModifiedAt time.Time `json:"-" gorm:"autoUpdateTime;column:modified_at"`
}

func (m *BaseModel) FormatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

func (m *BaseModel) Create(createData interface{}) error {
	if err := m.dbClient.Create(createData).Error; err != nil {
		return fmt.Errorf("建立資料發生錯誤: %w", err)
	}
	return nil
}

func (m *BaseModel) Update(id uint, updateData interface{}) error {
	if err := m.dbClient.Where("id = ?", id).Updates(updateData).Error; err != nil {
		return fmt.Errorf("更新資料發生錯誤 : %v", err)
	}

	return nil
}

func (m *BaseModel) Delete(id uint, model interface{}) error {

	if err := m.dbClient.Where("id = ?", id).Delete(model).Error; err != nil {
		return fmt.Errorf("刪除資料發生錯誤 : %v", err)
	}

	return nil

}

func (m *BaseModel) GetDetailById(id uint, model interface{}) error {
	err := m.dbClient.Where("id = ?", id).First(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *BaseModel) GetList(page int, size int, models interface{}) error {
	offset := (page - 1) * size
	err := m.dbClient.Offset(offset).Limit(size).Order("created_at desc").Find(models).Error
	if err != nil {
		return err
	}

	return nil

}

func (m *BaseModel) GetTotal(model interface{}) (int64, error) {
	var count int64
	if err := m.dbClient.Model(model).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (m *BaseModel) CheckExist(id uint, model interface{}) (int, error) {
	var count int64
	if err := m.dbClient.Model(model).Where("id = ?", id).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if count == 0 {
		return http.StatusNotFound, fmt.Errorf("找不到資料")
	}
	return http.StatusOK, nil

}
