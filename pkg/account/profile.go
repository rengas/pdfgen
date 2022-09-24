package account

type Profile struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	FirebaseId string `json:"firebaseId"`
	Provider   string `json:"provider"`
}
