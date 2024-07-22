-- +goose Up

CREATE TABLE operator_states (
  pipeline_id varchar(255) NOT NULL,
  operator_id  varchar(255) NOT NULL,
  state varchar(255) NOT NULL,
  container_id varchar(255),
  error varchar(255),

  PRIMARY KEY(pipeline_id, operator_id)
);

-- +goose Down
DROP TABLE operator_states;