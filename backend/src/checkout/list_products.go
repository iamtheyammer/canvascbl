package checkout

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

type products []product

type product struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	ShortName string  `json:"shortName"`
	Price     float64 `json:"price"`
	Type      string  `json:"type"`
}

func ListProductsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cp, err := db.CheckoutListProducts()
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing products"))
		util.SendInternalServerError(w)
		return
	}

	var prods products

	for _, p := range *cp {
		prods = append(prods, product{p.ID, p.Name, p.ShortName, p.Price, p.Type})
	}

	j, err := json.Marshal(prods)

	util.SendJSONResponse(w, j)
}
