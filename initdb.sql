DROP TABLE IF EXISTS access_logs;

CREATE TABLE access_logs (
  id int PRIMARY KEY AUTO_INCREMENT,
  ua text NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP
) COMMENT = 'access log';
