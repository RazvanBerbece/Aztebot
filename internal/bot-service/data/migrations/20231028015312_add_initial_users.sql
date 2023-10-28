-- +goose Up
CREATE TABLE IF NOT EXISTS Users (
  id                   INT AUTO_INCREMENT NOT NULL,
  discordTag           VARCHAR(255) NOT NULL,
  userId               VARCHAR(255) NOT NULL,
  currentRoleIds       VARCHAR(255) NOT NULL,
  currentCircle        VARCHAR(255) NOT NULL,   -- INNER/OUTER
  currentInnerOrder    INT,                     -- CAN BE NULL IF USER NOT IN THE INNER CIRCLE, NUMERAL OTHERWISE (1-3)
  currentLevel         int NOT NULL,
  currentExperience    int NOT NULL,            -- CUMULATION OF POINTS DRIVEN BY CONTRIBUTIONS, HOURS SPENT, ETC.
  PRIMARY KEY (`id`)
);
INSERT INTO Users
  (discordTag, userId, currentRoleIds, currentCircle, currentInnerOrder, currentLevel, currentExperience)
VALUES
  ('antoniozrd', '573659533361020941', '4, 10', 'OUTER', NULL, 1, 0),
  ('lordvixxen1337', '1077147870655950908', '1, 10, 20', 'INNER', 3, 100, 999),
  ('aztegramul', '526512064794066945', '1, 10, 20', 'INNER', 3, 100, 999);

-- +goose Down
DROP TABLE Users;