-- +goose Up
CREATE TABLE IF NOT EXISTS Roles (
  id                    INT AUTO_INCREMENT NOT NULL,
  roleName              VARCHAR(255) NOT NULL,
  displayName           VARCHAR(255) NOT NULL,
  colour                INT NOT NULL,
  info                  VARCHAR(510),
  perms                 VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- ADD ALL THE INITIAL ROLES TO THE DB
INSERT INTO Roles
  (roleName, displayName, colour, info, perms)
VALUES 
  ('arhitectii', 'Arhitectii', 0, 'Fondatorii acestei comunitatii.', ''),
  ('consul', 'Consul', 16776960, 'Autoritatea suprema asupra echipei staff si asupra administrarii si moderarii comunitatii de discord.', ''),
  ('close_friends', 'Close Friends', 16777215, 'Pentru prietenii apropiati ai comunitatii.', ''),
  ('senior_administrator', 'Senior Administrator', 16751052, 'Administratorii seniori ai comunitatii.', ''),
  ('administrator', 'Administrator', 255, 'Administratorii comunitatii.', ''),
  ('moderator', 'Moderator', 32768, 'Moderatorii comunitatii.', ''),
  ('top_contribuitori', 'Top Contribuitori', 4251856, 'Cei care au contribuit cel mai mult la dezvoltarea server-ului prin sugestiile si ideile pe care le-au exprimat si avut.', ''),
  ('server_booster', 'Server Booster', 16711935, 'Cei care au boostat server-ul. (Va multumim frumos)', ''),
  ('content_creator', 'Content Creator', 8388736, 'Diversi youtuberi/influenceri ai comunitatii.', ''),
  ('aztec', 'Aztec', 16711680, 'Membru oficial al acestei comunitatii.', ''),
  ('zelator', 'Zelator', 16777215, '', ''),
  ('theoricus', 'Theoricus', 16777215, '', ''),
  ('practicus', 'Practicus', 16777215, '', ''),
  ('philosophus', 'Philosophus', 16777215, '', ''),
  ('adeptus_minor', 'Adeptus Minor', 7895160, '', ''),
  ('adeptus_major', 'Adeptus Major', 0, '', ''),
  ('adeptus_exemptus', 'Adeptus Exemptus', 12632256, '', ''),
  ('magister_templi', 'Magister Templi', 16766720, '', ''),
  ('magus', 'Magus', 15329769, '', ''),
  ('ipsissimus', 'Ipsissimus', 8388608, '', '');

-- +goose Down
DROP TABLE Roles;