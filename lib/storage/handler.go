package storage

func GetOperatorStates() error
func SaveOperatorState(pipelineID, operatorID, state, containerID, errMsg string) error 
func DeleteOperator(pipelineID, operatorID string) error