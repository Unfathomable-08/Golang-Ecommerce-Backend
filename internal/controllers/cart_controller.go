package controllers

import (
	"encoding/json"
	"net/http"

	"api/internal/services"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetCarts(w http.ResponseWriter, r *http.Request) {
	carts, err := services.GetCarts()
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, carts)
}

func CreateCart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CartID    string `json:"cartId"`
		ProductID string `json:"productId"`
		Qty       int    `json:"qty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Qty < 1 {
		body.Qty = 1
	}

	oid, err := bson.ObjectIDFromHex(body.ProductID)
	if err != nil {
		errResp(w, http.StatusBadRequest, "Invalid productId")
		return
	}

	cart, err := services.UpsertCart(body.CartID, oid, body.Qty)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, cart)
}
