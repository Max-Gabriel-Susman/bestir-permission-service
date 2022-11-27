-- +goose Up
CREATE TABLE IF NOT EXISTS permission (
    id CHAR(36) NOT NULL, 
    [name] VARCHAR(255) NOT NULL, 
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- +goose Down 
DROP TABLE IF EXISTS permission;