package utils

import "forum/config"

func InitServices() error {
	err := config.ConnectDatabase()
	if err != nil {
		return err
	}
	err = config.NewTemplateManager()
	if err != nil {
		return err
	}
	config.NewSessionManager()
	return nil
}
