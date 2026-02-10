-- +goose Up
-- +goose StatementBegin
CREATE TABLE group_subjects
(
    id                uuid PRIMARY KEY                       NOT NULL,
    school_subject_id uuid                                   NOT NULL,
    group_id          uuid                                   NOT NULL,
    teacher_id        uuid,
    count             smallint,

    created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at        TIMESTAMP WITH TIME ZONE,
-- TODO check unique when we set deleted_at is case we want to soft delete( maybe need to add condition for uniq)
    CONSTRAINT group_subjects_key
        UNIQUE (school_subject_id, group_id),
    CONSTRAINT group_subjects_school_subject_id_fkey
        FOREIGN KEY (school_subject_id) REFERENCES school_subjects (id),
    CONSTRAINT group_subjects_group_id_fkey
        FOREIGN KEY (group_id) REFERENCES groups (id),
    CONSTRAINT group_subjects_teacher_id_fkey
        FOREIGN KEY (teacher_id) REFERENCES teachers (id)
);

COMMENT ON COLUMN group_subjects.id                 IS 'Group subject identifier';
COMMENT ON COLUMN group_subjects.school_subject_id  IS 'School subject identifier';
COMMENT ON COLUMN group_subjects.group_id           IS 'Group identifier';
COMMENT ON COLUMN group_subjects.teacher_id         IS 'Teacher identifier';
COMMENT ON COLUMN group_subjects.count              IS 'Group subject lesson count';

COMMENT ON COLUMN group_subjects.created_at         IS 'Date and time the group subject was created';
COMMENT ON COLUMN group_subjects.updated_at         IS 'Date and time the group subject was updated';
COMMENT ON COLUMN group_subjects.deleted_at         IS 'Date and time the group subject was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE group_subjects;
-- +goose StatementEnd
