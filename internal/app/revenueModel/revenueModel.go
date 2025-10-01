package revenueModel

import (
   "gorm.io/driver/postgres"
   "gorm.io/gorm"
)

type RevenueModel struct {
	db *gorm.DB
}

func New(dsn string) (*RevenueModel, error) {
   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) 
   if err != nil {
      return nil, err
   }

   return &RevenueModel{
      db: db,
   }, nil
}

func NewRevenueModel() (*RevenueModel, error) {
	return &RevenueModel{}, nil
}

