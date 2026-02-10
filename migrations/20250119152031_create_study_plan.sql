-- +goose Up
-- +goose StatementBegin
CREATE TYPE study_plan_status_t AS ENUM ('planned', 'ongoing', 'completed');

CREATE TABLE study_plans
(
    id               uuid PRIMARY KEY                                                NOT NULL,
    group_subject_id uuid                                                            NOT NULL,
    title            varchar(255)                                                    NOT NULL,
    description      text,
    plan_order       SMALLINT                                                        NOT NULL,
    status           study_plan_status_t      DEFAULT 'planned'::study_plan_status_t NOT NULL,

    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now()                          NOT NULL,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT now()                          NOT NULL,
    deleted_at       TIMESTAMP WITH TIME ZONE,

    CONSTRAINT study_plan_group_subject_id_fkey
        FOREIGN KEY (group_subject_id) REFERENCES group_subjects (id)
);

CREATE UNIQUE INDEX study_plan_group_subject_id_plan_order_key
    ON study_plans (group_subject_id, plan_order)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN study_plans.id                IS 'Study plan identifier';
COMMENT ON COLUMN study_plans.group_subject_id  IS 'School group subject identifier';
COMMENT ON COLUMN study_plans.title             IS 'Study plan title';
COMMENT ON COLUMN study_plans.description       IS 'Study plan description';
COMMENT ON COLUMN study_plans.plan_order        IS 'Study plan order in the group subject';
COMMENT ON COLUMN study_plans.status            IS 'Study plan status';

COMMENT ON COLUMN study_plans.created_at        IS 'Date and time the study plan was created';
COMMENT ON COLUMN study_plans.updated_at        IS 'Date and time the study plan was updated';
COMMENT ON COLUMN study_plans.deleted_at        IS 'Date and time the study plan was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE study_plans;

DROP TYPE study_plan_status_t;
-- +goose StatementEnd
