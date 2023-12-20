package models

import (
	"errors"
	"gorm.io/gorm"
)

func IgnoreRecordNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return err
}