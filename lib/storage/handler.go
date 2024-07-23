package storage

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

var tLayout = time.RFC3339Nano

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}


func (h *Handler) GetOperatorStates(ctx context.Context) ([]agent.OperatorState, error) {
	query := "SELECT pipeline_id, operator_id, state, container_id, error FROM operator_states"
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err 
	}
	defer rows.Close()
	allOperatorStates := []agent.OperatorState{}
	for rows.Next() {
		var operatorState agent.OperatorState
		if err = rows.Scan(&operatorState.PipelineID, &operatorState.OperatorID, &operatorState.State, &operatorState.ContainerID, &operatorState.ErrMsg); err != nil {
			return nil, err
		}
		allOperatorStates = append(allOperatorStates, operatorState)
	}
	return allOperatorStates, nil
}

func (h *Handler) SaveOperatorState(ctx context.Context, pipelineID, operatorID, state, containerID, errMsg string, txItf driver.Tx) error {
	tx, err := h.BeginTransaction(ctx, txItf)
	if err != nil {
		return err
	}

	rows, err := h.db.Query("SELECT COUNT(*) FROM operator_states  WHERE pipeline_id == ? AND operator_id == ?", pipelineID, operatorID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int

	for rows.Next() {   
		if err := rows.Scan(&count); err != nil {
			return err
		}
	}
	rows.Close()
	if count == 0 {
		_, err = tx.ExecContext(ctx, "INSERT INTO operator_states (pipeline_id, operator_id, state, container_id, error) VALUES (?, ?, ?, ?, ?);", pipelineID, operatorID, state, containerID, errMsg)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, "UPDATE operator_states SET state=?, container_id=?, error=? WHERE pipeline_id=? AND operator_id=?", state, containerID, errMsg, pipelineID, operatorID)
		if err != nil {
			return err
		}
	}

	return h.Commit(tx, txItf)
}

func (h *Handler) DeleteOperator(ctx context.Context, pipelineID, operatorID string, txItf driver.Tx) error {
	tx, err := h.BeginTransaction(ctx, txItf)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM operator_states WHERE pipeline_id == ? AND operator_id == ?", pipelineID, operatorID)
	if err != nil {
		return err
	}

	return h.Commit(tx, txItf)
}


func (h *Handler) BeginTransaction(ctx context.Context, txItf driver.Tx) (*sql.Tx, error) {
	if txItf != nil {
		tx := txItf.(*sql.Tx)
		return tx, nil
	}

	tx, e := h.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}
	return tx, nil
}

func (h *Handler) Commit(tx *sql.Tx, txItf driver.Tx) error {
	if txItf == nil {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}