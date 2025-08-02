package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
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
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	workout, exercises, err := api.WorkoutService.GetWorkoutByID(r.Context(), workoutID)
	if err != nil {
		api.Logger.Error("Failed to get workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout")
		return
	}

	if workout == nil || exercises == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Workout not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"workout":   workout,
		"exercises": exercises,
	})
}

func (api *API) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CreateWorkoutRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workout, err := api.WorkoutService.CreateWorkout(r.Context(), req, userID)
	if err != nil {
		api.Logger.Error("Failed to create workout", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, workout)
}

func (api *API) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateWorkoutRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workout, err := api.WorkoutService.UpdateWorkout(r.Context(), workoutID, req, userID)
	if err != nil {
		api.Logger.Error("Failed to update workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, workout)
}

func (api *API) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	err = api.WorkoutService.DeleteWorkout(r.Context(), workoutID, userID)
	if err != nil {
		api.Logger.Error("Failed to delete workout", "error", err, "workout_id", workoutID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Workout deleted successfully",
	})
}

// Exercise handlers (now part of workout service)

func (api *API) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := api.WorkoutService.GetAllExercises(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get exercises", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get exercises")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"exercises": exercises,
	})
}

func (api *API) GetExercise(w http.ResponseWriter, r *http.Request) {
	exerciseIDStr := chi.URLParam(r, "id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	exercise, err := api.WorkoutService.GetExerciseByID(r.Context(), exerciseID)
	if err != nil {
		api.Logger.Error("Failed to get exercise", "error", err, "exercise_id", exerciseID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get exercise")
		return
	}

	if exercise == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Exercise not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, exercise)
}

func (api *API) CreateExercise(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[pgstore.CreateExerciseRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	exercise, err := api.WorkoutService.CreateExercise(r.Context(), req)
	if err != nil {
		api.Logger.Error("Failed to create exercise", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, exercise)
}

func (api *API) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	exerciseIDStr := chi.URLParam(r, "id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateExerciseRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	exercise, err := api.WorkoutService.UpdateExercise(r.Context(), exerciseID, req)
	if err != nil {
		api.Logger.Error("Failed to update exercise", "error", err, "exercise_id", exerciseID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, exercise)
}

func (api *API) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	exerciseIDStr := chi.URLParam(r, "id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	err = api.WorkoutService.DeleteExercise(r.Context(), exerciseID)
	if err != nil {
		api.Logger.Error("Failed to delete exercise", "error", err, "exercise_id", exerciseID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete exercise")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Exercise deleted successfully",
	})
}

// Exercise templates

func (api *API) GetExerciseTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := api.WorkoutService.GetExerciseTemplates(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get exercise templates", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get exercise templates")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
	})
}

// Workout templates

func (api *API) GetWorkoutTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := api.WorkoutService.GetWorkoutTemplates(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get workout templates", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout templates")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
	})
}

// Workout execution and tracking

func (api *API) FinishWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")

	req, err := utils.DecodeValidJSON[struct {
		Duration  int                      `json:"duration"`
		Exercises []map[string]interface{} `json:"exercises"`
		Notes     string                   `json:"notes,omitempty"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.WorkoutService.FinishWorkout(r.Context(), userID.String(), workoutIDStr, req.Duration, req.Exercises, req.Notes)
	if err != nil {
		api.Logger.Error("Failed to finish workout", "error", err, "workout_id", workoutIDStr, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to finish workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Workout completed successfully",
	})
}

func (api *API) ExecuteWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")

	result, err := api.WorkoutService.ExecuteWorkout(r.Context(), userID.String(), workoutIDStr)
	if err != nil {
		api.Logger.Error("Failed to execute workout", "error", err, "workout_id", workoutIDStr, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to execute workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, result)
}

func (api *API) RateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")

	req, err := utils.DecodeValidJSON[struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment,omitempty"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.WorkoutService.RateWorkout(r.Context(), userID.String(), workoutIDStr, req.Rating, req.Comment)
	if err != nil {
		api.Logger.Error("Failed to rate workout", "error", err, "workout_id", workoutIDStr, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to rate workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Workout rated successfully",
	})
}

func (api *API) GetWorkoutHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	history, err := api.WorkoutService.GetWorkoutHistory(r.Context(), userID.String())
	if err != nil {
		api.Logger.Error("Failed to get workout history", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout history")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"history": history,
	})
}

// Exercise-workout relationships

func (api *API) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "id")
	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		ExerciseID string `json:"exercise_id"`
		Sets       int    `json:"sets"`
		Reps       int    `json:"reps"`
		RestTime   *int   `json:"rest_time,omitempty"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	exerciseID, err := uuid.Parse(req.ExerciseID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	err = api.WorkoutService.AddExerciseToWorkout(r.Context(), workoutID, exerciseID, userID, req.Sets, req.Reps, req.RestTime)
	if err != nil {
		api.Logger.Error("Failed to add exercise to workout", "error", err, "workout_id", workoutID, "exercise_id", exerciseID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add exercise to workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Exercise added to workout successfully",
	})
}

func (api *API) RemoveExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	workoutIDStr := chi.URLParam(r, "workoutId")
	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	exerciseIDStr := chi.URLParam(r, "exerciseId")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	err = api.WorkoutService.RemoveExerciseFromWorkout(r.Context(), workoutID, exerciseID, userID)
	if err != nil {
		api.Logger.Error("Failed to remove exercise from workout", "error", err, "workout_id", workoutID, "exercise_id", exerciseID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to remove exercise from workout")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Exercise removed from workout successfully",
	})
}
