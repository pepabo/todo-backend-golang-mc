SET NAMES utf8mb4 COLLATE utf8mb4_general_ci;

DROP TABLE IF EXISTS todos;

CREATE TABLE todos (
  id int PRIMARY KEY AUTO_INCREMENT,
  title text NOT NULL,
  completed boolean,
  `order` int
) COMMENT = 'TODO';

INSERT INTO todos (title, completed, `order`) VALUES ('マネージドクラウドに登録する', false, 1);
INSERT INTO todos (title, completed, `order`) VALUES ('マネージドクラウドでGoアプリケーションをデプロイする', false, 2);
