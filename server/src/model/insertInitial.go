package model

import (
	"gorm.io/gorm"
)

// 出入初始化数据
func InsertInitial(db *gorm.DB) error {
	// 需要运行初始化插入的sql
	var runList []func(db *gorm.DB) error = []func(db *gorm.DB) error{}
	for _, f := range runList {
		err := f(db)
		if err != nil {
			return err
		}
	}
	return nil
}

// 检查是否存在如果不存在就创建
func checkAndCreate[T any](db *gorm.DB, models ...T) error {
	for _, v := range models {
		// 先检查有没有
		res := db.Debug().Find(v)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected <= 0 {
			// 创建
			res = db.Debug().Create(v)
			if res.Error != nil {
				return res.Error
			}
		}
	}

	return nil
}
