-- +goose Up
-- +goose StatementBegin
ALTER TYPE user_roles_t ADD VALUE 'guardian';

CREATE TYPE student_guardians_relation_t AS ENUM ('mother', 'father', 'guardian', 'relative');

CREATE TABLE student_guardians
(
    id         UUID PRIMARY KEY                       NOT NULL,
    user_id    UUID                                   NOT NULL,
    student_id UUID                                   NOT NULL,
    relation   student_guardians_relation_t           NOT NULL,
    school_id  UUID                                   NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,

    CONSTRAINT student_guardians_user_id_fkey
        FOREIGN KEY (user_id)       REFERENCES users (id),
    CONSTRAINT student_guardians_school_id_fkey
        FOREIGN KEY (school_id)     REFERENCES schools (id),
    CONSTRAINT student_guardians_student_id_fkey
        FOREIGN KEY (student_id)    REFERENCES students (id),

    CONSTRAINT student_guardians_key
        UNIQUE (student_id, relation)
);

COMMENT ON COLUMN student_guardians.id          IS 'Student guardian identifier';
COMMENT ON COLUMN student_guardians.user_id     IS 'User identifier';
COMMENT ON COLUMN student_guardians.student_id  IS 'Student identifier';
COMMENT ON COLUMN student_guardians.relation    IS 'Student guardian relation';
COMMENT ON COLUMN student_guardians.school_id   IS 'School identifier';

COMMENT ON COLUMN student_guardians.created_at  IS 'Date and time the student guardian was created';
COMMENT ON COLUMN student_guardians.updated_at  IS 'Date and time the student guardian was updated';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE student_guardians;

DROP TYPE student_guardians_relation_t;
-- +goose StatementEnd