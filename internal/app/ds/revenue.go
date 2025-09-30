package ds

import "time"

type Period struct {
	ID                  int `gorm:"primaryKey"`
	Title               string
	Description         string
	Img                 string
	Duration            string
	ShortDescription    string
	DetailedDescription string
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type PeriodsApplication struct {
	ID           int `gorm:"primaryKey"`
	Year         int
	FullDuration int
	Status       string
	ModeratorID  *int
	CreatorID    int
	CompanyName  string
	CreatedAt    time.Time
	FormedAt     time.Time
	CompletedAt  *time.Time
}

type SelectedPeriod struct {
	ApplicationID     int `gorm:"primaryKey"`
	PeriodID          int `gorm:"primaryKey"`
	PreviousRevenue   string
	ForecastedRevenue float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type User struct {
	ID          int `gorm:"primaryKey"`
	Login       string
	IsModerator bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
