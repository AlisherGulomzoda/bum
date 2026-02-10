-- +goose Up
-- +goose StatementBegin
CREATE TABLE school_subjects
(
    id          uuid PRIMARY KEY        NOT NULL,
    subject_id  uuid                    NOT NULL,
    school_id   uuid                    NOT NULL,
    name        VARCHAR(255)            NOT NULL,
    description TEXT                    NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT  school_subjects_subject_id_fkey
        FOREIGN KEY (subject_id) REFERENCES subjects (id),
    CONSTRAINT  school_subjects_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id)
);

COMMENT ON COLUMN school_subjects.id                 IS 'School subject identifier';
COMMENT ON COLUMN school_subjects.subject_id         IS 'Subject identifier';
COMMENT ON COLUMN school_subjects.school_id          IS 'School identifier';
COMMENT ON COLUMN school_subjects.name               IS 'School subject name';
COMMENT ON COLUMN school_subjects.description        IS 'School subject description';
COMMENT ON COLUMN school_subjects.created_at         IS 'Date and time the grade standard was created';
COMMENT ON COLUMN school_subjects.updated_at         IS 'Date and time the grade standard was updated';
COMMENT ON COLUMN school_subjects.deleted_at         IS 'Date and time the grade standard was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE school_subjects;
-- +goose StatementEnd