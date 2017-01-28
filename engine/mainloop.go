package engine

import (
	"runtime"
	"time"
)

type engine struct {
	frequency         int
	mainLoopTimer     *time.Ticker
	mainLoopTimerChan <-chan time.Time
	stopped           bool


	preparation PreparationFunc
}

type PreparationFunc func(p PreparationInterface) error

func emptyPreaparation(p PreparationInterface) error {
	return nil
}

func NewAndStartMainLoop(f PreparationFunc) (*engine, error) {
	e := new(engine)
	e.frequency = 25
	e.mainLoopTimer = time.NewTicker(1000 / 25 * time.Millisecond)
	e.mainLoopTimerChan = e.mainLoopTimer.C
	if f == nil {
		f = emptyPreaparation
	}
	e.preparation = f
	err := f(e)
	if err != nil {
		return nil, err
	}
	go e.mainLoop()
	return e, nil
}

func (this *engine) mainLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	timerChan := this.mainLoopTimerChan
	for !this.stopped {
		<-timerChan
		//TODO
	}
}

func (this *engine) DoNothing() {
}
