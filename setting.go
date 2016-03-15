package widget

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/serializable_meta"
)

type QorWidgetSetting struct {
	gorm.Model
	Scope string
	serializable_meta.SerializableMeta
}
