-- +goose Up
-- +goose StatementBegin
CREATE TABLE teachers
(
    id         UUID PRIMARY KEY                       NOT NULL,
    role_id    UUID                                   NOT NULL,
    user_id    UUID                                   NOT NULL,
    school_id  UUID                                   NOT NULL,
    phone      VARCHAR(30),
    email      VARCHAR(100),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT teachers_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT teachers_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),

    CONSTRAINT teachers_role_id_fkey
        FOREIGN KEY (role_id) REFERENCES user_roles (id),

    CONSTRAINT teachers_role_id_key UNIQUE (role_id)
);

ALTER TYPE user_roles_t ADD VALUE 'teacher';

COMMENT ON COLUMN teachers.id               IS 'Teacher Identifier';
COMMENT ON COLUMN teachers.role_id          IS 'Role id identifier';
COMMENT ON COLUMN teachers.user_id          IS 'User Identifier';
COMMENT ON COLUMN teachers.school_id        IS 'School Identifier';
COMMENT ON COLUMN teachers.phone            IS 'Phone Number';
COMMENT ON COLUMN teachers.email            IS 'Email Address';
COMMENT ON COLUMN teachers.created_at       IS 'Date and time the teacher was created';
COMMENT ON COLUMN teachers.updated_at       IS 'Date and time the teacher was last updated';
COMMENT ON COLUMN teachers.deleted_at       IS 'Date and time the teacher was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE teachers;
-- +goose StatementEnd
