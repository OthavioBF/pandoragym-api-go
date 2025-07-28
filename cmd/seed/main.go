package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
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

	fmt.Println("Starting comprehensive database seeding...")

	// Clear existing data (similar to TypeScript seed)
	fmt.Println("Clearing existing data...")
	clearExistingData(ctx, pool)

	// Hash password for all users (same as original: 123456)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	// Data arrays from the original TypeScript seed
	specializations := []string{
		"Musculação Avançada",
		"Treinamento Funcional de Alta Intensidade",
		"Hipertrofia e Força",
		"Condicionamento Físico e Performance",
		"Treinamento para Perda de Peso",
		"Treinamento Pós-Reabilitação",
		"Avaliação Física e Biomecânica",
		"Treinamento de Mobilidade e Flexibilidade",
		"Treinamento para Atletas",
		"Treinamento Preventivo de Lesões",
	}

	qualifications := []string{
		"Licenciado em Educação Física",
		"Bacharelado em Ciências do Esporte",
		"Especialização em Fisiologia do Exercício",
		"Certificação Internacional de Treinamento Funcional",
		"Curso de Biomecânica Aplicada ao Treinamento",
		"Certificação em Avaliação Física e Prescrição de Treinos",
		"Curso de Personal Trainer com Ênfase em Hipertrofia",
		"Formação em Nutrição Esportiva Avançada",
		"Especialização em Treinamento para Atletas de Alto Rendimento",
		"Curso de Reabilitação Esportiva e Prevenção de Lesões",
	}

	unsplashImages := []string{
		"https://images.unsplash.com/photo-1548690312-e3b507d8c110?q=80&w=2187&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
		"https://plus.unsplash.com/premium_photo-1672784160185-ee68f5b84e4d?q=80&w=2063&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
		"https://images.unsplash.com/photo-1606902965551-dce093cda6e7?q=80&w=2187&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
		"https://ai-previews.123rf.com/ai-txt2img/600nwm/481adba4-a4f5-44e2-98ad-7ea52b06464f.jpg",
		"https://a.storyblok.com/f/97382/2000x1500/69d267834d/personal-trainer-cover.png/m/1200x900",
		"https://assets.setmore.com/website/v2/images/industry-pages/personal-trainer/man-ipad-smiling.png",
		"https://images.pexels.com/photos/414029/pexels-photo-414029.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=800&w=800",
		"https://media.istockphoto.com/id/856797530/pt/foto/portrait-of-a-beautiful-woman-at-the-gym.jpg?s=612x612&w=0&k=20&c=Uda8Fe8FGZGha16-zF9k0Zi_qYjPXWNiZlXfgyNl_Z4=",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQN_2vNXU--sTzcl-ChVHLggoRxeyGzG2LH2A&s",
		"https://as2.ftcdn.net/v2/jpg/02/47/31/79/1000_F_247317931_ew1afYkGhLQkptDEph6x38hRwScl03oU.jpg",
		"https://images.unsplash.com/photo-1605050824853-7fb0755face3?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
	}

	names := []struct {
		first string
		last  string
	}{
		{"Thelma", "Krause"},
		{"Julie", "Smith"},
		{"Julie", "Rekamie"},
		{"Katherine", "Davis"},
		{"David", "Wilson"},
		{"Olive", "Brown"},
		{"Alvaro", "Rodrigues"},
		{"Tiago", "Alveres"},
		{"Ramon", "Dino"},
		{"Julian", "Alvarez"},
		{"Brendan", "Bartell"},
	}

	descriptions := []string{
		"Com 20 anos de experiência, meu trabalho é pautado pela seriedade e compromisso. Desenvolvo metodologias que respeitam os limites e objetivos de cada aluno, garantindo um acompanhamento individualizado e eficaz.",
		"Meu objetivo é ajudar você a alcançar sua melhor forma física de maneira segura e eficaz, promovendo uma vida ativa e saudável! Acredito que a harmonia entre corpo e mente é fundamental para resultados duradouros e significativos na sua jornada de saúde.",
		"Meu principal foco é oferecer treinos totalmente personalizados. Cada prescrição é feita de forma individualizada, assegurando que você tenha a segurança e confiança de um profissional qualificado, o que resulta em uma melhor qualidade de vida.",
		"Como Personal Trainer, busco sempre atualizar meus conhecimentos técnicos e científicos. Com base nos princípios da beneficência e da não maleficência, minha missão é contribuir para a sua motivação e foco, através de suporte psicológico durante os treinos.",
		"Dedico-me diariamente ao meu desenvolvimento profissional, aprendendo constantemente para oferecer a melhor orientação aos meus alunos. Após vivenciar experiências negativas com profissionais que não direcionavam adequadamente, desenvolvi um método de ensino que prioriza a qualidade de vida e a conquista dos objetivos.",
		"Sou especialista em Cinesiologia e Biomecânica, atuando em Vila Velha em academias, residências ou ao ar livre. Ao longo da minha carreira como Personal Trainer, já ajudei muitas pessoas a transformarem suas vidas e a alcançarem uma saúde plena. Também ofereço Consultoria Esportiva Online, mas acredito que o treinamento presencial é a forma mais assertiva de alcançar resultados.",
		"Minha história de vida inspira meus clientes. Após enfrentar problemas de saúde relacionados ao sedentarismo, tomei a decisão de mudar minha vida sem medicamentos. Iniciei caminhadas, adotei uma dieta saudável e pratiquei musculação, conseguindo zerar meu colesterol em apenas dois anos. Estou aqui para ajudar você a dar o primeiro passo rumo à sua transformação!",
		"Sou professor de Educação Física formado pela UFES, uma das universidades mais renomadas do Brasil. Com duas pós-graduações — uma em Educação e outra em Condicionamento Físico —, minha missão é levar conhecimento e excelência ao treino de cada aluno, garantindo resultados efetivos.",
		"Meu nome é Patrícia Gomes e sou apaixonada por exercícios físicos, que considero uma verdadeira terapia. A prática regular de atividades físicas não só melhora a saúde, mas também proporciona autoconhecimento e ajuda a alcançar metas pessoais de maneira positiva em todas as áreas da vida.",
		"Como professora de aulas aeróbicas, ofereço atividades como dança, ritbox, step e funcional. Venha cuidar do seu corpo e da sua saúde de forma prática, utilizando poucos equipamentos em qualquer lugar que você desejar se exercitar.",
		"Acredito no poder da alimentação saudável como aliada na prática de exercícios. Com a minha orientação, você aprenderá a escolher os melhores alimentos que potencializam seu desempenho e melhoram sua saúde, criando um estilo de vida equilibrado e sustentável.",
	}

	// Students data from the original seed
	studentsData := []struct {
		name                  string
		email                 string
		phone                 string
		avatarURL             string
		bornDate              time.Time
		age                   int32
		weight                float64
		objective             string
		trainingFrequency     string
		didBodybuilding       bool
		physicalActivityLevel string
		medicalCondition      string
		observations          string
	}{
		{
			name:                  "Aline Garcia",
			email:                 "aline@garcia.com",
			phone:                 "555-555-5551",
			avatarURL:             "https://images.unsplash.com/photo-1594381898411-846e7d193883?q=80&w=2787&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
			bornDate:              time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			age:                   29,
			weight:                75.5,
			objective:             "Lose weight",
			trainingFrequency:     "3 times a week",
			didBodybuilding:       false,
			physicalActivityLevel: "Moderate",
			medicalCondition:      "None",
			observations:          "something",
		},
		{
			name:                  "Pedro Antônio",
			email:                 "jane@smith.com",
			phone:                 "555-555-5552",
			avatarURL:             "https://plus.unsplash.com/premium_photo-1663040472837-4d2051b93735?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTU2fHxneW18ZW58MHx8MHx8fDA%3D",
			bornDate:              time.Date(1993, 5, 15, 0, 0, 0, 0, time.UTC),
			age:                   31,
			weight:                68.2,
			objective:             "Gain muscle",
			trainingFrequency:     "4 times a week",
			didBodybuilding:       true,
			physicalActivityLevel: "High",
			medicalCondition:      "Asthma",
			observations:          "something",
		},
		{
			name:                  "André Silva",
			email:                 "michael@johnson.com",
			phone:                 "555-555-5553",
			avatarURL:             "https://images.unsplash.com/photo-1639653818737-7e884dc84954?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTYyfHxneW18ZW58MHx8MHx8fDA%3D",
			bornDate:              time.Date(1990, 9, 10, 0, 0, 0, 0, time.UTC),
			age:                   34,
			weight:                82.3,
			objective:             "Maintain fitness",
			trainingFrequency:     "5 times a week",
			didBodybuilding:       true,
			physicalActivityLevel: "Very High",
			medicalCondition:      "None",
			observations:          "something",
		},
		{
			name:                  "Emily Davis",
			email:                 "emily@davis.com",
			phone:                 "555-555-5554",
			avatarURL:             "https://plus.unsplash.com/premium_photo-1705018501151-4045c97658a3?w=700&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mjl8fHdvbWFufGVufDB8fDB8fHww",
			bornDate:              time.Date(1997, 11, 25, 0, 0, 0, 0, time.UTC),
			age:                   27,
			weight:                62.5,
			objective:             "Improve endurance",
			trainingFrequency:     "3 times a week",
			didBodybuilding:       false,
			physicalActivityLevel: "Moderate",
			medicalCondition:      "Knee injury",
			observations:          "something",
		},
		{
			name:                  "Chris Brown",
			email:                 "chris@brown.com",
			phone:                 "555-555-5555",
			avatarURL:             "https://images.unsplash.com/photo-1565104781149-275a5392dabc?q=80&w=2181&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
			bornDate:              time.Date(1988, 3, 30, 0, 0, 0, 0, time.UTC),
			age:                   36,
			weight:                90.0,
			objective:             "Lose fat",
			trainingFrequency:     "2 times a week",
			didBodybuilding:       false,
			physicalActivityLevel: "Low",
			medicalCondition:      "High blood pressure",
			observations:          "something",
		},
	}

	// Create students first
	fmt.Println("Creating students...")
	var studentIDs []uuid.UUID
	for _, student := range studentsData {
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
		err = queries.CreateStudent(ctx, pgstore.CreateStudentParams{
			ID:                    userID,
			BornDate:              student.bornDate,
			Age:                   student.age,
			Weight:                student.weight,
			Objective:             student.objective,
			TrainingFrequency:     student.trainingFrequency,
			DidBodybuilding:       student.didBodybuilding,
			MedicalCondition:      &student.medicalCondition,
			PhysicalActivityLevel: &student.physicalActivityLevel,
			Observations:          &student.observations,
		})
		if err != nil {
			log.Printf("Failed to create student profile %s: %v", student.name, err)
			continue
		}

		fmt.Printf("Created student: %s\n", student.name)
	}

	// Create 11 random personal trainers (like the original seed)
	fmt.Println("Creating personal trainers...")
	var personalIDs []uuid.UUID
	
	// Copy arrays for random selection
	availableImages := make([]string, len(unsplashImages))
	copy(availableImages, unsplashImages)
	availableNames := make([]struct{ first, last string }, len(names))
	copy(availableNames, names)
	availableDescriptions := make([]string, len(descriptions))
	copy(availableDescriptions, descriptions)

	for i := 0; i <= 10; i++ {
		userID := uuid.New()
		personalIDs = append(personalIDs, userID)

		// Get random data and remove from arrays
		randomName := getRandomNameAndRemove(&availableNames)
		randomDescription := getRandomAndRemove(&availableDescriptions)
		randomSpecialization := specializations[rand.Intn(len(specializations))]
		randomQualification := qualifications[rand.Intn(len(qualifications))]
		randomExperience := fmt.Sprintf("%d anos de experiência em Educação Física", 2+rand.Intn(14))

		// Create user
		_, err := queries.CreateUser(ctx, pgstore.CreateUserParams{
			ID:        userID,
			Name:      randomName,
			Email:     fmt.Sprintf("trainer%d@pandoragym.com", i),
			Phone:     fmt.Sprintf("11%09d", rand.Intn(1000000000)),
			Password:  string(passwordHash),
			Role:      pgstore.RolePersonal,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			log.Printf("Failed to create personal trainer user: %v", err)
			continue
		}

		// Create personal profile
		err = queries.CreatePersonal(ctx, pgstore.CreatePersonalParams{
			ID:             userID,
			Description:    &randomDescription,
			Experience:     &randomExperience,
			Specialization: &randomSpecialization,
			Qualifications: &randomQualification,
		})
		if err != nil {
			log.Printf("Failed to create personal trainer profile: %v", err)
			continue
		}

		fmt.Printf("Created personal trainer: %s\n", randomName)
	}

	// Create featured personal trainer (like Bianca Andrade in the original)
	fmt.Println("Creating featured personal trainer...")
	featuredUserID := uuid.New()
	personalIDs = append(personalIDs, featuredUserID)

	_, err = queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        featuredUserID,
		Name:      "Bianca Andrade",
		Email:     "bianca@pandoragym.com",
		Phone:     "11987654321",
		Password:  string(passwordHash),
		Role:      pgstore.RolePersonal,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Printf("Failed to create featured personal trainer user: %v", err)
	} else {
		err = queries.CreatePersonal(ctx, pgstore.CreatePersonalParams{
			ID:             featuredUserID,
			Description:    stringPtr("Como Personal Trainer, busco sempre atualizar meus conhecimentos técnicos e científicos. Com base nos princípios da beneficência e da não maleficência, minha missão é contribuir para a sua motivação e foco, através de suporte psicológico durante os treinos."),
			Experience:     stringPtr("7 anos de experiência na área de Educação Física"),
			Specialization: stringPtr("Treinamento de atletas"),
			Qualifications: stringPtr("Curso de Personal com foco em Hipertrofia"),
		})
		if err != nil {
			log.Printf("Failed to create featured personal trainer profile: %v", err)
		} else {
			fmt.Println("Created featured personal trainer: Bianca Andrade")
		}
	}

	// Create sample workouts for each personal trainer (like the original seed)
	fmt.Println("Creating sample workouts...")
	
	workoutNames := []string{
		"Treino de Superiores", "Treino de Quadriceps", "Treino de Glúteo",
		"Treino Funcional", "Treino de Peito", "Treino de Costas",
		"Treino de Braços", "Treino de Pernas", "Treino HIIT",
		"Treino de Core", "Treino Full Body",
	}
	
	workoutDescriptions := []string{
		"Treino focado no desenvolvimento dos membros superiores",
		"Treino intensivo para quadriceps e força nas pernas",
		"Treino especializado para glúteos e posterior de coxa",
		"Treinamento funcional para condicionamento geral",
		"Desenvolvimento do peitoral e músculos auxiliares",
		"Fortalecimento das costas e melhora da postura",
		"Treino completo para braços e antebraços",
		"Treino completo para membros inferiores",
		"Treino intervalado de alta intensidade",
		"Fortalecimento do core e estabilização",
		"Treino completo para todo o corpo",
	}
	
	workoutThumbnails := []string{
		"https://plenocorpo.com/wp-content/uploads/2023/07/treino_membros-superiores-feminino.jpg",
		"https://blog.gsuplementos.com.br/wp-content/uploads/2020/10/iStock-1345539108.jpg",
		"https://ciaathleticasjc.com.br/wp-content/uploads/2023/12/Cia-Athletica-SJC-Exercicios-para-levantar-gluteos-rapidamente-Autores-Grupo-S2-Marketing-Freepik.jpg",
		"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1581009146145-b5ef050c2e1e?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1581009146145-b5ef050c2e1e?q=80&w=2070&auto=format&fit=crop",
		"https://blog.gsuplementos.com.br/wp-content/uploads/2020/10/iStock-1345539108.jpg",
		"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1581009146145-b5ef050c2e1e?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?q=80&w=2070&auto=format&fit=crop",
	}

	levels := []pgstore.Level{pgstore.LevelBeginner, pgstore.LevelIntermediary, pgstore.LevelAdvanced}
	weekDaysOptions := [][]pgstore.Day{
		{pgstore.DaySeg},
		{pgstore.DayTer},
		{pgstore.DayQua},
		{pgstore.DayQui},
		{pgstore.DaySex},
		{pgstore.DaySeg, pgstore.DayQua},
		{pgstore.DayTer, pgstore.DayQui},
		{pgstore.DaySeg, pgstore.DayQua, pgstore.DaySex},
	}

	// Create 5 workouts for each personal trainer
	for i, personalID := range personalIDs {
		for j := 0; j < 5; j++ {
			workoutID := uuid.New()
			nameIndex := (i*5 + j) % len(workoutNames)
			selectedWeekDays := weekDaysOptions[rand.Intn(len(weekDaysOptions))]
			
			// Create workout using direct SQL to avoid the Day array scanning issue
			_, err := pool.Exec(ctx, `
				INSERT INTO workout (id, name, description, thumbnail, level, week_days, exclusive, is_template, modality, personal_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			`,
				workoutID,
				workoutNames[nameIndex],
				workoutDescriptions[nameIndex],
				workoutThumbnails[nameIndex],
				levels[rand.Intn(len(levels))],
				pq.Array(convertDaysToStringArray(selectedWeekDays)),
				rand.Float32() < 0.3, // 30% chance of being exclusive
				rand.Float32() < 0.7, // 70% chance of being template
				"Musculação",
				personalID,
				now,
				now,
			)
			if err != nil {
				log.Printf("Failed to create workout: %v", err)
				continue
			}
		}
		fmt.Printf("Created 5 workouts for personal trainer %d\n", i+1)
	}

	// Create some admin workouts (without personal trainer)
	fmt.Println("Creating admin workouts...")
	for j := 0; j < 5; j++ {
		workoutID := uuid.New()
		nameIndex := j % len(workoutNames)
		selectedWeekDays := weekDaysOptions[rand.Intn(len(weekDaysOptions))]
		
		_, err := pool.Exec(ctx, `
			INSERT INTO workout (id, name, description, thumbnail, level, week_days, exclusive, is_template, modality, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
			workoutID,
			fmt.Sprintf("Admin %s", workoutNames[nameIndex]),
			workoutDescriptions[nameIndex],
			workoutThumbnails[nameIndex],
			levels[rand.Intn(len(levels))],
			pq.Array(convertDaysToStringArray(selectedWeekDays)),
			false,
			false,
			"Musculação",
			now,
			now,
		)
		if err != nil {
			log.Printf("Failed to create admin workout: %v", err)
			continue
		}
	}
	fmt.Println("Created 5 admin workouts")

	// Create sample exercises for some workouts (like the original seed)
	fmt.Println("Creating sample exercises...")
	createSampleExercises(ctx, pool, now)

	// Create specific workouts from the original TypeScript seed
	fmt.Println("Creating specific workouts from original seed...")
	createSpecificWorkouts(ctx, pool, studentIDs, personalIDs, now)

	fmt.Println("Database seeding completed successfully!")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func clearExistingData(ctx context.Context, pool *pgxpool.Pool) {
	// Clear data in dependency order (similar to TypeScript seed)
	fmt.Println("Clearing workouts_history...")
	_, err := pool.Exec(ctx, "DELETE FROM workouts_history")
	if err != nil {
		log.Printf("Failed to clear workouts_history: %v", err)
	}
	
	fmt.Println("Clearing exercises_setup...")
	_, err = pool.Exec(ctx, "DELETE FROM exercises_setup")
	if err != nil {
		log.Printf("Failed to clear exercises_setup: %v", err)
	}
	
	fmt.Println("Clearing workout...")
	_, err = pool.Exec(ctx, "DELETE FROM workout")
	if err != nil {
		log.Printf("Failed to clear workout: %v", err)
	}
	
	fmt.Println("Clearing personal...")
	_, err = pool.Exec(ctx, "DELETE FROM personal")
	if err != nil {
		log.Printf("Failed to clear personal: %v", err)
	}
	
	fmt.Println("Clearing student...")
	_, err = pool.Exec(ctx, "DELETE FROM student")
	if err != nil {
		log.Printf("Failed to clear student: %v", err)
	}
	
	fmt.Println("Clearing users...")
	_, err = pool.Exec(ctx, "DELETE FROM users")
	if err != nil {
		log.Printf("Failed to clear users: %v", err)
	}
}

func getRandomAndRemove(slice *[]string) string {
	if len(*slice) == 0 {
		return ""
	}
	index := rand.Intn(len(*slice))
	item := (*slice)[index]
	*slice = append((*slice)[:index], (*slice)[index+1:]...)
	return item
}

func getRandomNameAndRemove(names *[]struct{ first, last string }) string {
	if len(*names) == 0 {
		return "Unknown User"
	}
	index := rand.Intn(len(*names))
	name := (*names)[index]
	*names = append((*names)[:index], (*names)[index+1:]...)
	return fmt.Sprintf("%s %s", name.first, name.last)
}

func convertDaysToStringArray(days []pgstore.Day) []string {
	result := make([]string, len(days))
	for i, day := range days {
		result[i] = string(day)
	}
	return result
}

func createSampleExercises(ctx context.Context, pool *pgxpool.Pool, now time.Time) {
	// Get some workout IDs to add exercises to
	rows, err := pool.Query(ctx, "SELECT id FROM workout LIMIT 10")
	if err != nil {
		log.Printf("Failed to get workout IDs: %v", err)
		return
	}
	defer rows.Close()

	var workoutIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			log.Printf("Failed to scan workout ID: %v", err)
			continue
		}
		workoutIDs = append(workoutIDs, id)
	}

	// Sample exercises from the original TypeScript seed
	exercises := []struct {
		title     string
		video     string
		thumbnail string
		sets      int32
		reps      int32
		restTime  int32
		load      int32
	}{
		{
			title:     "Push-Ups",
			video:     "https://www.youtube.com/watch?v=_l3ySVKYVJ8&t=1s&ab_channel=CrossFit",
			thumbnail: "https://www.fitnesseducation.edu.au/wp-content/uploads/2017/03/Pushups.jpg",
			sets:      3,
			reps:      20,
			restTime:  30,
			load:      0,
		},
		{
			title:     "Knee Push Ups",
			video:     "https://www.youtube.com/watch?app=desktop&v=jWxvty2KROs",
			thumbnail: "https://www.shutterstock.com/image-photo/young-attractive-woman-doing-warming-260nw-610195001.jpg",
			sets:      3,
			reps:      15,
			restTime:  30,
			load:      0,
		},
		{
			title:     "Rosca martelo",
			video:     "https://www.youtube.com/watch?v=RUZ7Med2Mg4&ab_channel=CloudGym",
			thumbnail: "https://vitat.com.br/wp-content/uploads/2023/08/rosca-martelo.jpg",
			sets:      3,
			reps:      15,
			restTime:  30,
			load:      0,
		},
		{
			title:     "Elevação frontal",
			video:     "https://www.youtube.com/watch?v=jhxLYSm_P-k&ab_channel=MyTrainingPRO",
			thumbnail: "https://i.ytimg.com/vi/RUzWF4aDt7g/maxresdefault.jpg",
			sets:      3,
			reps:      15,
			restTime:  30,
			load:      0,
		},
		{
			title:     "Agachamento Livre",
			video:     "https://www.youtube.com/watch?v=86ZW7tmmLuU&t=1s&ab_channel=SuporteSa%C3%BAde",
			thumbnail: "https://www.cnnbrasil.com.br/wp-content/uploads/sites/12/2023/09/tipos-de-agachamento-barra.jpg?w=1024",
			sets:      4,
			reps:      30,
			restTime:  30,
			load:      20,
		},
		{
			title:     "Leg Press 45°",
			video:     "https://www.youtube.com/watch?v=eg2R_x8uMxQ",
			thumbnail: "https://treinomestre.com.br/wp-content/uploads/2019/06/leg-press-45.jpg",
			sets:      4,
			reps:      20,
			restTime:  30,
			load:      50,
		},
		{
			title:     "Cadeira extensora",
			video:     "https://www.youtube.com/watch?v=XZA5xiAs8Es&ab_channel=BOAFORMA",
			thumbnail: "https://www.lmesportes.com.br/media/magefan_blog/cadeira_extensora.jpg",
			sets:      4,
			reps:      20,
			restTime:  30,
			load:      45,
		},
		{
			title:     "Elevação Pélvica",
			video:     "https://www.youtube.com/watch?v=E7iIqI8ldJ0&ab_channel=MarianaSardelli",
			thumbnail: "https://image.tuasaude.com/media/article/co/sw/elevacao-pelvica_63320_l.jpg",
			sets:      3,
			reps:      15,
			restTime:  45,
			load:      20,
		},
	}

	// Add 3-5 exercises to each workout
	for _, workoutID := range workoutIDs {
		numExercises := 3 + rand.Intn(3) // 3-5 exercises
		for i := 0; i < numExercises; i++ {
			exercise := exercises[rand.Intn(len(exercises))]
			exerciseID := uuid.New()
			
			_, err := pool.Exec(ctx, `
				INSERT INTO exercises_setup (id, name, video_url, thumbnail, sets, reps, rest_time_between_sets, load, workout_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`,
				exerciseID,
				exercise.title,
				exercise.video,
				exercise.thumbnail,
				exercise.sets,
				exercise.reps,
				exercise.restTime,
				exercise.load,
				workoutID,
			)
			if err != nil {
				log.Printf("Failed to create exercise: %v", err)
				continue
			}
		}
		fmt.Printf("Added %d exercises to workout\n", numExercises)
	}
}

func createSpecificWorkouts(ctx context.Context, pool *pgxpool.Pool, studentIDs, personalIDs []uuid.UUID, now time.Time) {
	// Create "Treino de Superiores" for a student (like the original)
	if len(studentIDs) > 0 {
		workoutID := uuid.New()
		_, err := pool.Exec(ctx, `
			INSERT INTO workout (id, name, description, thumbnail, level, week_days, exclusive, is_template, modality, student_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			workoutID,
			"Treino de Superiores",
			"A simple beginner HIIT workout to get started",
			"https://plenocorpo.com/wp-content/uploads/2023/07/treino_membros-superiores-feminino.jpg",
			"BEGINNER",
			pq.Array([]string{"Seg"}),
			false,
			false,
			"Musculação",
			studentIDs[0], // Assign to first student
			now,
			now,
		)
		if err != nil {
			log.Printf("Failed to create specific workout: %v", err)
		} else {
			// Add specific exercises from the original seed
			exercises := []struct {
				title     string
				video     string
				thumbnail string
				sets      int32
				reps      int32
				restTime  int32
				load      int32
			}{
				{
					title:     "Push-Ups",
					video:     "https://www.youtube.com/watch?v=_l3ySVKYVJ8&t=1s&ab_channel=CrossFit",
					thumbnail: "https://www.fitnesseducation.edu.au/wp-content/uploads/2017/03/Pushups.jpg",
					sets:      3,
					reps:      20,
					restTime:  30,
					load:      0,
				},
				{
					title:     "Knee Push Ups",
					video:     "https://www.youtube.com/watch?app=desktop&v=jWxvty2KROs",
					thumbnail: "https://www.shutterstock.com/image-photo/young-attractive-woman-doing-warming-260nw-610195001.jpg",
					sets:      3,
					reps:      15,
					restTime:  30,
					load:      0,
				},
				{
					title:     "Rosca martelo",
					video:     "https://www.youtube.com/watch?v=RUZ7Med2Mg4&ab_channel=CloudGym",
					thumbnail: "https://vitat.com.br/wp-content/uploads/2023/08/rosca-martelo.jpg",
					sets:      3,
					reps:      15,
					restTime:  30,
					load:      0,
				},
				{
					title:     "Elevação frontal",
					video:     "https://www.youtube.com/watch?v=jhxLYSm_P-k&ab_channel=MyTrainingPRO",
					thumbnail: "https://i.ytimg.com/vi/RUzWF4aDt7g/maxresdefault.jpg",
					sets:      3,
					reps:      15,
					restTime:  30,
					load:      0,
				},
			}

			for _, exercise := range exercises {
				exerciseID := uuid.New()
				_, err := pool.Exec(ctx, `
					INSERT INTO exercises_setup (id, name, video_url, thumbnail, sets, reps, rest_time_between_sets, load, workout_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				`,
					exerciseID,
					exercise.title,
					exercise.video,
					exercise.thumbnail,
					exercise.sets,
					exercise.reps,
					exercise.restTime,
					exercise.load,
					workoutID,
				)
				if err != nil {
					log.Printf("Failed to create specific exercise: %v", err)
				}
			}
			fmt.Println("Created specific workout: Treino de Superiores")
		}
	}

	// Create "Treino de Quadriceps" for a personal trainer and student (like the original)
	if len(personalIDs) > 0 && len(studentIDs) > 0 {
		workoutID := uuid.New()
		_, err := pool.Exec(ctx, `
			INSERT INTO workout (id, name, description, thumbnail, level, week_days, exclusive, is_template, modality, personal_id, student_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`,
			workoutID,
			"Treino de Quadriceps",
			"Treino Quadriceps para quem já tem alguma experiência.",
			"https://blog.gsuplementos.com.br/wp-content/uploads/2020/10/iStock-1345539108.jpg",
			"BEGINNER",
			pq.Array([]string{"Ter"}),
			true,  // exclusive
			false, // not template
			"Musculação",
			personalIDs[len(personalIDs)-1], // Use featured trainer (last one added)
			studentIDs[0],                   // Assign to first student
			now,
			now,
		)
		if err != nil {
			log.Printf("Failed to create quadriceps workout: %v", err)
		} else {
			// Add specific quadriceps exercises
			exercises := []struct {
				title     string
				video     string
				thumbnail string
				sets      int32
				reps      int32
				restTime  int32
				load      int32
			}{
				{
					title:     "Agachamento Livre",
					video:     "https://www.youtube.com/watch?v=86ZW7tmmLuU&t=1s&ab_channel=SuporteSa%C3%BAde",
					thumbnail: "https://www.cnnbrasil.com.br/wp-content/uploads/sites/12/2023/09/tipos-de-agachamento-barra.jpg?w=1024",
					sets:      4,
					reps:      30,
					restTime:  30,
					load:      20,
				},
				{
					title:     "Leg Press 45°",
					video:     "https://www.youtube.com/watch?v=eg2R_x8uMxQ",
					thumbnail: "https://treinomestre.com.br/wp-content/uploads/2019/06/leg-press-45.jpg",
					sets:      4,
					reps:      20,
					restTime:  30,
					load:      50,
				},
				{
					title:     "Cadeira extensora",
					video:     "https://www.youtube.com/watch?v=XZA5xiAs8Es&ab_channel=BOAFORMA",
					thumbnail: "https://www.lmesportes.com.br/media/magefan_blog/cadeira_extensora.jpg",
					sets:      4,
					reps:      20,
					restTime:  30,
					load:      45,
				},
				{
					title:     "Afundo com halteres ou barra",
					video:     "https://www.youtube.com/watch?v=ALP9JIXA-PA&ab_channel=MyTrainingPRO",
					thumbnail: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRWrrenBkz3yU7a2q3LqOk8Pu-iDeSOPR37W4Vd2heq2WgH8373Wy3MrGjU8HUSRCBOYho&usqp=CAU",
					sets:      4,
					reps:      20,
					restTime:  30,
					load:      10,
				},
				{
					title:     "Agachamento Búlgaro",
					video:     "https://www.youtube.com/watch?v=txBXA1dvlAQ&ab_channel=RodrigoZagoTreinador",
					thumbnail: "https://static1.minhavida.com.br/articles/12/85/89/20/homem-fazendo-agachamento-bulgaro-amp_card-1.jpg",
					sets:      4,
					reps:      20,
					restTime:  30,
					load:      10,
				},
			}

			for _, exercise := range exercises {
				exerciseID := uuid.New()
				_, err := pool.Exec(ctx, `
					INSERT INTO exercises_setup (id, name, video_url, thumbnail, sets, reps, rest_time_between_sets, load, workout_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				`,
					exerciseID,
					exercise.title,
					exercise.video,
					exercise.thumbnail,
					exercise.sets,
					exercise.reps,
					exercise.restTime,
					exercise.load,
					workoutID,
				)
				if err != nil {
					log.Printf("Failed to create quadriceps exercise: %v", err)
				}
			}
			fmt.Println("Created specific workout: Treino de Quadriceps")
		}
	}
}
