package api

import (
	"net"
	"net/http"

	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) AuthenticateWithPassword(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.AuthenticateRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Extract IP address from request
	if req.IPAddress == nil {
		ipAddr := ExtractIPAddress(r.RemoteAddr)
		req.IPAddress = ipAddr
	}

	// Extract device info from User-Agent if not provided
	if req.DeviceInfo == nil {
		userAgent := r.Header.Get("User-Agent")
		if userAgent != "" {
			req.DeviceInfo = &userAgent
		}
	}

	response, err := api.AuthService.Authenticate(r.Context(), req)
	if err != nil {
		api.Logger.Error("Authentication failed", "error", err)
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// ExtractIPAddress extracts IP address from request
func ExtractIPAddress(remoteAddr string) *string {
	if remoteAddr == "" {
		return nil
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If splitting fails, assume it's just an IP
		host = remoteAddr
	}

	return &host
}

func (api *API) CreateStudentAccount(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.CreateStudentWithUserRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.UserService.CreateStudent(r.Context(), req)
	if err != nil {
		api.Logger.Error("Failed to create student account", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"user": user,
	})
}

func (api *API) CreatePersonalAccount(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.CreatePersonalWithUserRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.UserService.CreatePersonal(r.Context(), req)
	if err != nil {
		api.Logger.Error("Failed to create personal account", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"user": user,
	})
}

func (api *API) PasswordRecover(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.PasswordRecoverRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.AuthService.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		api.Logger.Error("Failed to request password reset", "error", err)
		// Don't reveal if email exists or not
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "If the email exists, a password reset link has been sent",
	})
}

func (api *API) ResetPassword(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.ResetPasswordRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.AuthService.ResetPassword(r.Context(), req.Token, req.Password)
	if err != nil {
		api.Logger.Error("Failed to reset password", "error", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid or expired token")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Password reset successfully",
	})
}

func (api *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.RefreshTokenRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := api.AuthService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		api.Logger.Error("Failed to refresh token", "error", err)
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (api *API) RevokeToken(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.RefreshTokenRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.AuthService.RevokeRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		api.Logger.Error("Failed to revoke token", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Token revoked successfully",
	})
}
