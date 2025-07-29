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

	// File upload (public)
	api.Router.Post("/upload", api.UploadFile)

	// Auth routes (public)
	api.Router.Route("/auth", func(r chi.Router) {
		r.Post("/session", api.AuthenticateWithPassword)
		r.Post("/register/student", api.CreateStudentAccount)
		r.Post("/register/personal", api.CreatePersonalAccount)
		r.Post("/password/recover", api.PasswordRecover)
		r.Post("/password/reset", api.ResetPassword)
		r.Post("/refresh", api.RefreshToken)
		r.Post("/revoke", api.RevokeToken)
		
		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(api.JWTMiddleware)
			r.Get("/profile", api.GetProfile)
			r.Put("/profile", api.UpdateProfile)
			r.Post("/avatar", api.UploadAvatar)
		})
	})

	// Protected API routes
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
			
			// Workout execution and tracking
			r.Post("/{id}/finish", api.FinishWorkout)
			r.Post("/{id}/execute", api.ExecuteWorkout)
			r.Post("/{id}/rate", api.RateWorkout)
			
			// Workout history and analytics
			r.Get("/history", api.GetWorkoutHistory)
			r.Get("/templates", api.GetWorkoutTemplates)
		})

		// Exercise routes
		r.Route("/exercises", func(r chi.Router) {
			r.Get("/", api.GetExercises)
			r.Post("/", api.CreateExercise)
			r.Get("/{id}", api.GetExercise)
			r.Put("/{id}", api.UpdateExercise)
			r.Delete("/{id}", api.DeleteExercise)
			r.Get("/templates", api.GetExerciseTemplates)
		})

		// Training programs (free and premium)
		r.Route("/programs", func(r chi.Router) {
			r.Get("/", api.GetAllTrainingPrograms)
			r.Get("/free", api.GetFreeTrainingPrograms)
			r.Get("/free/{id}", api.GetFreeTrainingProgramByID)
		})

		// Personal trainer routes
		r.Route("/trainers", func(r chi.Router) {
			r.Get("/", api.GetPersonalTrainersList)
			r.Get("/{id}", api.GetPersonalTrainerByID)
			r.Post("/{id}/comments", api.AddPersonalTrainerComment)
			r.Get("/{id}/comments", api.GetPersonalTrainerComments)
			
			// Trainer-specific routes (requires personal trainer role)
			r.Group(func(r chi.Router) {
				r.Use(api.PersonalOnlyMiddleware)
				r.Get("/profile", api.GetPersonalProfile)
				r.Put("/profile", api.UpdatePersonalProfile)
				
				// Student management
				r.Get("/students", api.GetPersonalStudents)
				r.Post("/students", api.CreateStudent)
				r.Get("/students/{id}", api.GetStudentByID)
				r.Get("/students/{id}/workouts", api.GetStudentWorkouts)
				r.Get("/students/{id}/evolution", api.GetStudentEvolution)
				r.Delete("/students/{id}", api.RemoveStudent)
				
				// Plan management
				r.Get("/plans", api.GetTrainerPlans)
				r.Post("/plans", api.CreatePlan)
				r.Put("/plans/{id}", api.UpdatePlan)
				r.Delete("/plans/{id}", api.DeletePlan)
				
				// Messaging
				r.Post("/messages", api.SendMessage)
				r.Get("/schedule", api.GetPersonalSchedule)
				r.Post("/schedule", api.CreatePersonalSchedule)
			})
		})

		// Subscription management
		r.Route("/subscriptions", func(r chi.Router) {
			r.Post("/", api.SubscribeToPlan)
			r.Delete("/", api.CancelPlan)
		})

		// Scheduling routes
		r.Route("/schedulings", func(r chi.Router) {
			r.Get("/", api.GetSchedulings)
			r.Post("/", api.CreateScheduling)
			r.Get("/{id}", api.GetScheduling)
			r.Put("/{id}", api.UpdateScheduling)
			r.Delete("/{id}", api.CancelScheduling)
		})

		// Analytics routes
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/workout-frequency", api.GetWorkoutFrequency)
			r.Get("/workout-history", api.GetWorkoutHistoryExercises)
			r.Get("/workout-performance", api.GetWorkoutExercisePerformanceComparison)
			
			// Trainer analytics for specific users (requires trainer role)
			r.Group(func(r chi.Router) {
				r.Use(api.PersonalOnlyMiddleware)
				r.Get("/users/{userId}/workout-frequency", api.GetWorkoutFrequencyForUser)
				r.Get("/users/{userId}/workout-history", api.GetWorkoutHistoryForUser)
				r.Get("/users/{userId}/workout-performance", api.GetWorkoutPerformanceForUser)
			})
		})
	})

	// Admin routes
	api.Router.Route("/admin", func(r chi.Router) {
		r.Use(api.JWTMiddleware)
		r.Use(api.AdminOnlyMiddleware)
		
		// User management
		r.Route("/users", func(r chi.Router) {
			r.Get("/", api.GetAllUsers)
			r.Post("/status", api.UpdateUserStatus)
			r.Delete("/", api.DeleteUser)
			r.Get("/statistics", api.GetUserStatistics)
		})
		
		// Platform management
		r.Get("/statistics", api.GetStatistics)
		r.Get("/reports", api.GetReports)
		r.Get("/system/health", api.GetSystemHealth)
		
		// Template management
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
}
