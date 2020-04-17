package email

import (
	"strings"
)

// SendWelcome sends a welcome email!
func SendWelcome(email string, name string) {
	firstName := strings.Split(name, " ")[0]

	send(
		welcome,
		map[string]interface{}{"first_name": firstName},
		email,
		name,
	)
}
