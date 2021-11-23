package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	http2 "otpapp-native/http"
	"time"
)

func Init() *chi.Mux {
	rt := chi.NewRouter()
	rt.Use(middleware.Throttle(1))
	rt.Use(httprate.LimitAll(5, 1*time.Minute))
	rt.Route("/otp", func(r chi.Router) {
		r.Post("/request", http2.PhoneRequest)
		r.Post("/validate", http2.PhoneValidate)
	})

	return rt
}
