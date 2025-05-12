package middlewarex

import (
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/utils"
	"net/http"
)

func UserID(Key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.URL.Query().Get("user_id")
			if userID == "" {
				next.ServeHTTP(w, r)
				return
			}
			maskUserID, _ := utils.Encrypt(userID, Key)
			r = r.WithContext(contextx.WithUserID(r.Context(), contextx.UserID(userID)))
			r = r.WithContext(contextx.WithMaskUserID(r.Context(), contextx.MaskUserID(maskUserID)))
			next.ServeHTTP(w, r)
		})
	}
}
