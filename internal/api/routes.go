package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	// Global middleware
	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.Recoverer)
	api.Router.Use(middleware.RequestID)
	api.Router.Use(middleware.RealIP)
	api.Router.Use(CORSMiddleware)

	// Health check
	api.Router.Get("/health", api.HealthCheck)

	// Auth routes (public)
	api.Router.Route("/auth", func(r chi.Router) {
		r.Post("/session", api.AuthenticateWithPassword)
		r.Post("/register/student", api.CreateStudentAccount)
		r.Post("/register/personal", api.CreatePersonalAccount)
		r.Post("/password/recover", api.PasswordRecover)
		r.Post("/password/reset", api.ResetPassword)
		r.Post("/refresh", api.RefreshToken)
		r.Post("/revoke", api.RevokeToken)
	})

	// Protected routes
	api.Router.Route("/api", func(r chi.Router) {
		r.Use(api.JWTMiddleware)

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/profile", api.GetProfile)
			r.Put("/profile", api.UpdateProfile)
			r.Post("/avatar", api.UploadAvatar)
		})

		// Workout routes
		r.Route("/workouts", func(r chi.Router) {
			r.Get("/", api.GetWorkouts)
			r.Post("/", api.CreateWorkout)
			r.Get("/{id}", api.GetWorkout)
			r.Put("/{id}", api.UpdateWorkout)
			r.Delete("/{id}", api.DeleteWorkout)
			r.Post("/{id}/exercises", api.AddExerciseToWorkout)
			r.Delete("/{workoutId}/exercises/{exerciseId}", api.RemoveExerciseFromWorkout)
		})

		// Exercise routes
		r.Route("/exercises", func(r chi.Router) {
			r.Get("/", api.GetExercises)
			r.Post("/", api.CreateExercise)
			r.Get("/{id}", api.GetExercise)
			r.Put("/{id}", api.UpdateExercise)
			r.Delete("/{id}", api.DeleteExercise)
		})

		// Scheduling routes
		r.Route("/schedulings", func(r chi.Router) {
			r.Get("/", api.GetSchedulings)
			r.Post("/", api.CreateScheduling)
			r.Get("/{id}", api.GetScheduling)
			r.Put("/{id}", api.UpdateScheduling)
			r.Delete("/{id}", api.CancelScheduling)
		})

		// Personal trainer specific routes
		r.Route("/personal", func(r chi.Router) {
			r.Use(api.PersonalOnlyMiddleware)
			r.Get("/students", api.GetPersonalStudents)
			r.Get("/students/{id}/evolution", api.GetStudentEvolution)
			r.Post("/messages", api.SendMessage)
			r.Get("/schedule", api.GetPersonalSchedule)
			r.Post("/schedule", api.CreatePersonalSchedule)
		})

		// Student specific routes
		r.Route("/student", func(r chi.Router) {
			r.Use(api.StudentOnlyMiddleware)
			r.Get("/workouts/history", api.GetWorkoutHistory)
			r.Post("/workouts/{id}/execute", api.ExecuteWorkout)
			r.Post("/workouts/{id}/rate", api.RateWorkout)
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(api.AdminOnlyMiddleware)
			r.Get("/users", api.GetAllUsers)
			r.Get("/statistics", api.GetStatistics)
		})

		// File upload
		r.Post("/upload", api.UploadFile)
	})
}
