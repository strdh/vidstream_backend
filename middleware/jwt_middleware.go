package middleware

import (
    "os"
    "context"
    "net/http"
    "xyzstream/utils"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")

        claims, ok := utils.VerifyJWT(tokenString, os.Getenv("JWT_KEY"))
        if !ok {
            utils.WriteResponse(w, r, http.StatusUnauthorized, "Unauthorized", nil)
            return
        }

        idUser, ok := claims["sub"]
        if !ok {
            utils.WriteResponse(w, r, http.StatusUnauthorized, "Unauthorized", nil)
            return
        }

        ctx := context.WithValue(r.Context(), "idUser", idUser)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    }
}