-- +goose Up
-- +goose StatementBegin
CREATE TABLE subjects
(
    id          UUID PRIMARY KEY                       NOT NULL,
    name        VARCHAR(255)                           NOT NULL
        CONSTRAINT subjects_name_key UNIQUE,
    description TEXT                                   NOT NULL,
    
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at  TIMESTAMP
);

COMMENT ON COLUMN subjects.id               IS 'subject Identifier';
COMMENT ON COLUMN subjects.name             IS 'subject name';
COMMENT ON COLUMN subjects.description      IS 'subject description';
COMMENT ON COLUMN subjects.created_at       IS 'Date and time the teacher was created';
COMMENT ON COLUMN subjects.updated_at       IS 'Date and time the teacher was last updated';
COMMENT ON COLUMN subjects.deleted_at       IS 'Date and time the teacher was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subjects;
-- +goose StatementEnd