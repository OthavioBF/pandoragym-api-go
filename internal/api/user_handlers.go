package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

// Authentication handlers

func (api *API) AuthenticateWithPassword(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.AuthenticateWithPasswordRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, return mock response
	// TODO: Implement actual authentication
	api.Logger.Info("Authentication requested", "email", req.Email)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message":       "Authentication not fully implemented yet",
		"access_token":  "mock-access-token",
		"refresh_token": "mock-refresh-token",
	})
}

func (api *API) CreateStudentAccount(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.CreateStudentWithUserRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.UserService.CreateStudentWithUser(r.Context(), req)
	if err != nil {
		api.Logger.Error("Failed to create student account", "error", err, "email", req.Email)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, user)
}

func (api *API) CreatePersonalAccount(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.CreatePersonalWithUserRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.UserService.CreatePersonalWithUser(r.Context(), req)
	if err != nil {
		api.Logger.Error("Failed to create personal trainer account", "error", err, "email", req.Email)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, user)
}

func (api *API) PasswordRecover(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Password recovery not implemented yet",
	})
}

func (api *API) ResetPassword(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Password reset not implemented yet",
	})
}

func (api *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Token refresh not implemented yet",
	})
}

func (api *API) RevokeToken(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Token revocation not implemented yet",
	})
}

// Profile management handlers

func (api *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := api.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get user profile", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	if user == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, user)
}

func (api *API) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateProfileRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.UserService.UpdateUserProfile(r.Context(), userID, &req)
	if err != nil {
		api.Logger.Error("Failed to update user profile", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Profile updated successfully",
	})
}

func (api *API) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// TODO: Implement file upload to storage service (S3, Supabase, etc.)
	// For now, just return success with mock URL
	avatarURL := "https://example.com/avatars/" + userID.String() + "_" + header.Filename

	api.Logger.Info("Avatar upload requested", "user_id", userID, "filename", header.Filename)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
	})
}

func (api *API) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// TODO: Implement file upload to storage service (S3, Supabase, etc.)
	// For now, just return success with mock URL
	fileURL := "https://example.com/files/" + header.Filename

	api.Logger.Info("File upload requested", "filename", header.Filename)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
}

// Personal trainer discovery and interaction

func (api *API) GetPersonalTrainersList(w http.ResponseWriter, r *http.Request) {
	// For now, return empty array
	// TODO: Implement actual trainer retrieval
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"trainers": []interface{}{},
	})
}

func (api *API) GetPersonalTrainerByID(w http.ResponseWriter, r *http.Request) {
	trainerID := chi.URLParam(r, "id")

	// For now, return not found
	// TODO: Implement actual trainer retrieval
	api.Logger.Info("Personal trainer requested", "trainer_id", trainerID)
	utils.WriteErrorResponse(w, http.StatusNotFound, "Personal trainer not found")
}

func (api *API) AddPersonalTrainerComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	trainerID := chi.URLParam(r, "id")

	req, err := utils.DecodeValidJSON[struct {
		Comment string `json:"comment"`
		Rating  int    `json:"rating"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, just return success
	// TODO: Implement actual comment addition
	api.Logger.Info("Trainer comment requested", "trainer_id", trainerID, "user_id", userID, "rating", req.Rating)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Comment added successfully",
	})
}

func (api *API) GetPersonalTrainerComments(w http.ResponseWriter, r *http.Request) {
	trainerID := chi.URLParam(r, "id")

	// For now, return empty array
	// TODO: Implement actual comment retrieval
	api.Logger.Info("Trainer comments requested", "trainer_id", trainerID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"comments": []interface{}{},
	})
}

// Personal trainer profile management (requires trainer role)

func (api *API) GetPersonalProfile(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// For now, return mock profile
	// TODO: Implement actual profile retrieval
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"id":           trainerID.String(),
		"name":         "John Doe",
		"email":        "john@example.com",
		"bio":          "Certified personal trainer",
		"specialties":  []string{"Weight Loss", "Muscle Building"},
		"experience":   5,
	})
}

func (api *API) UpdatePersonalProfile(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Phone       string   `json:"phone"`
		Bio         string   `json:"bio"`
		Specialties []string `json:"specialties"`
		Experience  int      `json:"experience"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, just return success
	// TODO: Implement actual profile update
	api.Logger.Info("Trainer profile update requested", "trainer_id", trainerID, "name", req.Name)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Profile updated successfully",
	})
}

// Student management (requires trainer role)

func (api *API) GetPersonalStudents(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// For now, return empty array
	// TODO: Implement actual student retrieval
	api.Logger.Info("Trainer students requested", "trainer_id", trainerID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"students": []interface{}{},
	})
}

func (api *API) CreateStudent(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		BornDate string `json:"born_date"`
		Age      int    `json:"age"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, return mock student
	// TODO: Implement actual student creation
	api.Logger.Info("Student creation requested", "trainer_id", trainerID, "student_name", req.Name)

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"id":      uuid.New().String(),
		"name":    req.Name,
		"email":   req.Email,
		"message": "Student created successfully",
	})
}

func (api *API) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	studentID := chi.URLParam(r, "id")

	// For now, return not found
	// TODO: Implement actual student retrieval
	api.Logger.Info("Student details requested", "trainer_id", trainerID, "student_id", studentID)
	utils.WriteErrorResponse(w, http.StatusNotFound, "Student not found")
}

func (api *API) GetStudentWorkouts(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	studentID := chi.URLParam(r, "id")

	// For now, return empty array
	// TODO: Implement actual student workout retrieval
	api.Logger.Info("Student workouts requested", "trainer_id", trainerID, "student_id", studentID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"workouts": []interface{}{},
	})
}

func (api *API) GetStudentEvolution(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Student evolution tracking not implemented yet",
	})
}

func (api *API) RemoveStudent(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	studentID := chi.URLParam(r, "id")

	// For now, just return success
	// TODO: Implement actual student removal
	api.Logger.Info("Student removal requested", "trainer_id", trainerID, "student_id", studentID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Student removed successfully",
	})
}

// Communication and messaging (requires trainer role)

func (api *API) SendMessage(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Messaging system not implemented yet",
	})
}

func (api *API) GetPersonalSchedule(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Personal trainer schedule not implemented yet",
	})
}

func (api *API) CreatePersonalSchedule(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Personal trainer schedule creation not implemented yet",
	})
}

// Admin user management (admin only)

func (api *API) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination and filtering
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	role := r.URL.Query().Get("role")
	search := r.URL.Query().Get("search")

	users, total, err := api.UserService.GetAllUsers(page, limit, role, search)
	if err != nil {
		api.Logger.Error("Failed to get all users", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (api *API) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[struct {
		UserID string `json:"user_id"`
		Status string `json:"status"` // active, suspended, banned
		Reason string `json:"reason,omitempty"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.UserService.UpdateUserStatus(req.UserID, req.Status, req.Reason)
	if err != nil {
		api.Logger.Error("Failed to update user status", "error", err, "user_id", req.UserID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "User status updated successfully",
	})
}

func (api *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.UserService.DeleteUser(req.UserID, req.Reason)
	if err != nil {
		api.Logger.Error("Failed to delete user", "error", err, "user_id", req.UserID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

func (api *API) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	stats, err := api.AnalyticsService.GetUserStatistics(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get user statistics", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get user statistics")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
// Plan management (requires trainer role)

func (api *API) GetTrainerPlans(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	plans, err := api.PlanService.GetTrainerPlans(r.Context(), trainerID)
	if err != nil {
		api.Logger.Error("Failed to get trainer plans", "error", err, "trainer_id", trainerID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get plans")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

func (api *API) CreatePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"` // in days
		Features    []string `json:"features"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := api.PlanService.CreatePlan(r.Context(), trainerID, req.Name, req.Description, req.Price, req.Duration, req.Features)
	if err != nil {
		api.Logger.Error("Failed to create plan", "error", err, "trainer_id", trainerID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, plan)
}

func (api *API) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	planID := chi.URLParam(r, "id")

	req, err := utils.DecodeValidJSON[struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"`
		Features    []string `json:"features"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := api.PlanService.UpdatePlan(r.Context(), trainerID, planID, req.Name, req.Description, req.Price, req.Duration, req.Features)
	if err != nil {
		api.Logger.Error("Failed to update plan", "error", err, "trainer_id", trainerID, "plan_id", planID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, plan)
}

func (api *API) DeletePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	planID := chi.URLParam(r, "id")

	err := api.PlanService.DeletePlan(r.Context(), trainerID, planID)
	if err != nil {
		api.Logger.Error("Failed to delete plan", "error", err, "trainer_id", trainerID, "plan_id", planID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Plan deleted successfully",
	})
}
