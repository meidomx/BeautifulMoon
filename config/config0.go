package config

type InternalConfig struct {
	ACTION_FRAME_PER_SECOND int
	MAX_ACTION__FRAME_SKIP  int

	BASE_BLOCK_SIZE_PIXEL_CNT                           int
	MAX_RATE_OF_HIGH_SPEED_OBJECT_FIX_TO_BASE_BLOCK_FIX int //e.g. 1/n , -1 means fix to base block

	PHASE_PROCESSOR_TASK_COUNT_PER_QUEUE int
}

func NewInternalConfig() *InternalConfig {
	i := new(InternalConfig)
	i.ACTION_FRAME_PER_SECOND = 50
	i.MAX_ACTION__FRAME_SKIP = 0
	i.BASE_BLOCK_SIZE_PIXEL_CNT = 1
	i.MAX_RATE_OF_HIGH_SPEED_OBJECT_FIX_TO_BASE_BLOCK_FIX = -1

	i.PHASE_PROCESSOR_TASK_COUNT_PER_QUEUE = 128
	return i
}
