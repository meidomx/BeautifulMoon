package api

import "time"

type PhaseTask struct {
	dispatchId int
	Attachment interface{}
}

func (this PhaseTask) HashCodeInt() int {
	return this.dispatchId
}

type PhaseProcessor interface {
	ProcessInitPhase(*LoopTriggeredEvent)
	ProcessConcurrentPhase()
	ProcessFinalPhase()
	ProcessBackgroundPhase()
	Shutdown() //blocking
}

type PhaseController interface {
	SubmitTask(task PhaseTask)
	BatchSubmitTask(tasks []PhaseTask)
}

type PhaseHandler interface {
	DoInitPhase(PhaseController, *LoopTriggeredEvent)
	DoConcurrentPhase(t PhaseTask, processorNo int)
	DoFinalPhase(PhaseController)
}

type PhaseInitProcessFunction func(PhaseController, *LoopTriggeredEvent)
type PhaseConcurrentProcessFunction func(t PhaseTask, processorNo int)
type PhaseFinalProcessFunction func(PhaseController)

type LoopTriggeredEvent struct {
	TriggeredTime     time.Time
	LastTriggeredTime time.Time
}
