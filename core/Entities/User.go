package entities

import (
	"fmt"
	"os"
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

func NewUserFromEnv(isCommitter bool) User {
	var nameEnv, emailEnv, dateEnv string

	if isCommitter {
		nameEnv = "SHIP_COMMITTER_NAME"
		emailEnv = "SHIP_COMMITTER_EMAIL"
		dateEnv = "SHIP_COMMITTER_DATE"
	} else {
		nameEnv = "SHIP_AUTHOR_NAME"
		emailEnv = "SHIP_AUTHOR_EMAIL"
		dateEnv = "SHIP_AUTHOR_DATE"
	}
	name := os.Getenv(nameEnv)
	email := os.Getenv(emailEnv)
	date := os.Getenv(dateEnv)

	if name == "" {
		name = os.Getenv("USER")
		if name == "" {
			name = "unknown"
		}
	}

	if email == "" {
		email = name + "@localhost"
	}

	var ts time.Time
	if date != "" {
		if t, err := time.Parse(time.RFC3339, date); err == nil {
			ts = t
		}
	}

	if ts.IsZero() {
		ts = time.Now()
	}

	return User{
		Name:      name,
		Email:     email,
		Timestamp: ts,
	}
}
