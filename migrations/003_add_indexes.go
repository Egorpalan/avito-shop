package migrations

import (
	"gorm.io/gorm"
)

func AddIndexes(db *gorm.DB) error {
	if err := db.Exec("CREATE INDEX idx_merches_name ON merches(name)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX idx_inventories_user_id ON inventories(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX idx_inventories_merch_id ON inventories(merch_id)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX idx_transactions_from_user_id ON transactions(from_user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX idx_transactions_to_user_id ON transactions(to_user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX idx_transactions_amount ON transactions(amount)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX idx_users_username ON users(username)").Error; err != nil {
		return err
	}

	return nil
}
