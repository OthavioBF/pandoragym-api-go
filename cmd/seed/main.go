package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize database connection
	pool, err := pgstore.InitDB(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create queries instance
	queries := pgstore.NewQueries(pool)
	ctx := context.Background()

	fmt.Println("Starting database seeding...")

	// Hash password for all users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	now := time.Now()

	// Create Personal Trainers
	personalTrainers := []struct {
		name           string
		email          string
		phone          string
		description    string
		experience     string
		specialization string
		qualifications string
	}{
		{
			name:           "Carlos Silva",
			email:          "carlos@pandoragym.com",
			phone:          "11999887766",
			description:    "Personal trainer especializado em hipertrofia e condicionamento físico",
			experience:     "5 anos de experiência em academias de alto padrão",
			specialization: "Hipertrofia e Força",
			qualifications: "CREF 123456-G/SP, Especialização em Fisiologia do Exercício",
		},
		{
			name:           "Ana Costa",
			email:          "ana@pandoragym.com",
			phone:          "11988776655",
			description:    "Especialista em treinamento funcional e reabilitação",
			experience:     "7 anos trabalhando com atletas e reabilitação",
			specialization: "Treinamento Funcional",
			qualifications: "CREF 654321-G/SP, Curso de Biomecânica Aplicada",
		},
		{
			name:           "Roberto Santos",
			email:          "roberto@pandoragym.com",
			phone:          "11977665544",
			description:    "Personal trainer focado em emagrecimento e condicionamento",
			experience:     "4 anos de experiência com público diverso",
			specialization: "Emagrecimento e Condicionamento",
			qualifications: "CREF 789012-G/SP, Certificação em Nutrição Esportiva",
		},
	}

	var personalIDs []uuid.UUID

	for _, pt := range personalTrainers {
		userID := uuid.New()
		personalIDs = append(personalIDs, userID)

		// Create user
		_, err := queries.CreateUser(ctx, pgstore.CreateUserParams{
			ID:        userID,
			Name:      pt.name,
			Email:     pt.email,
			Phone:     pt.phone,
			Password:  string(passwordHash),
			Role:      pgstore.RolePersonal,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			log.Printf("Failed to create personal trainer user %s: %v", pt.name, err)
			continue
		}

		// Create personal profile
		err = queries.CreatePersonal(ctx, pgstore.CreatePersonalParams{
			ID:             userID,
			Description:    &pt.description,
			Experience:     &pt.experience,
			Specialization: &pt.specialization,
			Qualifications: &pt.qualifications,
		})
		if err != nil {
			log.Printf("Failed to create personal trainer profile %s: %v", pt.name, err)
			continue
		}

		fmt.Printf("Created personal trainer: %s\n", pt.name)
	}

	// Create Students
	students := []struct {
		name                  string
		email                 string
		phone                 string
		age                   int
		weight                float64
		objective             string
		trainingFrequency     string
		didBodybuilding       bool
		medicalCondition      *string
		physicalActivityLevel *string
		observations          *string
	}{
		{
			name:                  "João Silva",
			email:                 "joao@email.com",
			phone:                 "11987654321",
			age:                   25,
			weight:                75.5,
			objective:             "Ganhar massa muscular",
			trainingFrequency:     "5x por semana",
			didBodybuilding:       false,
			physicalActivityLevel: stringPtr("Intermediário"),
			observations:          stringPtr("Sem restrições médicas"),
		},
		{
			name:                  "Maria Santos",
			email:                 "maria@email.com",
			phone:                 "11976543210",
			age:                   30,
			weight:                65.0,
			objective:             "Perder peso e tonificar",
			trainingFrequency:     "4x por semana",
			didBodybuilding:       false,
			medicalCondition:      stringPtr("Hipertensão controlada"),
			physicalActivityLevel: stringPtr("Iniciante"),
			observations:          stringPtr("Prefere exercícios de baixo impacto"),
		},
		{
			name:                  "Pedro Costa",
			email:                 "pedro@email.com",
			phone:                 "11965432109",
			age:                   28,
			weight:                80.0,
			objective:             "Melhorar condicionamento físico",
			trainingFrequency:     "3x por semana",
			didBodybuilding:       true,
			physicalActivityLevel: stringPtr("Avançado"),
			observations:          stringPtr("Ex-atleta de futebol"),
		},
	}

	var studentIDs []uuid.UUID

	for _, student := range students {
		userID := uuid.New()
		studentIDs = append(studentIDs, userID)

		// Create user
		_, err := queries.CreateUser(ctx, pgstore.CreateUserParams{
			ID:        userID,
			Name:      student.name,
			Email:     student.email,
			Phone:     student.phone,
			Password:  string(passwordHash),
			Role:      pgstore.RoleStudent,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			log.Printf("Failed to create student user %s: %v", student.name, err)
			continue
		}

		// Create student profile
		bornDate := time.Now().AddDate(-student.age, 0, 0)
		err = queries.CreateStudent(ctx, pgstore.CreateStudentParams{
			ID:                    userID,
			BornDate:              bornDate,
			Age:                   int32(student.age),
			Weight:                student.weight,
			Objective:             student.objective,
			TrainingFrequency:     student.trainingFrequency,
			DidBodybuilding:       student.didBodybuilding,
			MedicalCondition:      student.medicalCondition,
			PhysicalActivityLevel: student.physicalActivityLevel,
			Observations:          student.observations,
		})
		if err != nil {
			log.Printf("Failed to create student profile %s: %v", student.name, err)
			continue
		}

		fmt.Printf("Created student: %s\n", student.name)
	}

	// Create Sample Exercises
	exercises := []struct {
		name        string
		thumbnail   string
		videoURL    string
		sets        int32
		reps        int32
		restTime    *int32
		personalID  *uuid.UUID
	}{
		{
			name:      "Supino Reto",
			thumbnail: "https://example.com/supino-reto.jpg",
			videoURL:  "https://example.com/supino-reto.mp4",
			sets:      4,
			reps:      12,
			restTime:  int32Ptr(90),
		},
		{
			name:      "Agachamento Livre",
			thumbnail: "https://example.com/agachamento.jpg",
			videoURL:  "https://example.com/agachamento.mp4",
			sets:      4,
			reps:      15,
			restTime:  int32Ptr(120),
		},
		{
			name:      "Puxada Frontal",
			thumbnail: "https://example.com/puxada.jpg",
			videoURL:  "https://example.com/puxada.mp4",
			sets:      3,
			reps:      12,
			restTime:  int32Ptr(90),
		},
		{
			name:      "Desenvolvimento com Halteres",
			thumbnail: "https://example.com/desenvolvimento.jpg",
			videoURL:  "https://example.com/desenvolvimento.mp4",
			sets:      3,
			reps:      10,
			restTime:  int32Ptr(90),
		},
		{
			name:      "Rosca Direta",
			thumbnail: "https://example.com/rosca.jpg",
			videoURL:  "https://example.com/rosca.mp4",
			sets:      3,
			reps:      12,
			restTime:  int32Ptr(60),
		},
	}

	var exerciseIDs []uuid.UUID

	for _, exercise := range exercises {
		exerciseID := uuid.New()
		exerciseIDs = append(exerciseIDs, exerciseID)

		_, err := queries.CreateExercise(ctx, pgstore.CreateExerciseParams{
			ID:                  exerciseID,
			Name:                exercise.name,
			Thumbnail:           exercise.thumbnail,
			VideoURL:            exercise.videoURL,
			Sets:                exercise.sets,
			Reps:                exercise.reps,
			RestTimeBetweenSets: exercise.restTime,
			PersonalID:          exercise.personalID,
			CreatedAt:           now,
			UpdatedAt:           now,
		})
		if err != nil {
			log.Printf("Failed to create exercise %s: %v", exercise.name, err)
			continue
		}

		fmt.Printf("Created exercise: %s\n", exercise.name)
	}

	// Create Sample Workouts
	workouts := []struct {
		name        string
		description string
		thumbnail   string
		modality    string
		level       pgstore.Level
		weekDays    []pgstore.Day
		personalID  *uuid.UUID
	}{
		{
			name:        "Treino de Peito e Tríceps",
			description: "Treino focado no desenvolvimento do peitoral e tríceps",
			thumbnail:   "https://example.com/treino-peito.jpg",
			modality:    "Musculação",
			level:       pgstore.LevelIntermediary,
			weekDays:    []pgstore.Day{pgstore.DaySeg, pgstore.DayQui},
			personalID:  &personalIDs[0],
		},
		{
			name:        "Treino de Pernas",
			description: "Treino completo para membros inferiores",
			thumbnail:   "https://example.com/treino-pernas.jpg",
			modality:    "Musculação",
			level:       pgstore.LevelIntermediary,
			weekDays:    []pgstore.Day{pgstore.DayTer, pgstore.DaySex},
			personalID:  &personalIDs[1],
		},
		{
			name:        "Treino Funcional",
			description: "Treino funcional para condicionamento geral",
			thumbnail:   "https://example.com/treino-funcional.jpg",
			modality:    "Funcional",
			level:       pgstore.LevelBeginner,
			weekDays:    []pgstore.Day{pgstore.DaySeg, pgstore.DayQua, pgstore.DaySex},
			personalID:  &personalIDs[2],
		},
	}

	for _, workout := range workouts {
		workoutID := uuid.New()

		_, err := queries.CreateWorkout(ctx, pgstore.CreateWorkoutParams{
			ID:          workoutID,
			Name:        workout.name,
			Description: &workout.description,
			Thumbnail:   workout.thumbnail,
			Level:       &workout.level,
			WeekDays:    workout.weekDays,
			Exclusive:   false,
			IsTemplate:  true,
			Modality:    workout.modality,
			PersonalID:  workout.personalID,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		if err != nil {
			log.Printf("Failed to create workout %s: %v", workout.name, err)
			continue
		}

		fmt.Printf("Created workout: %s\n", workout.name)
	}

	fmt.Println("Database seeding completed successfully!")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
