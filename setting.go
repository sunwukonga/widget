package widget

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/serializable_meta"
)

type QorWidgetSetting struct {
	gorm.Model
	Scope string
	key   string
	serializable_meta.SerializableMeta
}

func findSettingByNameAndKey(kind string, key string) (setting *QorWidgetSetting) {
	return &QorWidgetSetting{}
}
