package middleware

import (
	"context"
	"fmt"
	"iot-platform/internal/repository"
	"net/http"
)

func AuthMiddleware(repo repository.DevicesRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			ctx := r.Context()
			apiKey := r.Header.Get("Api-Key")

			if id == "" {
				http.Error(w, "no device id is given", http.StatusBadRequest)
				return
			}
			if apiKey == "" {
				http.Error(w, "no api key is provided", http.StatusUnauthorized)
			}

			device, err := repo.FindDeviceById(ctx, id)
			if err != nil || device.ApiKey != apiKey {
				http.Error(w, fmt.Sprint("invalid device id or api key"), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "device", device)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
