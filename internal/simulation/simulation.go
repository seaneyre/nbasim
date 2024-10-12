package simulation

import "fmt"

type Simulation struct {
	...
}

func New() *Simulation {
	return &Simulation{}
}

func (s *Simulation) Run() error {
	fmt.Println("Running simulation")
	...
	return nil
}
