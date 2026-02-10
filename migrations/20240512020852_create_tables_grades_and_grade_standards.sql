-- +goose Up
-- +goose StatementBegin
CREATE TABLE grade_standards
(
    id              UUID PRIMARY KEY         NOT NULL,
    organization_id UUID,
    name            VARCHAR(50)              NOT NULL,
    education_years SMALLINT                 NOT NULL,
    description     TEXT                     NOT NULL,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT grade_standards_name_key UNIQUE (name),
    CONSTRAINT grade_standards_organization_id_fkey
        FOREIGN KEY (organization_id) REFERENCES educational_organizations (id)
);

COMMENT ON COLUMN grade_standards.id                IS 'Grade standard identifier';
COMMENT ON COLUMN grade_standards.organization_id   IS 'Organization identifier';
COMMENT ON COLUMN grade_standards.name              IS 'Template name';
COMMENT ON COLUMN grade_standards.education_years   IS 'Number of years of study';
COMMENT ON COLUMN grade_standards.description       IS 'Template description';
COMMENT ON COLUMN grade_standards.created_at        IS 'Date and time the grade standard was created';
COMMENT ON COLUMN grade_standards.updated_at        IS 'Date and time the grade standard was updated';
COMMENT ON COLUMN grade_standards.deleted_at        IS 'Date and time the grade standard was deleted';

CREATE TABLE grades
(
    id                UUID PRIMARY KEY         NOT NULL,
    grade_standard_id UUID                     NOT NULL,
    name              VARCHAR(40)              NOT NULL,
    education_year    SMALLINT,

    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT grades_name_key UNIQUE (name, grade_standard_id),
    CONSTRAINT grades_grade_standard_id_fkey
        FOREIGN KEY (grade_standard_id) REFERENCES grade_standards (id)
);

COMMENT ON COLUMN grades.id                 IS 'Grade identifier';
COMMENT ON COLUMN grades.grade_standard_id  IS 'Grade standard identifier';
COMMENT ON COLUMN grades.name               IS 'Grade name';
COMMENT ON COLUMN grades.education_year     IS 'Ordinal year of study';
COMMENT ON COLUMN grades.created_at         IS 'Date and time the grade was created';
COMMENT ON COLUMN grades.updated_at         IS 'Date and time the grade was updated';
COMMENT ON COLUMN grades.deleted_at         IS 'Date and time the grade was deleted';

ALTER TABLE schools
    ADD COLUMN grade_standard_id UUID;

ALTER TABLE schools
    ADD CONSTRAINT schools_grade_standard_id_fkey
        FOREIGN KEY (grade_standard_id) REFERENCES grade_standards(id);

COMMENT ON COLUMN schools.grade_standard_id  IS 'grade standard identification';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE schools
    DROP COLUMN grade_standard_id;

DROP TABLE grades;

DROP TABLE grade_standards;

-- +goose StatementEnd
