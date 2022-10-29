package user

import (
	"database/sql/driver"
	"fmt"
)

type Role int

const (
	RoleNormal Role = iota
	RoleAdmin
)

var roleName = map[Role]string{
	RoleNormal: "normal",
	RoleAdmin:  "admin",
}

var roleValue = map[string]Role{
	"normal": RoleNormal,
	"admin":  RoleAdmin,
}

func (r Role) String() string {
	return roleName[r]
}

func (r *Role) FromString(str string) {
	*r = roleValue[str]
}

func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Role) UnmarshalText(data []byte) error {
	r.FromString(string(data))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return driver.Value(r.String()), nil
}

func (r *Role) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok {
		r.FromString(string(v))
		return nil
	}
	return fmt.Errorf("failed to convert %v to Role value", value)
}
