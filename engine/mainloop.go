package engine

import (
	"runtime"
	"time"

	"github.com/meidomx/BeautifulMoon/config"
	"github.com/meidomx/BeautifulMoon/engine/api"
	"github.com/meidomx/BeautifulMoon/engine/internal/phase"
)

//type check
var _ PreparationInterface = new(engine)

type engine struct {
	frequency         int
	mainLoopTimer     *time.Ticker
	mainLoopTimerChan <-chan time.Time
	stopped           bool

	preparation PreparationFunc

	phaseProcess api.PhaseProcessor

	fps float32
}

func NewAndStartMainLoop(f PreparationFunc, c *config.InternalConfig, ph api.PhaseHandler) (*engine, error) {
	e := new(engine)
	e.frequency = c.ACTION_FRAME_PER_SECOND
	e.mainLoopTimer = time.NewTicker(1000 / time.Duration(e.frequency) * time.Millisecond)
	e.mainLoopTimerChan = e.mainLoopTimer.C
	if f == nil {
		f = func(p PreparationInterface) error { return nil }
	}
	e.preparation = f
	err := f(e)
	if err != nil {
		return nil, err
	}

	e.phaseProcess = phase.NewAndStartPhaseProcessor(c, ph)

	go e.mainLoop()
	return e, nil
}

func (this *engine) mainLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	last := time.Now().UnixNano() / 1000 / 1000
	timerChan := this.mainLoopTimerChan
	phaseProcessor := this.phaseProcess

	event := new(api.LoopTriggeredEvent)
	event.LastTriggeredTime = <-timerChan

	for !this.stopped {
		invokeTime := <-timerChan

		event.TriggeredTime = invokeTime

		phaseProcessor.ProcessInitPhase(event)
		phaseProcessor.ProcessConcurrentPhase()
		phaseProcessor.ProcessFinalPhase()
		phaseProcessor.ProcessBackgroundPhase()

		curTime := time.Now().UnixNano() / 1000 / 1000
		this.fps = 1000 / float32(curTime-last)
		last = curTime
		event.LastTriggeredTime = event.TriggeredTime
	}
	phaseProcessor.Shutdown()
	this.mainLoopTimer.Stop()
}

func (this *engine) DoNothing() {
}

func (this *engine) GetFPS() float32 {
	return this.fps
}
