package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/middleware/auth"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workouts, err := api.WorkoutService.GetWorkouts(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get workouts", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workouts")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"workouts": workouts,
	})
}

func (api *API) GetWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	if workoutIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Workout ID is required")
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID format")
		return
	}

	workout, err := api.WorkoutService.GetWorkoutByID(r.Context(), workoutID, userID)
	if err != nil {
		api.Logger.Error("Failed to get workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "Workout not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"workout": workout,
	})
}

func (api *API) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CreateWorkoutParams](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workout, err := api.WorkoutService.CreateWorkout(r.Context(), userID, &req)
	if err != nil {
		api.Logger.Error("Failed to create workout", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"workout": workout,
	})
}

func (api *API) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	if workoutIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Workout ID is required")
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID format")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateWorkoutParams](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workout, err := api.WorkoutService.UpdateWorkout(r.Context(), workoutID, userID, &req)
	if err != nil {
		api.Logger.Error("Failed to update workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"workout": workout,
	})
}

func (api *API) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	if workoutIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Workout ID is required")
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID format")
		return
	}

	err = api.WorkoutService.DeleteWorkout(r.Context(), workoutID, userID)
	if err != nil {
		api.Logger.Error("Failed to delete workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Workout deleted successfully",
	})
}

// Exercise handlers
func (api *API) GetExercises(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exercises, err := api.WorkoutService.GetExercises(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get exercises", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get exercises")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"exercises": exercises,
	})
}

func (api *API) GetExercise(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exerciseIDStr := chi.URLParam(r, "id")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Exercise ID is required")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	exercise, err := api.WorkoutService.GetExerciseByID(r.Context(), exerciseID, userID)
	if err != nil {
		api.Logger.Error("Failed to get exercise", "error", err, "exercise_id", exerciseID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "Exercise not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"exercise": exercise,
	})
}

func (api *API) CreateExercise(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CreateExerciseParams](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	exercise, err := api.WorkoutService.CreateExercise(r.Context(), userID, &req)
	if err != nil {
		api.Logger.Error("Failed to create exercise", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"exercise": exercise,
	})
}

func (api *API) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exerciseIDStr := chi.URLParam(r, "id")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Exercise ID is required")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateExerciseParams](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	exercise, err := api.WorkoutService.UpdateExercise(r.Context(), exerciseID, userID, &req)
	if err != nil {
		api.Logger.Error("Failed to update exercise", "error", err, "exercise_id", exerciseID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"exercise": exercise,
	})
}

func (api *API) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exerciseIDStr := chi.URLParam(r, "id")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Exercise ID is required")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	err = api.WorkoutService.DeleteExercise(r.Context(), exerciseID, userID)
	if err != nil {
		api.Logger.Error("Failed to delete exercise", "error", err, "exercise_id", exerciseID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Exercise deleted successfully",
	})
}

// Exercise to workout handlers
func (api *API) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	if workoutIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Workout ID is required")
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID format")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.AddExerciseToWorkoutRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement adding exercise to workout
	_ = workoutID
	_ = req

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Exercise added to workout successfully",
	})
}

func (api *API) RemoveExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "workoutId")
	if workoutIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Workout ID is required")
		return
	}

	exerciseIDStr := chi.URLParam(r, "exerciseId")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Exercise ID is required")
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID format")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	// TODO: Implement removing exercise from workout
	_ = workoutID
	_ = exerciseID

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Exercise removed from workout successfully",
	})
}
