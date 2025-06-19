package middleware

import (
	"context"
	"errors"
	"iot-platform/internal/database/postgres/device"
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
				return
			}

			dev, err := repo.FindDeviceById(ctx, id)
			if err != nil {
				if errors.Is(device.ErrDeviceNotFound, err) || apiKey != dev.ApiKey {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}

				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "device", dev)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
