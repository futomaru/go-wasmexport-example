package db

import (
	"context"
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Open() (*gorm.DB, func(), error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, nil, errors.New("DATABASE_URL required")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, _ := db.DB()
	return db, func() { sqlDB.Close() }, nil
}

func LoadRecords(ctx context.Context, db *gorm.DB) ([]Record, error) {
	var records []Record
	db.WithContext(ctx).Order("id").Find(&records)
	return records, nil
}
