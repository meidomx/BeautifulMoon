package phase

import (
	"runtime"
	"sync"

	"github.com/meidomx/BeautifulMoon/config"
	"github.com/meidomx/BeautifulMoon/engine/api"
)

var _ api.PhaseProcessor = new(phaseProcessor)
var _ api.PhaseController = new(phaseProcessor)

type _PhaseType int

const (
	_PHASE_INIT _PhaseType = iota
	_PHASE_CONCURRENT
	_PHASE_FINAL
	_PHASE_BACKGROUND

	_MAX_THREAD_CNT = 8
)

func NewAndStartPhaseProcessor(c *config.InternalConfig, ph api.PhaseHandler) *phaseProcessor {
	processor := new(phaseProcessor)
	processor.closeChan = make(chan bool, 1)
	processor.taskChan = make([]chan api.PhaseTask, _MAX_THREAD_CNT)
	for i := 0; i < len(processor.taskChan); i++ {
		processor.taskChan[i] = make(chan api.PhaseTask, c.PHASE_PROCESSOR_TASK_COUNT_PER_QUEUE)
	}
	processor.threadCnt = _MAX_THREAD_CNT
	if runtime.NumCPU() < _MAX_THREAD_CNT {
		processor.threadCnt = runtime.NumCPU()
	}
	processor.initPhaseFunction = ph.DoInitPhase
	processor.concurrentPhaseFunction = ph.DoConcurrentPhase
	processor.finalPhaseFunction = ph.DoFinalPhase
	processor.phaseHandler = ph
	processor.closeWaitGroup = new(sync.WaitGroup)
	for i := 0; i < processor.threadCnt; i++ {
		go processor.concurrentProcessingLoop(i)
	}

	return processor
}

type phaseProcessor struct {
	closeChan      chan bool
	closeWaitGroup *sync.WaitGroup
	taskChan       []chan api.PhaseTask
	threadCnt      int

	phaseHandler            api.PhaseHandler
	initPhaseFunction       api.PhaseInitProcessFunction
	concurrentPhaseFunction api.PhaseConcurrentProcessFunction
	finalPhaseFunction      api.PhaseFinalProcessFunction

	currentPhase _PhaseType
}

func (this *phaseProcessor) SubmitTask(task api.PhaseTask) {
	this.taskChan[task.HashCodeInt()%_MAX_THREAD_CNT] <- task
}

func (this *phaseProcessor) BatchSubmitTask(tasks []api.PhaseTask) {
	for _, v := range tasks {
		this.SubmitTask(v)
	}
}

// init phase
// may be blocked when: task channel is full
func (this *phaseProcessor) ProcessInitPhase(p *api.LoopTriggeredEvent) {
	this.currentPhase = _PHASE_INIT
	this.initPhaseFunction(this, p)
}

// concurrent phase
// blocked until: task channel is empty & no task run
func (this *phaseProcessor) ProcessConcurrentPhase() {
	//TODO
}

// final phase
// may be blocked when: task channel is full
func (this *phaseProcessor) ProcessFinalPhase() {
	this.currentPhase = _PHASE_FINAL
	this.finalPhaseFunction(this)
}

// background phase
// nonblocking
func (this *phaseProcessor) ProcessBackgroundPhase() {
	// nothing
}

// blocking
func (this *phaseProcessor) Shutdown() {
	close(this.closeChan)
	this.closeWaitGroup.Wait()
}

func (this *phaseProcessor) concurrentProcessingLoop(fromIdx int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	this.closeWaitGroup.Add(1)
	defer this.closeWaitGroup.Done()

	no := fromIdx
	ch1 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch2 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch3 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch4 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch5 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch6 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch7 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	ch8 := this.taskChan[fromIdx%_MAX_THREAD_CNT]
	fromIdx++
	closeCh := this.closeChan
	concurrentFunction := this.concurrentPhaseFunction
	var t api.PhaseTask

LOOP:
	for true {
		select {
		case t = <-ch1:
			concurrentFunction(t, no)
		case t = <-ch2:
			concurrentFunction(t, no)
		case t = <-ch3:
			concurrentFunction(t, no)
		case t = <-ch4:
			concurrentFunction(t, no)
		case t = <-ch5:
			concurrentFunction(t, no)
		case t = <-ch6:
			concurrentFunction(t, no)
		case t = <-ch7:
			concurrentFunction(t, no)
		case t = <-ch8:
			concurrentFunction(t, no)
		case <-closeCh:
			break LOOP
		}
	}
}
