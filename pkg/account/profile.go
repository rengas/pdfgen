package account

import "time"

type Profile struct {
	Id         string    `json:"id"`
	Email      string    `json:"email"`
	FirebaseId string    `json:"firebaseId"`
	Provider   string    `json:"provider"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
