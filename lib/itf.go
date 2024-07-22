package lib

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

type StorageHandler interface {
	GetOperatorStates() ([]operator.OperatorState, error)
	SaveOperatorState(pipelineID, operatorID, state, containerID, errMsg string) error 
	DeleteOperator(pipelineID, operatorID string) error
}