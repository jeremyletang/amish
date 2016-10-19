USE amish;

CREATE TABLE IF NOT EXISTS users
(
  id                  VARCHAR(36)                                                     NOT NULL,
  login               VARCHAR(512)                                                    NOT NULL,
  name                VARCHAR(512)                                                    NOT NULL,
  company             VARCHAR(512)                                                    NOT NULL,
  avatar_url          VARCHAR(512)                                                    NOT NULL,
  location            VARCHAR(512)                                                    NOT NULL,
  blog                VARCHAR(512)                                                    NOT NULL,
  email               VARCHAR(512)                                                    NOT NULL,
  content_updated     INTEGER  DEFAULT 0                                              NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS repositories
(
  id                  VARCHAR(36)                                                     NOT NULL,
  owner               VARCHAR(512)                                                    NOT NULL,
  name                VARCHAR(512)                                                    NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS stars
(
  id                  VARCHAR(36)                                                     NOT NULL,
  repository_id       VARCHAR(36)                                                     NOT NULL,
  user_id             VARCHAR(36)                                                     NOT NULL,
  starred_at          DATETIME                                                        NOT NULL,
  valid               INTEGER                                                         NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE stars
      ADD FOREIGN KEY (repository_id) REFERENCES repositories (id),
      ADD FOREIGN KEY (user_id) REFERENCES users (id);
      
