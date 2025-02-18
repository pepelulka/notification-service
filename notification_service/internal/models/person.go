package models

type Person struct {
	PersonId    int        `json:"person_id"`
	Email       NullString `json:"email"`
	TelegramId  NullString `json:"telegram_id"`
	PhoneNumber NullString `json:"phone_number"`
}

type PersonCreate struct {
	Email       NullString `json:"email"`
	TelegramId  NullString `json:"telegram_id"`
	PhoneNumber NullString `json:"phone_number"`
}

type PersonCreateResult struct {
	PersonId int `json:"person_id"`
}
