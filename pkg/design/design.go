package design

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Design struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UserId    string    `json:"userId"`
	Fields    *Attrs    `json:"fields"`
	Template  string    `json:"design"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
