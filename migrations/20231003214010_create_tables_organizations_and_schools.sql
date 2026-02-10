-- +goose Up
-- +goose StatementBegin
CREATE TABLE educational_organizations
(
    id              UUID    PRIMARY KEY      NOT NULL,
    name            VARCHAR(150)             NOT NULL,
    logo            VARCHAR(255),
    description     TEXT,

    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT educational_organizations_name_key UNIQUE (name)
);

COMMENT ON COLUMN educational_organizations.id          IS 'Educational Organizations Identifier';
COMMENT ON COLUMN educational_organizations.name        IS 'Name';
COMMENT ON COLUMN educational_organizations.logo        IS 'Link to the organizations logo';
COMMENT ON COLUMN educational_organizations.created_at  IS 'Date and time the entry was created';
COMMENT ON COLUMN educational_organizations.updated_at  IS 'Date and time the entry was updated';
COMMENT ON COLUMN educational_organizations.deleted_at  IS 'Date and time the entry was deleted';

CREATE TABLE schools
(
    id              UUID  PRIMARY KEY           NOT NULL,
    organization_id UUID                        NOT NULL,
    name            VARCHAR(150)                NOT NULL,
    location        VARCHAR(255)                NOT NULL,
    phone           VARCHAR(20),
    email           VARCHAR(100),
    
    created_at      TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT schools_name_key UNIQUE (name),
    CONSTRAINT schools_organization_id_fkey
        FOREIGN KEY (organization_id) REFERENCES educational_organizations(id)
);

COMMENT ON COLUMN schools.id                IS 'Identification';
COMMENT ON COLUMN schools.organization_id   IS 'Organizations identification';
COMMENT ON COLUMN schools.name              IS 'Name';
COMMENT ON COLUMN schools.location          IS 'Address';
COMMENT ON COLUMN schools.phone             IS 'Phone';
COMMENT ON COLUMN schools.email             IS 'Email';
COMMENT ON COLUMN schools.created_at        IS 'Date and time the entry was created';
COMMENT ON COLUMN schools.updated_at        IS 'Date and time the entry was updated';
COMMENT ON COLUMN schools.deleted_at        IS 'Date and time the entry was deleted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE schools;

DROP TABLE educational_organizations;

-- +goose StatementEnd
