package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api/internal/models"
	"api/internal/services"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	result, err := services.GetOrders(page, limit)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, result)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	order, err := services.GetOrder(id)
	if err != nil {
		errResp(w, http.StatusNotFound, "Order not found")
		return
	}
	respond(w, http.StatusOK, order)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.CreateOrder(&o)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, map[string]any{"success": true, "order": result})
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var update bson.M
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.UpdateOrder(id, update)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]any{"success": true, "order": result})
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	order, err := services.CancelOrder(id)
	if err != nil {
		if err == services.ErrCannotCancel {
			errResp(w, http.StatusBadRequest, err.Error())
		} else {
			errResp(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respond(w, http.StatusOK, map[string]any{"success": true, "order": order})
}
