package design

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Design struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ProfileId string `json:"profileId"`
	Fields    *Attrs `json:"fields"`
	Template  string `json:"design"`
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
