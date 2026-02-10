-- +goose Up
-- +goose StatementBegin
CREATE TABLE marks
(
    id                  uuid PRIMARY KEY                        NOT NULL,
    lesson_id           uuid                                    NOT NULL,
    student_id          uuid                                    NOT NULL,
    mark                uuid                                    NOT NULL,
    description         text,

    created_at          TIMESTAMP WITH TIME ZONE DEFAULT now()  NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT now()  NOT NULL,
    deleted_at          TIMESTAMP WITH TIME ZONE,

    CONSTRAINT marks_lesson_id_fkey
        FOREIGN KEY (lesson_id) REFERENCES lessons (id),
    CONSTRAINT marks_student_id_fkey
        FOREIGN KEY (student_id) REFERENCES students (id),
    CONSTRAINT marks_lesson_key UNIQUE (lesson_id, student_id)
);

COMMENT ON COLUMN marks.id IS 'mark identifier';
COMMENT ON COLUMN marks.lesson_id IS 'lessons identifier';
COMMENT ON COLUMN marks.student_id IS 'student identifier';
COMMENT ON COLUMN marks.mark IS 'mark';
COMMENT ON COLUMN marks.description IS 'mark description';

COMMENT ON COLUMN marks.created_at IS 'Date and time the lesson mark was created';
COMMENT ON COLUMN marks.updated_at IS 'Date and time the lesson mark was updated';
COMMENT ON COLUMN marks.deleted_at IS 'Date and time the lesson mark was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE marks;
-- +goose StatementEnd