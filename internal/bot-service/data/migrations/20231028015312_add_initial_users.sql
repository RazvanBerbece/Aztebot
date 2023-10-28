-- +goose Up
CREATE TABLE IF NOT EXISTS Users (
  id                   INT AUTO_INCREMENT NOT NULL,
  discordTag           VARCHAR(255) NOT NULL,
  userId               VARCHAR(255) NOT NULL,
  currentRoleId        int NOT NULL,
  currentCircle        VARCHAR(255) NOT NULL,   -- INNER/OUTER
  currentInnerOrder    VARCHAR(255),            -- CAN BE NULL IF USER NOT IN THE INNER CIRCLE
  currentLevel         int NOT NULL,
  currentExperience    int NOT NULL,            -- CUMULATION OF POINTS DRIVEN BY CONTRIBUTIONS, HOURS SPENT, ETC.
  PRIMARY KEY (`id`),
  FOREIGN KEY (currentRoleId) REFERENCES Roles(id)
);
INSERT INTO Users
  (discordTag, userId, currentRoleId, currentCircle, currentInnerOrder, currentLevel, currentExperience)
VALUES
  ('antoniozrd', '573659533361020941', 4, 'OUTER', '', 1, 0);

-- +goose Down
DROP TABLE Users;