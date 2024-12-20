-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS deployments (
    id int not null primary key,
    app_id text not null,
    strategy text not null,
    status text not null default 'pending',
    metadata text not null default '{}',
    state text,
    snapshot text default '{}',
    worker_id text,
    last_heartbeat integer not NULL default 0,
    created_at DATETIME not null DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME not null DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX deployments_app ON deployments (app_id);

CREATE TABLE IF NOT EXISTS events (
    id int not null primary key,
    deployment_id text,
    entity_id text,
    message text,
    action text not null default '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX deployment_events ON events (deployment_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE deployments;
DROP TABLE events;
-- +goose StatementEnd
