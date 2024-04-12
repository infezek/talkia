package usecase

type UseCaseReturn[Input any, Output any] interface {
	Execute(input Input) (Output, error)
}

type UseCaseNoReturn[Input any] interface {
	Execute(input Input) error
}

type UseCase interface {
	Execute() error
}
