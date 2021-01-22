package main

type Application struct {
	state *State
}

func NewApplication() *Application {
	ret := &Application{
		state: NewState(),
	}

	return ret
}

func (a *Application) GetState() *State {
	return a.state
}
