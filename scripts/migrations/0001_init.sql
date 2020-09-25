-- +goose Up
CREATE TABLE events
(
    id          BIGSERIAL                NOT NULL PRIMARY KEY,
    user_id     BIGINT                   NOT NULL,
    title       CHARACTER VARYING        NOT NULL,
    description text                     NOT NULL,
    start_time  timestamp WITH TIME ZONE NOT NULL,
    duration    BIGINT                   NOT NULL,
    is_deleted  BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE events;