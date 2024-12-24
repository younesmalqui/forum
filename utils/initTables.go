package utils

import "forum/database"

func InitTables() error {
	tablesFn := database.Tables

	for _, fn := range tablesFn {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
