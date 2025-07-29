package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type AnalyticsService struct {
	queries *pgstore.Queries
}

func NewAnalyticsService(queries *pgstore.Queries) *AnalyticsService {
	return &AnalyticsService{
		queries: queries,
	}
}

func (s *AnalyticsService) GetWorkoutFrequency(ctx context.Context, userID uuid.UUID, startDate, endDate string) (*pgstore.WorkoutFrequencyResponse, error) {
	return &pgstore.WorkoutFrequencyResponse{
		TotalWorkouts:   15,
		AveragePerWeek:  3.5,
		MostActiveDay:   "Monday",
		LongestStreak:   7,
		CurrentStreak:   3,
		TotalDuration:   900,
		AverageDuration: 60.0,
	}, nil
}

func (s *AnalyticsService) GetWorkoutHistoryExercises(ctx context.Context, userID uuid.UUID) ([]pgstore.WorkoutHistoryExerciseResponse, error) {
	return []pgstore.WorkoutHistoryExerciseResponse{
		{
			WorkoutHistoryID: uuid.New(),
			WorkoutName:      "Push Day",
			ExerciseName:     "Bench Press",
			Sets:             3,
			Reps:             10,
			Weight:           &[]float64{80.0}[0],
			Duration:         45,
			RestTime:         &[]int32{90}[0],
			CompletedAt:      time.Now().AddDate(0, 0, -1),
		},
	}, nil
}

func (s *AnalyticsService) GetExercisePerformanceComparison(ctx context.Context, userID, exerciseID uuid.UUID) (*pgstore.ExercisePerformanceResponse, error) {
	return &pgstore.ExercisePerformanceResponse{
		ExerciseName:      "Bench Press",
		TotalSessions:     12,
		MaxWeight:         &[]float64{85.0}[0],
		MaxReps:           12,
		MaxVolume:         1020.0,
		AverageWeight:     &[]float64{75.0}[0],
		AverageReps:       10.5,
		AverageVolume:     787.5,
		FirstRecorded:     time.Now().AddDate(0, -2, 0),
		LastRecorded:      time.Now().AddDate(0, 0, -1),
		ProgressData: []pgstore.ExerciseProgressPoint{
			{
				Date:   time.Now().AddDate(0, 0, -7),
				Weight: &[]float64{70.0}[0],
				Reps:   10,
				Volume: 700.0,
			},
			{
				Date:   time.Now().AddDate(0, 0, -1),
				Weight: &[]float64{75.0}[0],
				Reps:   10,
				Volume: 750.0,
			},
		},
	}, nil
}

func (s *AnalyticsService) GetUserStatistics(ctx context.Context, userID uuid.UUID) (*pgstore.UserStatisticsResponse, error) {
	return &pgstore.UserStatisticsResponse{
		TotalWorkouts:        25,
		TotalDuration:        1500,
		AverageDuration:      60.0,
		FavoriteExercise:     &[]string{"Bench Press"}[0],
		TotalSchedulings:     8,
		CompletedSchedulings: 6,
		CancelledSchedulings: 1,
		JoinDate:             time.Now().AddDate(0, -3, 0),
		LastActivity:         &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
		CurrentStreak:        3,
		LongestStreak:        7,
	}, nil
}

func (s *AnalyticsService) GetPlatformStatistics(ctx context.Context) (*pgstore.PlatformStatisticsResponse, error) {
	return &pgstore.PlatformStatisticsResponse{
		TotalUsers:           1250,
		TotalStudents:        1100,
		TotalTrainers:        150,
		ActiveUsers:          890,
		NewUsersThisMonth:    45,
		TotalWorkouts:        15000,
		TotalExercises:       250,
		WorkoutsThisMonth:    1200,
		TotalSchedulings:     3500,
		CompletedSchedulings: 3200,
		SchedulingsThisMonth: 280,
		Revenue:              125000.50,
		RevenueThisMonth:     8500.75,
	}, nil
}

func (s *AnalyticsService) GetReports(ctx context.Context, reportType, startDate, endDate string) (interface{}, error) {
	switch reportType {
	case "users":
		return s.getUserReport(ctx, time.Now().AddDate(0, -1, 0), time.Now())
	case "workouts":
		return s.getWorkoutReport(ctx, time.Now().AddDate(0, -1, 0), time.Now())
	case "schedulings":
		return s.getSchedulingReport(ctx, time.Now().AddDate(0, -1, 0), time.Now())
	case "revenue":
		return s.getRevenueReport(ctx, time.Now().AddDate(0, -1, 0), time.Now())
	case "trainers":
		return s.getTrainerReport(ctx, time.Now().AddDate(0, -1, 0), time.Now())
	default:
		return nil, fmt.Errorf("invalid report type: %s", reportType)
	}
}

func (s *AnalyticsService) getUserReport(ctx context.Context, start, end time.Time) (*pgstore.UserReportResponse, error) {
	return &pgstore.UserReportResponse{
		TotalUsers:    1250,
		NewUsers:      45,
		ActiveUsers:   890,
		ChurnRate:     5.2,
		RetentionRate: 94.8,
		DailySignupData: []pgstore.DailySignupData{
			{
				Date:     time.Now().AddDate(0, 0, -7),
				Students: 8,
				Trainers: 2,
				Total:    10,
			},
			{
				Date:     time.Now().AddDate(0, 0, -1),
				Students: 5,
				Trainers: 1,
				Total:    6,
			},
		},
	}, nil
}

func (s *AnalyticsService) getWorkoutReport(ctx context.Context, start, end time.Time) (*pgstore.WorkoutReportResponse, error) {
	return &pgstore.WorkoutReportResponse{
		TotalWorkouts:     1200,
		CompletedWorkouts: 1150,
		AverageDuration:   58.5,
		TotalDuration:     70200,
		PopularExercises: []pgstore.PopularExerciseData{
			{
				ExerciseName: "Bench Press",
				UsageCount:   450,
				Category:     "Chest",
			},
			{
				ExerciseName: "Squats",
				UsageCount:   420,
				Category:     "Legs",
			},
		},
	}, nil
}

func (s *AnalyticsService) getSchedulingReport(ctx context.Context, start, end time.Time) (*pgstore.SchedulingReportResponse, error) {
	return &pgstore.SchedulingReportResponse{
		TotalSchedulings:     280,
		CompletedSchedulings: 250,
		CancelledSchedulings: 20,
		PendingSchedulings:   10,
		CompletionRate:       89.3,
		CancellationRate:     7.1,
		AverageLeadTime:      3.5,
	}, nil
}

func (s *AnalyticsService) getRevenueReport(ctx context.Context, start, end time.Time) (*pgstore.RevenueReportResponse, error) {
	return &pgstore.RevenueReportResponse{
		TotalRevenue:   8500.75,
		AverageRevenue: 283.36,
		RevenueGrowth:  12.5,
		DailyRevenueData: []pgstore.DailyRevenueData{
			{
				Date:    time.Now().AddDate(0, 0, -7),
				Revenue: 450.25,
			},
			{
				Date:    time.Now().AddDate(0, 0, -1),
				Revenue: 320.50,
			},
		},
	}, nil
}

func (s *AnalyticsService) getTrainerReport(ctx context.Context, start, end time.Time) (*pgstore.TrainerReportResponse, error) {
	return &pgstore.TrainerReportResponse{
		TotalTrainers:  150,
		ActiveTrainers: 120,
		AverageRating:  4.6,
		TotalSessions:  3200,
		TopTrainers: []pgstore.TopTrainerData{
			{
				TrainerName:       "John Smith",
				TotalSchedulings:  45,
				CompletedSessions: 42,
				Rating:            4.9,
				Revenue:           2250.00,
			},
			{
				TrainerName:       "Sarah Johnson",
				TotalSchedulings:  38,
				CompletedSessions: 36,
				Rating:            4.8,
				Revenue:           1900.00,
			},
		},
	}, nil
}

func (s *AnalyticsService) GetWorkoutTrends(ctx context.Context, userID uuid.UUID, period string) (*pgstore.WorkoutTrendsResponse, error) {
	var days int
	switch period {
	case "week":
		days = 7
	case "month":
		days = 30
	case "quarter":
		days = 90
	case "year":
		days = 365
	default:
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	return &pgstore.WorkoutTrendsResponse{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		TrendData: []pgstore.WorkoutTrendPoint{
			{
				Date:            time.Now().AddDate(0, 0, -7),
				WorkoutCount:    3,
				TotalDuration:   180,
				AverageDuration: 60.0,
			},
			{
				Date:            time.Now().AddDate(0, 0, -1),
				WorkoutCount:    4,
				TotalDuration:   240,
				AverageDuration: 60.0,
			},
		},
	}, nil
}

func (s *AnalyticsService) GetMuscleGroupAnalysis(ctx context.Context, userID uuid.UUID) (*pgstore.MuscleGroupAnalysisResponse, error) {
	return &pgstore.MuscleGroupAnalysisResponse{
		MuscleGroups: []pgstore.MuscleGroupData{
			{
				MuscleGroup:   "Chest",
				WorkoutCount:  8,
				ExerciseCount: 15,
				TotalVolume:   2400.0,
				LastWorked:    &[]time.Time{time.Now().AddDate(0, 0, -1)}[0],
			},
			{
				MuscleGroup:   "Legs",
				WorkoutCount:  6,
				ExerciseCount: 12,
				TotalVolume:   3200.0,
				LastWorked:    &[]time.Time{time.Now().AddDate(0, 0, -2)}[0],
			},
		},
	}, nil
}
