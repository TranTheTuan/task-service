package migration

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func migrationVersion2022032916010(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&Task{}) {
		if err := tx.AutoMigrate(&Task{}); err != nil {
			return err
		}
	}
	return nil
}
