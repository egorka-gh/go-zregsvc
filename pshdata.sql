--
-- Скрипт сгенерирован Devart dbForge Studio for MySQL, Версия 7.3.131.0
-- Домашняя страница продукта: http://www.devart.com/ru/dbforge/mysql/studio
-- Дата скрипта: 13.07.2018 11:13:39
-- Версия сервера: 5.1.73-community
-- Версия клиента: 4.1
--


-- 
-- Отключение внешних ключей
-- 
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;

-- 
-- Установить режим SQL (SQL mode)
-- 
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 
-- Установка кодировки, с использованием которой клиент будет посылать запросы на сервер
--
SET NAMES 'utf8';

--
-- Создать база данных "pshdata"
--
CREATE DATABASE pshdata
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Установка базы данных по умолчанию
--
USE pshdata;

--
-- Удалить таблицу "programs"
--
DROP TABLE IF EXISTS programs;

--
-- Удалить таблицу "program_cards"
--
DROP TABLE IF EXISTS program_cards;

--
-- Удалить таблицу "log_action"
--
DROP TABLE IF EXISTS log_action;

--
-- Удалить таблицу "gender"
--
DROP TABLE IF EXISTS gender;

--
-- Удалить таблицу "clients"
--
DROP TABLE IF EXISTS clients;

--
-- Удалить таблицу "client_state_msg"
--
DROP TABLE IF EXISTS client_state_msg;

--
-- Удалить таблицу "client_state_log"
--
DROP TABLE IF EXISTS client_state_log;

--
-- Удалить таблицу "client_state"
--
DROP TABLE IF EXISTS client_state;

--
-- Установка базы данных по умолчанию
--
USE pshdata;

--
-- Создать таблицу "client_state"
--
CREATE TABLE client_state (
  id int(5) NOT NULL,
  name varchar(50) DEFAULT NULL,
  PRIMARY KEY (id)
)
ENGINE = INNODB
AVG_ROW_LENGTH = 1638
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "client_state_log"
--
CREATE TABLE client_state_log (
  id int(11) NOT NULL AUTO_INCREMENT,
  program int(5) NOT NULL,
  card varchar(50) NOT NULL,
  state int(5) NOT NULL,
  state_date datetime NOT NULL,
  comment varchar(250) DEFAULT NULL,
  action tinyint(4) DEFAULT 0,
  PRIMARY KEY (id),
  INDEX IDX_clients_state_log_card (program, card),
  INDEX IDX_clients_state_log_date (state_date),
  INDEX IDX_clients_state_log_state (state)
)
ENGINE = INNODB
AUTO_INCREMENT = 41
AVG_ROW_LENGTH = 409
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "client_state_msg"
--
CREATE TABLE client_state_msg (
  id int(5) NOT NULL,
  web_comment varchar(500) DEFAULT NULL,
  PRIMARY KEY (id),
  CONSTRAINT FK_client_state_msg_id FOREIGN KEY (id)
  REFERENCES client_state (id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE = INNODB
AVG_ROW_LENGTH = 2340
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "clients"
--
CREATE TABLE clients (
  program int(5) NOT NULL,
  card varchar(50) NOT NULL,
  state int(5) NOT NULL DEFAULT 1,
  state_date datetime DEFAULT NULL,
  surname varchar(100) DEFAULT NULL,
  name varchar(100) DEFAULT NULL,
  patronymic varchar(100) DEFAULT NULL,
  phone_code varchar(10) DEFAULT NULL,
  phone varchar(14) DEFAULT NULL,
  email varchar(50) DEFAULT NULL,
  gender int(5) DEFAULT 0,
  birthday datetime DEFAULT NULL,
  pet varchar(100) DEFAULT NULL,
  send_promo tinyint(1) DEFAULT 1,
  sync tinyint(1) DEFAULT 0,
  sync2 tinyint(1) DEFAULT 0,
  PRIMARY KEY (program, card),
  INDEX IDX_clients (state),
  INDEX IDX_clients_state_date (state, state_date),
  INDEX IDX_clients_sync (sync),
  INDEX IDX_clients_sync2 (state, sync2),
  CONSTRAINT FK_clients_program FOREIGN KEY (program)
  REFERENCES programs (id) ON DELETE NO ACTION ON UPDATE RESTRICT,
  CONSTRAINT FK_clients_state FOREIGN KEY (state)
  REFERENCES client_state (id) ON DELETE NO ACTION ON UPDATE RESTRICT
)
ENGINE = INNODB
AVG_ROW_LENGTH = 4096
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "gender"
--
CREATE TABLE gender (
  id int(5) NOT NULL,
  name varchar(50) DEFAULT NULL,
  PRIMARY KEY (id)
)
ENGINE = INNODB
AVG_ROW_LENGTH = 2048
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "log_action"
--
CREATE TABLE log_action (
  id int(5) NOT NULL,
  name varchar(50) DEFAULT NULL,
  PRIMARY KEY (id)
)
ENGINE = INNODB
AVG_ROW_LENGTH = 2048
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "program_cards"
--
CREATE TABLE program_cards (
  id int(11) NOT NULL AUTO_INCREMENT,
  program int(5) NOT NULL,
  card_start varchar(50) NOT NULL,
  card_end varchar(50) NOT NULL,
  active tinyint(1) DEFAULT 1,
  card_len int(11) NOT NULL,
  check_issued tinyint(1) DEFAULT 0,
  PRIMARY KEY (id),
  INDEX IDX_program_cards (card_len, card_start, card_end),
  CONSTRAINT FK_program_cards_program FOREIGN KEY (program)
  REFERENCES programs (id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE = INNODB
AUTO_INCREMENT = 3
AVG_ROW_LENGTH = 8192
CHARACTER SET utf8
COLLATE utf8_general_ci;

--
-- Создать таблицу "programs"
--
CREATE TABLE programs (
  id int(5) NOT NULL AUTO_INCREMENT,
  name varchar(50) DEFAULT NULL,
  alias varchar(50) DEFAULT NULL,
  external tinyint(1) DEFAULT 0,
  active tinyint(1) DEFAULT 1,
  PRIMARY KEY (id)
)
ENGINE = INNODB
AUTO_INCREMENT = 2
AVG_ROW_LENGTH = 16384
CHARACTER SET utf8
COLLATE utf8_general_ci;

-- 
-- Вывод данных для таблицы client_state
--
INSERT INTO client_state VALUES
(-1001, 'Соглашение'),
(-1000, 'Инициализация'),
(-12, 'Карта не выдана'),
(-11, 'Не верный статус'),
(-10, 'Указана не верная карта'),
(-1, 'Ошибка базы данных'),
(1, 'Выдана'),
(5, 'Регистрация'),
(10, 'Уточнение анкеты'),
(100, 'Зарегистрирован');

-- 
-- Вывод данных для таблицы client_state_log
--
INSERT INTO client_state_log VALUES
(1, 1, '44551', 1, '2017-09-18 14:29:52', NULL, 0),
(2, 1, '44551', 10, '2017-09-18 14:37:13', NULL, 0),
(3, 1, '44551', 100, '2017-09-18 14:37:32', NULL, 0),
(4, 1, '020001', 1, '2017-09-18 14:38:51', NULL, 0),
(5, 1, '44552', 1, '2017-09-18 14:46:56', NULL, 0),
(6, 1, '44553', 1, '2017-10-06 10:18:46', NULL, 0),
(7, 1, '44553', 10, '2017-10-06 10:19:15', NULL, 0),
(8, 1, '44553', 1, '2017-10-06 10:19:27', NULL, 0),
(9, 1, '44552', 10, '2017-10-06 16:08:13', NULL, 0),
(10, 1, '020001', 100, '2017-10-06 17:01:24', NULL, 0),
(11, 1, '44553', 100, '2017-10-06 17:01:50', NULL, 0),
(12, 1, '020001', 1, '2017-10-06 17:02:19', NULL, 0),
(13, 1, '44551', 1, '2017-10-06 17:02:20', NULL, 0),
(14, 1, '44552', 1, '2017-10-06 17:02:21', NULL, 0),
(15, 1, '44553', 1, '2017-10-06 17:02:22', NULL, 0),
(16, 1, '020001', 100, '2017-10-06 17:02:46', NULL, 0),
(17, 1, '44551', 100, '2017-10-06 17:05:53', NULL, 0),
(18, 1, '44552', 10, '2017-10-06 17:05:55', NULL, 0),
(19, 1, '44553', 100, '2017-10-06 17:05:59', NULL, 0),
(20, 1, '44552', 100, '2017-10-06 17:06:11', NULL, 0),
(21, 1, '020001', 1, '2017-10-06 17:10:21', NULL, 0),
(22, 1, '44551', 1, '2017-10-06 17:10:21', NULL, 0),
(23, 1, '44552', 1, '2017-10-06 17:10:22', NULL, 0),
(24, 1, '44553', 1, '2017-10-06 17:10:24', NULL, 0),
(25, 1, '020001', 100, '2017-10-06 17:10:39', NULL, 0),
(26, 1, '44551', 10, '2017-10-06 17:10:41', NULL, 0),
(27, 1, '44552', 100, '2017-10-06 17:10:43', NULL, 0),
(28, 1, '44553', 10, '2017-10-06 17:10:45', NULL, 0),
(29, 1, '44551', 100, '2017-10-06 17:11:03', NULL, 0),
(30, 1, '44553', 100, '2017-10-06 17:11:05', NULL, 0),
(31, 1, '020001', 1, '2017-10-06 17:11:36', NULL, 0),
(32, 1, '44551', 1, '2017-10-06 17:11:41', NULL, 0),
(33, 1, '44552', 1, '2017-10-06 17:11:42', NULL, 0),
(34, 1, '44553', 1, '2017-10-06 17:11:42', NULL, 0),
(35, 1, '020001', 100, '2017-10-06 17:13:02', NULL, 0),
(36, 1, '44552', 1, '2017-10-13 11:38:37', '192.168.30.196:igorz', 2),
(37, 1, '44552', 100, '2017-10-13 11:39:10', '192.168.30.196:igorz', 1),
(38, 1, '44551', 10, '2017-10-13 11:39:17', '192.168.30.196:igorz', 1),
(39, 1, '44553', 100, '2017-10-13 11:39:26', '192.168.30.196:igorz', 1),
(40, 1, '020001', 100, '2017-10-13 11:41:21', '192.168.30.196:igorz', 2);

-- 
-- Вывод данных для таблицы client_state_msg
--
INSERT INTO client_state_msg VALUES
(-1001, 'С <a href="http://www.tut.by" target="_blank" >правилами и условиями</a>  пользования дисконтной картой ознакомлен(а) и согласен(а). Приложение №1.'),
(-1000, 'Укажите код указанный на карте'),
(-12, 'Нет данных о выдаче карты пользователю. Повторите попытку позже. '),
(-10, ' <span color="blue"> Проверьте код карты</span>'),
(-1, 'Сервис не доступен. Попробуйте повторить попытку позже.'),
(5, 'Ваша анкета ожидает поверки на корректность заполнения.'),
(10, 'Ваша анкета не корректно заполнена. Для уточнения анкетных данных обратитесь в место получения карты.'),
(100, 'Ваша анкета зарегистрирована.');

-- 
-- Вывод данных для таблицы clients
--
INSERT INTO clients VALUES
(1, '020001', 100, '2017-10-06 17:13:02', 'gsd', 'ggggg', NULL, '44', '1234567', '', 2, '2002-12-31 00:00:00', '4646dfgdf', 0, 0, 0),
(1, '44551', 5, '2017-11-17 15:22:31', '44551', '44551', '44551', '28', '7777777', '', 1, NULL, NULL, 1, 1, 0),
(1, '44552', 100, '2017-10-13 11:39:10', '44552', '44552', NULL, '17', '9874561', NULL, 1, NULL, '44552', 1, 1, 0),
(1, '44553', 100, '2017-10-13 11:39:26', 'Test', 'Тест', 'Тестович', '33', '4455667', NULL, 0, NULL, NULL, 1, 1, 0);

-- 
-- Вывод данных для таблицы gender
--
INSERT INTO gender VALUES
(0, '-'),
(1, 'М'),
(2, 'Ж');

-- 
-- Вывод данных для таблицы log_action
--
INSERT INTO log_action VALUES
(-1, 'Удаление'),
(0, 'Добавление'),
(1, 'Смена статуса'),
(2, 'Сохранение');

-- 
-- Вывод данных для таблицы program_cards
--
INSERT INTO program_cards VALUES
(1, 1, '25000', '80000', 1, 5, 0),
(2, 1, '020001', '025000', 1, 6, 1);

-- 
-- Вывод данных для таблицы programs
--
INSERT INTO programs VALUES
(1, 'Собственная ПЛ', NULL, 0, 1);

--
-- Установка базы данных по умолчанию
--
USE pshdata;

DELIMITER $$

--
-- Удалить триггер "tg_clients_bu"
--
DROP TRIGGER IF EXISTS tg_clients_bu$$

--
-- Удалить триггер "tg_clients_bi"
--
DROP TRIGGER IF EXISTS tg_clients_bi$$

--
-- Удалить триггер "tg_clients_ai"
--
DROP TRIGGER IF EXISTS tg_clients_ai$$

DELIMITER ;

--
-- Установка базы данных по умолчанию
--
USE pshdata;

DELIMITER $$

--
-- Создать триггер "tg_clients_ai"
--
CREATE
DEFINER = 'root'@'localhost'
TRIGGER tg_clients_ai
AFTER INSERT
ON clients
FOR EACH ROW
BEGIN
  INSERT INTO client_state_log (program
  , card
  , state
  , state_date)
    VALUES (NEW.program
    , NEW.card
    , NEW.state
    , NOW());
END
$$

--
-- Создать триггер "tg_clients_bi"
--
CREATE
DEFINER = 'root'@'localhost'
TRIGGER tg_clients_bi
BEFORE INSERT
ON clients
FOR EACH ROW
BEGIN
  SET NEW.state_date = NOW();
END
$$

--
-- Создать триггер "tg_clients_bu"
--
CREATE
DEFINER = 'root'@'localhost'
TRIGGER tg_clients_bu
BEFORE UPDATE
ON clients
FOR EACH ROW
BEGIN
  IF NOT (OLD.state <=> NEW.state)
  THEN
    SET NEW.state_date = NOW();
  END IF;
END
$$

DELIMITER ;
-- 
-- Восстановить предыдущий режим SQL (SQL mode)
-- 
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;

-- 
-- Включение внешних ключей
-- 
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
