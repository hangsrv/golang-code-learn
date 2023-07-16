package main

import "fmt"

// Strategy 策略接口
type Strategy interface {
	Action()
}

// Strategy1 策略1
type Strategy1 struct{}

func (s *Strategy1) Action() {
	fmt.Println("Strategy1")
}

// Strategy2 策略2
type Strategy2 struct{}

func (s *Strategy2) Action() {
	fmt.Println("Strategy2")
}

// Strategy3 策略3
type Strategy3 struct{}

func (s *Strategy3) Action() {
	fmt.Println("Strategy3")
}

// StrategyContext 策略上下文
type StrategyContext struct {
	strategy Strategy
}

func (sc *StrategyContext) DoAction() {
	sc.strategy.Action()
}

func NewStrategyContext(strategy string) *StrategyContext {
	switch strategy {
	case "Strategy1":
		return &StrategyContext{&Strategy1{}}
	case "Strategy2":
		return &StrategyContext{&Strategy2{}}
	case "Strategy3":
		return &StrategyContext{&Strategy3{}}
	default:
		return nil
	}
}

func main() {
	sc := NewStrategyContext("Strategy2")
	sc.DoAction()
}
