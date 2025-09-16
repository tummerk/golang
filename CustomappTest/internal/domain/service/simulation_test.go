package service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_Generator_Simulation(t *testing.T) {

	targetRTP := 0.3

	sequenceSize := 400_000

	generator, err := NewGenerator(targetRTP)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	sequence := generateRandomSequence(sequenceSize, minMultiplier, maxMultiplier)

	fmt.Println("--- Starting Simulation ---")
	fmt.Printf("Target RTP:            %.4f\n", targetRTP)
	fmt.Printf("Sequence Size:         %d\n", sequenceSize)
	fmt.Printf("Average value of seq:  %.2f\n", calculateAverage(sequence))
	fmt.Println("---------------------------")

	// --- Запуск симуляции ---
	startTime := time.Now()

	sum0 := float64(len(sequence))
	sum1 := 0.0

	for _, x := range sequence {
		multiplier := generator.Generate()

		if multiplier > x {
			sum1 += x
		}
	}

	duration := time.Since(startTime)

	initialSum := 0.0
	for _, x := range sequence {
		initialSum += x
	}

	rtpByCount := sum1 / sum0
	rtpBySum := sum1 / initialSum

	fmt.Printf("\n--- Simulation Results ---\n")
	fmt.Printf("Execution time:        %s\n", duration)
	fmt.Printf("Sum of initial values: %.2f (sum0 по P.S.)\n", initialSum)
	fmt.Printf("Sum of final values:   %.2f (sum1)\n", sum1)
	fmt.Println("---------------------------")
	fmt.Printf("RTP (calculated by count):  %.5f\n", rtpByCount)
	fmt.Printf("   -> Target RTP:           %.5f\n", targetRTP)
	fmt.Printf("   -> Absolute Error:       %+.5f\n", rtpByCount-targetRTP)
	fmt.Println("---------------------------")
	fmt.Printf("RTP (calculated by sum):    %.5f (это настоящий Return To Player, процент от вложенного)\n", rtpBySum)

	// Проверка погрешности. Например, мы ожидаем, что отклонение не будет больше 2%.
	// Тест упадет, если погрешность слишком большая.
	allowedError := 0.02 // 2%
	if abs(rtpByCount-targetRTP) > allowedError {
		t.Errorf("Error is too high! Got RTP %.5f, want %.5f (error %.5f > %.5f)",
			rtpByCount, targetRTP, abs(rtpByCount-targetRTP), allowedError)
	}
}

func generateRandomSequence(size int, min, max float64) []float64 {

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	seq := make([]float64, size)
	for i := 0; i < size; i++ {
		seq[i] = min + r.Float64()*(max-min)
	}
	return seq
}

func calculateAverage(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
