package ds

import "time"

type Period struct {
	ID                  int     `gorm:"primaryKey" json:"id"`
	Title               *string `json:"title,omitempty"`
	Description         *string `json:"description,omitempty"`
	Img                 *string `json:"img",omitempty`
	Duration            string  `json:"duration"`
	ShortDescription    *string `json:"short_description"`
	DetailedDescription *string `json:"detailed_description"`
	IsActive            *bool   `json:"is_active,omitempty"`
}

type PeriodsApplication struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Year        *int       `json:"year,omitempty"`
	Status      string     `json:"status"`
	ModeratorID *int       `json:"moderator_id,omitempty"`
	CreatorID   int        `json:"creator_id"`
	CompanyName *string    `json:"company_name,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	FormedAt    *time.Time `json:"formed_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type SelectedPeriod struct {
	ApplicationID     int       `gorm:"primaryKey" json:"applicationId"`
	PeriodID          int       `gorm:"primaryKey" json:"periodId"`
	PreviousRevenue   string    `json:"previousRevenue"`
	ForecastedRevenue float64   `json:"forecastedRevenue"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type User struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	IsModerator bool      `json:"isModerator"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
