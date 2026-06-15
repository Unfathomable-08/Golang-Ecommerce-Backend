package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api/internal/models"
	"api/internal/services"

	"github.com/go-chi/chi/v5"
)

func GetReviews(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	result, err := services.GetReviews(page, limit, q.Get("status"))
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, result)
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	var rv models.Review
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.CreateReview(&rv)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, map[string]any{"success": true, "review": result})
}

func UpdateReviewStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.UpdateReviewStatus(id, body.Status)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]any{"success": true, "review": result})
}
