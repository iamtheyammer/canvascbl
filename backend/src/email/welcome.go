package email

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"strings"
	"time"
)

func SendWelcome(email string, name string) {
	firstName := strings.Split(name, " ")[0]

	send(
		welcome,
		map[string]interface{}{"first_name": firstName},
		email,
		name,
	)
}

func SendWelcomeIfNecessary(usersResponse *string) {
	user := db.GetUserFromCanvasProfileResponseJSON(usersResponse)
	if user == nil {
		return
	}

	if !user.InsertedAt.Add(time.Second * 5).After(time.Now()) {
		return
	}

	SendWelcome(user.Email, user.Name)
}
