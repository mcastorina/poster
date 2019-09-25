package store

import "github.com/mcastorina/poster/internal/models"

func GetResourceByName(name string) (models.Runnable, error) {
	resource, err := GetRequestByName(name)
	return &resource, err
}
