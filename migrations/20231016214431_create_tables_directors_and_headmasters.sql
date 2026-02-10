-- +goose Up
-- +goose StatementBegin
CREATE TABLE directors
(
    id         UUID PRIMARY KEY         NOT NULL,
    role_id    UUID                     NOT NULL,
    user_id    UUID                     NOT NULL,
    school_id  UUID                     NOT NULL,
    phone      VARCHAR(20),
    email      VARCHAR(100),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT directors_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT directors_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),

    CONSTRAINT directors_role_id_fkey
        FOREIGN KEY (role_id) REFERENCES user_roles (id),

    CONSTRAINT directors_role_id_key UNIQUE (role_id)
);

COMMENT ON COLUMN directors.id              IS 'Director identifier';
COMMENT ON COLUMN directors.role_id         IS 'Role id identifier';
COMMENT ON COLUMN directors.user_id         IS 'User identifier';
COMMENT ON COLUMN directors.school_id       IS 'School identifier';
COMMENT ON COLUMN directors.phone           IS 'Phone';
COMMENT ON COLUMN directors.email           Is 'Email';
COMMENT ON COLUMN directors.created_at      IS 'Date and time the director was created';
COMMENT ON COLUMN directors.updated_at      IS 'Date and time the director was updated';
COMMENT ON COLUMN directors.deleted_at      IS 'Date and time the director was deleted';

ALTER TYPE user_roles_t ADD VALUE 'director';

CREATE TABLE headmasters
(
    id         UUID PRIMARY KEY         NOT NULL,
    role_id    UUID                     NOT NULL,
    user_id    UUID                     NOT NULL,
    school_id  UUID                     NOT NULL,
    phone      VARCHAR(20),
    email      VARCHAR(100),
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT headmasters_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT headmasters_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),

    CONSTRAINT headmasters_role_id_fkey
        FOREIGN KEY (role_id) REFERENCES user_roles (id),

    CONSTRAINT headmasters_role_id_key UNIQUE (role_id)
);

COMMENT ON COLUMN headmasters.id            IS 'Headmaster identifier';
COMMENT ON COLUMN headmasters.role_id       IS 'Role id identifier';
COMMENT ON COLUMN headmasters.user_id       IS 'User identifier';
COMMENT ON COLUMN headmasters.school_id     IS 'School identifier';
COMMENT ON COLUMN headmasters.phone         IS 'Phone';
COMMENT ON COLUMN headmasters.email         Is 'Email';
COMMENT ON COLUMN headmasters.created_at    IS 'Date and time the headmaster was created';
COMMENT ON COLUMN headmasters.updated_at    IS 'Date and time the headmaster was updated';
COMMENT ON COLUMN headmasters.deleted_at    IS 'Date and time the headmaster was deleted';

ALTER TYPE user_roles_t ADD VALUE 'headmaster';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE directors;

DROP TABLE headmasters;

-- +goose StatementEnd
