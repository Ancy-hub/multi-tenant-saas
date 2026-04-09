package middleware

import (
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// RequireRole checks if the user has one of the allowed roles in the organization.
func RequireRole(service *services.MembershipService, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			orgIDParam := chi.URLParam(r, "id")
			if orgIDParam == "" {
				orgIDParam = chi.URLParam(r, "org_id")
			}

			orgID, err := uuid.Parse(orgIDParam)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, "Invalid org_id")
				return
			}

			role, err := service.GetUserRole(r.Context(), userID, orgID)
			if err != nil {
				utils.WriteError(w, http.StatusForbidden, "Access denied")
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.WriteError(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}
