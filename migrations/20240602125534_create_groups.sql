-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups
(
    id         UUID PRIMARY KEY        NOT NULL,
    school_id  UUID                    NOT NULL,
    name       VARCHAR(10)             NOT NULL,
    grade_id   UUID                    NOT NULL,
    class_teacher_id UUID,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT groups_name_key UNIQUE (name, grade_id, school_id),
    CONSTRAINT groups_class_teacher_id_fkey
        FOREIGN KEY (class_teacher_id) REFERENCES teachers (id),
    CONSTRAINT groups_grade_id_fkey
        FOREIGN KEY (grade_id) REFERENCES grades (id),
    CONSTRAINT groups_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id)
);

COMMENT ON COLUMN groups.id                 IS 'Grade standard identifier';
COMMENT ON COLUMN groups.school_id          IS 'School identifier';
COMMENT ON COLUMN groups.name               IS 'group name';
COMMENT ON COLUMN groups.grade_id           IS 'Grade identifier';
COMMENT ON COLUMN groups.class_teacher_id   IS 'class teacher identifier';
COMMENT ON COLUMN groups.created_at         IS 'Date and time the grade standard was created';
COMMENT ON COLUMN groups.updated_at         IS 'Date and time the grade standard was updated';
COMMENT ON COLUMN groups.deleted_at         IS 'Date and time the grade standard was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE groups;
-- +goose StatementEnd