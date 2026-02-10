-- +goose Up
-- +goose StatementBegin
CREATE TYPE users_gender_t AS ENUM ('male', 'female');

CREATE TABLE users
(
    id          UUID PRIMARY KEY         NOT NULL,
    first_name  VARCHAR(25)              NOT NULL,
    middle_name VARCHAR(25),
    last_name   VARCHAR(25)              NOT NULL,
    gender      users_gender_t           NOT NULL,
    phone       VARCHAR(20),
    email       VARCHAR(100)             NOT NULL,
    password    VARCHAR(60)              NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users_phone_key UNIQUE (phone),
    CONSTRAINT users_email_key UNIQUE (email)
);

COMMENT ON COLUMN users.id              IS 'Users identifier';
COMMENT ON COLUMN users.first_name      IS 'Users first name';
COMMENT ON COLUMN users.middle_name     IS 'Users middle name';
COMMENT ON COLUMN users.last_name       IS 'Users last name';
COMMENT ON COLUMN users.gender          IS 'Gender';
COMMENT ON COLUMN users.phone           IS 'Phone';
COMMENT ON COLUMN users.email           IS 'Email';
COMMENT ON COLUMN users.password        IS 'Hashed password';
COMMENT ON COLUMN users.created_at      IS 'Date and time the user was created';
COMMENT ON COLUMN users.updated_at      IS 'Date and time the user was updated';
COMMENT ON COLUMN users.deleted_at      IS 'Date and time the user was deleted';

CREATE TYPE user_roles_t AS ENUM('admin');

CREATE TABLE user_roles
(
    id                  UUID PRIMARY KEY            NOT NULL,
    user_id             UUID                        NOT NULL,
    role                user_roles_t                NOT NULL,
    school_id           UUID,
    organization_id     UUID,

    created_at  TIMESTAMP WITH TIME ZONE            NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE            NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT user_roles_user_id_role_school_id_organization_id_key
        UNIQUE NULLS NOT DISTINCT (user_id, role, school_id, organization_id),
    CONSTRAINT user_roles_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT user_roles_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),
    CONSTRAINT user_roles_organization_id_fkey
        FOREIGN KEY (organization_id) REFERENCES educational_organizations (id)
);

COMMENT ON COLUMN user_roles.id         IS 'Record identifier';
COMMENT ON COLUMN user_roles.user_id    IS 'User identifier';
COMMENT ON COLUMN user_roles.role       IS 'User role';
COMMENT ON COLUMN user_roles.school_id  IS 'School identifier';

COMMENT ON COLUMN user_roles.created_at      IS 'Date and time the user role was created';
COMMENT ON COLUMN user_roles.updated_at      IS 'Date and time the user role was updated';
COMMENT ON COLUMN user_roles.deleted_at      IS 'Date and time the user role was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_roles;

DROP TABLE users;

DROP TYPE users_gender_t;

DROP TYPE user_roles_t;

-- +goose StatementEnd
