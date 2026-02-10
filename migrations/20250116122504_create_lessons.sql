-- +goose Up
-- +goose StatementBegin
CREATE TABLE lessons
(
    id                  uuid PRIMARY KEY                        NOT NULL,
    school_id           uuid                                    NOT NULL,
    group_subject_id    uuid                                    NOT NULL,
    teacher_id          uuid,
    auditorium_id       uuid                                    NOT NULL,
    start_time          TIMESTAMP WITH TIME ZONE                NOT NULL,
    end_time            TIMESTAMP WITH TIME ZONE                NOT NULL,
    description         text,

    created_at          TIMESTAMP WITH TIME ZONE DEFAULT now()  NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT now()  NOT NULL,
    deleted_at          TIMESTAMP WITH TIME ZONE,

    CONSTRAINT lessons_school_id_fkey
        FOREIGN KEY (school_id) REFERENCES schools (id),
    CONSTRAINT lessons_group_subject_id_fkey
        FOREIGN KEY (group_subject_id) REFERENCES group_subjects (id),
    CONSTRAINT lessons_teacher_id_fkey
        FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    CONSTRAINT lessons_auditorium_id_fkey
        FOREIGN KEY (auditorium_id) REFERENCES auditoriums (id)
);

COMMENT ON COLUMN lessons.id                IS 'lesson identifier';
COMMENT ON COLUMN lessons.school_id         IS 'lessons school id identifier';
COMMENT ON COLUMN lessons.group_subject_id  IS 'group subject identifier';
COMMENT ON COLUMN lessons.teacher_id        IS 'teacher identifier';
COMMENT ON COLUMN lessons.auditorium_id     IS 'auditorium identifier';
COMMENT ON COLUMN lessons.start_time        IS 'lesson start time';
COMMENT ON COLUMN lessons.end_time          IS 'lesson end time';
COMMENT ON COLUMN lessons.description       IS 'lesson description';

COMMENT ON COLUMN lessons.created_at IS 'Date and time the lesson was created';
COMMENT ON COLUMN lessons.updated_at IS 'Date and time the lesson was updated';
COMMENT ON COLUMN lessons.deleted_at IS 'Date and time the lesson was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lessons;
-- +goose StatementEnd
