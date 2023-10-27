-- +goose Up
CREATE TABLE IF NOT EXISTS Roles (
  id                    INT AUTO_INCREMENT NOT NULL,
  fullName              VARCHAR(255) NOT NULL,
  displayName           VARCHAR(255) NOT NULL,
  colour                INT NOT NULL,
  minimumXp             INT NOT NULL,
  maximumXp             INT NOT NULL,
  perms                 VARCHAR(255),
  PRIMARY KEY (`id`)
);
INSERT INTO Roles
  (fullName, displayName, colour, minimumXp, maximumXp, perms)
VALUES
  ('initiate', 'Initiate', 16711680, 0, 1000, "");

-- +goose Down
DROP TABLE Roles;