package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) BindRoutes() {
	api.Router.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(api.CORSMiddleware)
		r.Use(api.SessionManager.LoadAndSave)

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"service":   "pandoragym-api",
				"version":   "1.0.0",
			})
		})

		r.Post("/upload", api.UploadFile)

		r.Post("/session", api.AuthenticateWithPassword)
		r.Post("/register/student", api.CreateStudentAccount)
		r.Post("/register/personal", api.CreatePersonalAccount)
		r.Post("/password/recover", api.PasswordRecover)
		r.Post("/password/reset", api.ResetPassword)
		r.Post("/refresh", api.RefreshSession)
		r.Post("/revoke", api.RevokeToken)

		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware)

			r.Route("/users", func(r chi.Router) {
				r.Get("/profile", api.GetProfile)
				r.Put("/profile", api.UpdateProfile)
				r.Post("/avatar", api.UploadAvatar)
			})

			r.Route("/workouts", func(r chi.Router) {
				r.Get("/", api.GetWorkouts)
				r.Post("/", api.CreateWorkout)
				r.Get("/{id}", api.GetWorkout)
				r.Put("/{id}", api.UpdateWorkout)
				r.Delete("/{id}", api.DeleteWorkout)
				r.Post("/{id}/exercises", api.AddExerciseToWorkout)
				r.Delete("/{workoutId}/exercises/{exerciseId}", api.RemoveExerciseFromWorkout)

				r.Post("/{id}/finish", api.FinishWorkout)
				r.Post("/{id}/execute", api.ExecuteWorkout)
				r.Post("/{id}/rate", api.RateWorkout)

				r.Get("/history", api.GetWorkoutHistory)
				r.Get("/templates", api.GetWorkoutTemplates)
			})

			r.Route("/exercises", func(r chi.Router) {
				r.Get("/", api.GetExercises)
				r.Post("/", api.CreateExercise)
				r.Get("/{id}", api.GetExercise)
				r.Put("/{id}", api.UpdateExercise)
				r.Delete("/{id}", api.DeleteExercise)
				r.Get("/templates", api.GetExerciseTemplates)
			})

			r.Route("/programs", func(r chi.Router) {
				r.Get("/", api.GetAllTrainingPrograms)
				r.Get("/free", api.GetFreeTrainingPrograms)
				r.Get("/free/{id}", api.GetFreeTrainingProgramByID)
			})

			r.Route("/trainers", func(r chi.Router) {
				r.Get("/", api.GetPersonalTrainersList)
				r.Get("/{id}", api.GetPersonalTrainerByID)
				r.Post("/{id}/comments", api.AddPersonalTrainerComment)
				r.Get("/{id}/comments", api.GetPersonalTrainerComments)

				r.Group(func(r chi.Router) {
					r.Use(api.RequirePersonal)
					r.Get("/profile", api.GetPersonalProfile)
					r.Put("/profile", api.UpdatePersonalProfile)

					r.Get("/students", api.GetPersonalStudents)
					r.Post("/students", api.CreateStudent)
					r.Get("/students/{id}", api.GetStudentByID)
					r.Get("/students/{id}/workouts", api.GetStudentWorkouts)
					r.Get("/students/{id}/evolution", api.GetStudentEvolution)
					r.Delete("/students/{id}", api.RemoveStudent)

					r.Get("/plans", api.GetTrainerPlans)
					r.Post("/plans", api.CreatePlan)
					r.Put("/plans/{id}", api.UpdatePlan)
					r.Delete("/plans/{id}", api.DeletePlan)

					r.Post("/messages", api.SendMessage)
					r.Get("/schedule", api.GetPersonalSchedule)
					r.Post("/schedule", api.CreatePersonalSchedule)
				})
			})

			r.Route("/subscriptions", func(r chi.Router) {
				r.Post("/", api.SubscribeToPlan)
				r.Delete("/", api.CancelPlan)
			})

			r.Route("/schedulings", func(r chi.Router) {
				r.Get("/", api.GetSchedulings)
				r.Post("/", api.CreateScheduling)
				r.Get("/{id}", api.GetScheduling)
				r.Put("/{id}", api.UpdateScheduling)
				r.Delete("/{id}", api.CancelScheduling)
			})

			r.Route("/analytics", func(r chi.Router) {
				r.Get("/workout-frequency", api.GetWorkoutFrequency)
				r.Get("/workout-history", api.GetWorkoutHistoryExercises)
				r.Get("/workout-performance", api.GetWorkoutExercisePerformanceComparison)

				r.Group(func(r chi.Router) {
					r.Use(api.RequirePersonal)
					r.Get("/users/{userId}/workout-frequency", api.GetWorkoutFrequencyForUser)
					r.Get("/users/{userId}/workout-history", api.GetWorkoutHistoryForUser)
					r.Get("/users/{userId}/workout-performance", api.GetWorkoutPerformanceForUser)
				})
			})
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(api.AuthMiddleware)
			r.Use(api.RequireAdmin)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", api.GetAllUsers)
				r.Post("/status", api.UpdateUserStatus)
				r.Delete("/", api.DeleteUser)
				r.Get("/statistics", api.GetUserStatistics)
			})

			r.Get("/statistics", api.GetStatistics)
			r.Get("/reports", api.GetReports)
			r.Get("/system/health", api.GetSystemHealth)

			r.Route("/templates", func(r chi.Router) {
				r.Route("/exercises", func(r chi.Router) {
					r.Get("/", api.GetExerciseTemplatesAdmin)
					r.Post("/", api.CreateExerciseTemplate)
					r.Delete("/{id}", api.DeleteExerciseTemplate)
				})

				r.Route("/workouts", func(r chi.Router) {
					r.Get("/", api.GetWorkoutTemplatesAdmin)
					r.Post("/", api.CreateWorkoutTemplate)
					r.Delete("/{id}", api.DeleteWorkoutTemplate)
				})
			})
		})
	})
}
