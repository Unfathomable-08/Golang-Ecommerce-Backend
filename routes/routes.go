package routes

import (
	"api/internal/controllers"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes sets up all the API endpoints expected by the frontend
func RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		// Products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", controllers.GetProducts)
			r.Get("/trending", controllers.GetTrendingProducts) // must come before /{id}
			r.Get("/{id}", controllers.GetProduct)
			r.Post("/", controllers.CreateProduct)
			r.Put("/{id}", controllers.UpdateProduct)
			r.Delete("/{id}", controllers.DeleteProduct)
		})

		// Carts
		r.Route("/carts", func(r chi.Router) {
			r.Get("/", controllers.GetCarts)
			r.Post("/", controllers.CreateCart)
		})

		// Orders
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", controllers.GetOrders)
			r.Post("/", controllers.CreateOrder)
			r.Get("/{id}", controllers.GetOrder)
			r.Put("/{id}", controllers.UpdateOrder)
			r.Post("/{id}/cancel", controllers.CancelOrder)
		})

		// Reviews
		r.Route("/reviews", func(r chi.Router) {
			r.Get("/", controllers.GetReviews)
			r.Post("/", controllers.CreateReview)
			r.Put("/{id}/status", controllers.UpdateReviewStatus)
		})
	})
}
