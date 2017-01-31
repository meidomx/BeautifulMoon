package api

type PhaseTask struct {
	dispatchId int
	Attachment interface{}
}

func (this PhaseTask) HashCodeInt() int {
	return this.dispatchId
}

type PhaseProcessor interface {
	ProcessInitPhase()
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
	DoInitPhase(PhaseController)
	DoConcurrentPhase(t PhaseTask, processorNo int)
	DoFinalPhase(PhaseController)
}

type PhaseInitProcessFunction func(PhaseController)
type PhaseConcurrentProcessFunction func(t PhaseTask, processorNo int)
type PhaseFinalProcessFunction func(PhaseController)
