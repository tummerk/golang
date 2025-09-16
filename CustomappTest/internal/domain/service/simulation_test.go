package service

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

// --- НОВЫЙ МОДЕРНИЗИРОВАННЫЙ ТЕСТ ---

// TestGeneratorDeviation запускает симуляцию несколько раз для разных размеров
// последовательности и вычисляет среднее, минимальное и максимальное отклонение
// фактического RTP от целевого.
func TestGeneratorDeviation(t *testing.T) {
	// --- Параметры тестирования ---
	targetRTP := 0.3
	runsPerSize := 5 // Запускаем по 5 раз для каждого размера
	sequenceSizes := []int{50_000, 100_000, 1_000_000}

	fmt.Printf("--- Запуск теста на отклонение (Target RTP: %.2f) ---\n", targetRTP)

	// Внешний цикл по заданным размерам последовательности
	for _, size := range sequenceSizes {
		fmt.Printf("\n--- Результаты для len = %d ---\n", size)

		deviations := make([]float64, 0, runsPerSize)

		// Внутренний цикл для многократного запуска симуляции
		for i := 0; i < runsPerSize; i++ {
			// Запускаем одну симуляцию
			actualRTP, err := runSingleSimulation(targetRTP, size)
			if err != nil {
				t.Fatalf("Симуляция провалена для len=%d: %v", size, err)
			}

			// Вычисляем и сохраняем абсолютное отклонение
			deviation := math.Abs(actualRTP - targetRTP)
			deviations = append(deviations, deviation)
		}

		// Рассчитываем и выводим статистику по отклонениям
		avgDeviation, minDeviation, maxDeviation := calculateDeviationStats(deviations)

		fmt.Printf("Среднее абсолютное отклонение: %.5f\n", avgDeviation)
		fmt.Printf("Мин. абсолютное отклонение:   %.5f\n", minDeviation)
		fmt.Printf("Макс. абсолютное отклонение:  %.5f\n", maxDeviation)
		fmt.Println("---------------------------------")
	}
}

// runSingleSimulation выполняет один полный прогон симуляции и возвращает
// рассчитанный RTP (по количеству).
func runSingleSimulation(targetRTP float64, sequenceSize int) (float64, error) {
	// Предполагаем, что NewGenerator, minMultiplier и maxMultiplier определены в этом пакете
	generator, err := NewGenerator(targetRTP)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать генератор: %w", err)
	}

	sequence := generateRandomSequence(sequenceSize, minMultiplier, maxMultiplier)

	sum1 := 0.0
	for _, x := range sequence {
		multiplier := generator.Generate()
		if multiplier > x {
			sum1 += x
		}
	}

	// Рассчитываем RTP по количеству (отношение выигранной суммы к количеству ставок)
	rtpByCount := sum1 / float64(len(sequence))
	return rtpByCount, nil
}

// calculateDeviationStats вычисляет среднее, минимальное и максимальное значение в срезе.
func calculateDeviationStats(data []float64) (avg, min, max float64) {
	if len(data) == 0 {
		return 0, 0, 0
	}

	sum := 0.0
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, v := range data {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	avg = sum / float64(len(data))
	return
}

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ (остались без изменений) ---

func generateRandomSequence(size int, min, max float64) []float64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	seq := make([]float64, size)
	for i := 0; i < size; i++ {
		seq[i] = min + r.Float64()*(max-min)
	}
	return seq
}

// Примечание: предполагается, что эти переменные и функция NewGenerator
// определены где-то в вашем пакете `service`.
// Если нет, раскомментируйте и адаптируйте этот блок.
/*
const (
	minMultiplier = 1.01
	maxMultiplier = 100.0
)

type Generator struct {
	targetRTP float64
}

func NewGenerator(rtp float64) (*Generator, error) {
	if rtp <= 0 {
		return nil, fmt.Errorf("RTP must be positive")
	}
	return &Generator{targetRTP: rtp}, nil
}

func (g *Generator) Generate() float64 {
	// Это заглушка, замените на вашу реальную логику генерации
	return minMultiplier + rand.Float64()*(maxMultiplier-minMultiplier)*/
