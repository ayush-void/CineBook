package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	CustomerRole    Role = "CUSTOMER"
	HallManagerRole Role = "HALL_MANAGER"
	AdminRole       Role = "ADMIN"
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone             string    `gorm:"uniqueIndex"`
	Name              string
	Role              Role     `gorm:"default:'CUSTOMER'"`
	IsDisabled        bool     `gorm:"default:false"`
	AssignedScreenIDs []string `gorm:"type:text[]"` // For Hall Managers
	CreatedAt         time.Time
}

type Movie struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string
	Description string
	RuntimeMins int
	PosterURL   string
	ReleaseDate time.Time
	AgeRating   string
	Shows       []Show `gorm:"foreignKey:MovieID"`
}

type Screen struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TheatreID  uuid.UUID `gorm:"type:uuid"`
	Name       string
	ScreenType string // STANDARD, IMAX, FOUR_DX
	Capacity   int
	Seats      []Seat `gorm:"foreignKey:ScreenID"`
}

type Seat struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ScreenID  uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_screen_row_num"`
	Row       string    `gorm:"uniqueIndex:idx_screen_row_num"`
	Number    int       `gorm:"uniqueIndex:idx_screen_row_num"`
	Category  string
	BasePrice float64
}

type Show struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MovieID   uuid.UUID `gorm:"type:uuid"`
	ScreenID  uuid.UUID `gorm:"type:uuid;index:idx_screen_time"`
	StartTime time.Time `gorm:"index:idx_screen_time"`
	EndTime   time.Time
	BasePrice float64
	Bookings  []Booking `gorm:"foreignKey:ShowID"`
}

type Booking struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	ShowID      uuid.UUID `gorm:"type:uuid"`
	Status      string    `gorm:"default:'PENDING'"`
	TotalAmount float64
	CreatedAt   time.Time
}
