package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/subscriptions"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

func GetSubscriptions(req *subscriptions.GetRequest) (*[]subscriptions.Subscription, error) {
	return subscriptions.Get(util.DB, req)
}
