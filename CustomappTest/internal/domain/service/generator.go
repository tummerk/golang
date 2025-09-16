package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	minMultiplier float64 = 1
	maxMultiplier float64 = 10000
)

type Generator struct {
	rtp float64

	mu   sync.Mutex
	rand *rand.Rand
}

func NewGenerator(rtp float64) (*Generator, error) {
	if rtp <= 0 || rtp > 1.0 {
		return nil, fmt.Errorf("rtp must be in (0, 1.0], but got %f", rtp)
	}

	return &Generator{
		rtp:  rtp,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (g *Generator) Generate() float64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Генерируем случайное число u в [0, 1)
	u := g.rand.Float64()
	//ЗАДАЧА ГЛАВНАЯ СДЕЛАТЬ ТАК ЧТОБЫ S(x)=P(M>x)=rtp/x функция выживания; F(x)=P(M<=x)=1-S(x)
	//чтобы это сделать подставим x=1 S(x)=rtp => x будет умирать c шансом F(1)=1-rtp
	if u < (1.0 - g.rtp) {
		return 1.0
	}
	// взяв производную из функции F(x) мы найдем саму функцию плотности  f(x) = 1/x^2.

	// растянем [1-rtp, 1) до [0, 1).
	v := (u - (1.0 - g.rtp)) / g.rtp

	// а вот уже от функции плотности берём интеграл от 1 до x получаем F(x)=1-1/x а отсюда
	// находим обратную функцию x по F(x) по Методу обратного преобразования и получим следующее

	multiplier := 1.0 / (1.0 - v) // от 1 до +беск генерируем число

	if multiplier > maxMultiplier {
		return maxMultiplier
	}

	return multiplier
}
