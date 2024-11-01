package middlewares

import (
	"net/http"
	"strings"

	"github.com/Maden-in-haven/crmlib/pkg/myjwt"
)

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		// Извлекаем сам токен, убирая "Bearer " перед ним
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Парсим токен и проверяем его валидность
		claims, err := myjwt.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Проверяем claims токена
		tokenType, exists := claims["typ"].(string)
		if !exists {
			http.Error(w, "Неверный токен: тип не определён", http.StatusForbidden)
			return
		}

		// Проверяем, что это access токен
		if tokenType != "access" {
			http.Error(w, "Неверный токен: ожидался access токен", http.StatusForbidden)
			return
		}

		userID, exists := claims["sub"].(string)
		if !exists {
			http.Error(w, "Неверный токен: тип не определён", http.StatusForbidden)
			return
		}
		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}
