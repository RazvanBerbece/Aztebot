USE aztecaDiscordDb;

CREATE TABLE IF NOT EXISTS User (
  id                   INT AUTO_INCREMENT NOT NULL,
  discordTag           VARCHAR(255) NOT NULL,
  userId               VARCHAR(255) NOT NULL,
  currentRole          VARCHAR(255) NOT NULL,
  currentLevel         VARCHAR(255) NOT NULL,
  currentExperience    VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO User
  (discordTag, userId, currentRole, currentLevel, currentExperience)
VALUES
  ('antoniozrd', '573659533361020941', 'none', 'nil', 'nil'),