package lib

import (
	"context"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	"database/sql/driver"
)

type StorageHandler interface {
	GetOperatorStates(ctx context.Context) ([]agent.OperatorState, error)
	SaveOperatorState(ctx context.Context, pipelineID, operatorID, state, containerID, errMsg string, txItf driver.Tx) error 
	DeleteOperator(ctx context.Context, pipelineID, operatorID string, txItf driver.Tx) error
}