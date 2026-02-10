-- +goose Up
-- +goose StatementBegin
CREATE TABLE students
(
    id          UUID PRIMARY KEY        NOT NULL,
    role_id     UUID                    NOT NULL,
    group_id    UUID                    NOT NULL,
    user_id     UUID                    NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT  students_group_id_fkey
        FOREIGN KEY (group_id) REFERENCES groups (id),
    CONSTRAINT  students_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),

    CONSTRAINT students_role_id_fkey
            FOREIGN KEY (role_id) REFERENCES user_roles (id),

    CONSTRAINT students_role_id_key UNIQUE (role_id)
);

COMMENT ON COLUMN students.id                   IS 'Group student identifier';
COMMENT ON COLUMN students.role_id              IS 'Role id identifier';
COMMENT ON COLUMN students.group_id             IS 'Student group identifier';
COMMENT ON COLUMN students.user_id              IS 'user identifier';

COMMENT ON COLUMN students.created_at           IS 'Date and time the group student was created';
COMMENT ON COLUMN students.updated_at           IS 'Date and time the group student was updated';
COMMENT ON COLUMN students.deleted_at           IS 'Date and time the group student was deleted';

ALTER TYPE user_roles_t ADD VALUE 'student';

-- adding class president
ALTER TABLE groups
    ADD COLUMN class_president_id UUID;

COMMENT ON COLUMN groups.class_president_id     IS 'Student identifier of class president';

ALTER TABLE groups
    ADD CONSTRAINT groups_class_president_id_fkey
        FOREIGN KEY (class_president_id) REFERENCES students(id);

-- adding deputy for class president
ALTER TABLE groups
    ADD COLUMN deputy_class_president_id UUID;

COMMENT ON COLUMN groups.deputy_class_president_id     IS 'Student identifier of deputy class president';

ALTER TABLE groups
    ADD CONSTRAINT groups_deputy_class_president_id_fkey
        FOREIGN KEY (deputy_class_president_id) REFERENCES students(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE groups
    DROP COLUMN deputy_class_president_id;

ALTER TABLE groups
    DROP COLUMN class_president_id;

DROP TABLE students;
-- +goose StatementEnd