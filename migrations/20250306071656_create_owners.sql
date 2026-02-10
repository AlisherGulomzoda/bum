-- +goose Up
-- +goose StatementBegin
CREATE TABLE owners
(
    id              UUID PRIMARY KEY                        NOT NULL,
    role_id         UUID                                    NOT NULL,
    user_id         UUID                                    NOT NULL,
    organization_id UUID                                    NOT NULL,
    phone           VARCHAR(30),
    email           VARCHAR(100),

    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()  NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()  NOT NULL,
    deleted_at      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT owners_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT owners_organization_id_fkey
        FOREIGN KEY (organization_id) REFERENCES educational_organizations (id),

    CONSTRAINT owners_role_id_fkey
        FOREIGN KEY (role_id) REFERENCES user_roles (id),

    CONSTRAINT owners_role_id_key UNIQUE (role_id)
);

ALTER TYPE user_roles_t ADD VALUE 'owner';

COMMENT ON COLUMN owners.id               IS 'Owner Identifier';
COMMENT ON COLUMN owners.role_id          IS 'Role id identifier';
COMMENT ON COLUMN owners.user_id          IS 'User Identifier';
COMMENT ON COLUMN owners.organization_id  IS 'Edu organization Identifier';
COMMENT ON COLUMN owners.phone            IS 'Phone Number';
COMMENT ON COLUMN owners.email            IS 'Email Address';

COMMENT ON COLUMN owners.created_at       IS 'Date and time the owner was created';
COMMENT ON COLUMN owners.updated_at       IS 'Date and time the owner was last updated';
COMMENT ON COLUMN owners.deleted_at       IS 'Date and time the owner was deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE owners;
-- +goose StatementEnd
