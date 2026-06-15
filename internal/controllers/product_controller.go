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

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errResp(w http.ResponseWriter, status int, msg string) {
	respond(w, status, map[string]string{"error": msg})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 12
	}

	var minP, maxP *float64
	if v := q.Get("minPrice"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			minP = &f
		}
	}
	if v := q.Get("maxPrice"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			maxP = &f
		}
	}

	result, err := services.GetProducts(page, limit, q.Get("category"), q.Get("search"), minP, maxP)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, result)
}

func GetTrendingProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 8
	}
	products, err := services.GetTrendingProducts(limit)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := services.GetProduct(id)
	if err != nil {
		errResp(w, http.StatusNotFound, "Product not found")
		return
	}
	respond(w, http.StatusOK, product)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.CreateProduct(&p)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, map[string]any{"success": true, "product": result})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var update bson.M
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		errResp(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.UpdateProduct(id, update)
	if err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]any{"success": true, "product": result})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := services.DeleteProduct(id); err != nil {
		errResp(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]any{"success": true})
}
