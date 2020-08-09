package gophercon2020

//go:generate repogen

//repogen:entity
type User struct {
	ID           uint `gorm:"primary_key"`
	Email        string
	PasswordHash string
}
