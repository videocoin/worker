package proto

import (
	time "time"

	"github.com/jinzhu/gorm"
)

func (w *WorkOrder) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("created_at", time.Now().Unix())
}
func (w *WorkOrder) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_at", time.Now().Unix())
}
