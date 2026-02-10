-- +goose Up
-- +goose StatementBegin
CREATE TABLE auditoriums
(
    id                uuid PRIMARY KEY                       NOT NULL,
    school_id         uuid                                   NOT NULL,
    name              varchar(50)                            NOT NULL,
    school_subject_id uuid,
    description       text,

    created_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at        TIMESTAMP WITH TIME ZONE,

    CONSTRAINT auditoriums_name_key
        UNIQUE (school_id, name),
    CONSTRAINT auditoriums_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),
    CONSTRAINT auditoriums_school_subject_id_fkey
        FOREIGN KEY (school_subject_id) REFERENCES school_subjects (id)
);

COMMENT ON COLUMN auditoriums.id                IS 'Auditorium identifier';
COMMENT ON COLUMN auditoriums.school_id         IS 'School identifier';
COMMENT ON COLUMN auditoriums.name              IS 'Auditorium name';
COMMENT ON COLUMN auditoriums.school_subject_id IS 'Assignment auditorium school subject identifier';
COMMENT ON COLUMN auditoriums.description       IS 'School auditorium description';

COMMENT ON COLUMN auditoriums.created_at IS 'Date and time the auditorium was created';
COMMENT ON COLUMN auditoriums.updated_at IS 'Date and time the auditorium was updated';
COMMENT ON COLUMN auditoriums.deleted_at IS 'Date and time the auditorium was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auditoriums;
-- +goose StatementEnd
