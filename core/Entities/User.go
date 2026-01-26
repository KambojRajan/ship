package entities

import (
	"fmt"
	"time"
)

type User struct {
	Name      string
	Email     string
	Timestamp time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("%s <%s> %d %s", u.Name, u.Email, u.Timestamp.Unix(), u.Timestamp.Format("-0700"))
}
