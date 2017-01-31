package engine

type PreparationInterface interface {
}

type PreparationFunc func(p PreparationInterface) error
