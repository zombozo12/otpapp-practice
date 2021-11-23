package router

import (
	"github.com/go-chi/chi"
	http2 "otpapp-native/http"
)

func Init() *chi.Mux {
	rt := chi.NewRouter()
	rt.Route("/otp", func(r chi.Router) {
		r.Post("/request", http2.PhoneRequest)
		r.Post("/validate", http2.PhoneValidate)
	})

	return rt
}
