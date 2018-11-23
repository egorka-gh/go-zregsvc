SET NAMES 'utf8';
SET SESSION sql_mode='NO_AUTO_VALUE_ON_ZERO';

--
-- Изменить таблицу "clients"
--
ALTER TABLE clients
  ADD COLUMN version INT(11) NOT NULL DEFAULT 0 AFTER sync2,
  ADD COLUMN deleted TINYINT(1) NOT NULL DEFAULT 0 AFTER version;
ALTER TABLE clients
  ADD INDEX IDX_clients_version (version);

--
-- Изменить таблицу "programs"
--
ALTER TABLE programs
  ADD COLUMN version INT(11) NOT NULL DEFAULT 0 AFTER active,
  ADD COLUMN deleted TINYINT(1) NOT NULL DEFAULT 0 AFTER version;

ALTER TABLE programs
  ADD INDEX IDX_programs_version (version); 

--
-- Изменить таблицу "program_cards"
--
ALTER TABLE program_cards
  ADD COLUMN version INT(11) NOT NULL DEFAULT 0 AFTER check_issued,
  ADD COLUMN deleted TINYINT(1) NOT NULL DEFAULT 0 AFTER version;

ALTER TABLE program_cards
  ADD INDEX IDX_program_cards_version (version);

--
-- Создать таблицу "cnv_source"
--
CREATE TABLE cnv_source (
  id CHAR(2) DEFAULT NULL,
  name VARCHAR(50) DEFAULT NULL
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

INSERT INTO cnv_source(id, name) VALUES('00', 'Центральная база');

--
-- Создать таблицу "cnv_version"
--
CREATE TABLE cnv_version (
  source char(2) NOT NULL DEFAULT '',
  table_name varchar(40) NOT NULL,
  latest_version int(11) NOT NULL DEFAULT 0,
  syncorder int(5) NOT NULL DEFAULT 0,
  PRIMARY KEY (source, table_name)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

CREATE TABLE client_activity (
  source char(2) NOT NULL,
  doc_id varchar(50) NOT NULL,
  card varchar(50) NOT NULL,
  doc_date datetime NOT NULL,
  doc_sum decimal(16, 2) NOT NULL DEFAULT 0.00,
  bonuce_sum decimal(16, 2) NOT NULL DEFAULT 0.00,
  version int(11) NOT NULL DEFAULT 0,
  deleted tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (source, doc_id),
  INDEX IDX_client_activity (source, version),
  INDEX IDX_client_activity_card (card)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

CREATE TABLE pshdata.client_balance (
  card varchar(50) NOT NULL,
  balance_date date NOT NULL COMMENT 'the last day of the corresponding month',
  doc_sum decimal(16, 2) DEFAULT 0.00,
  bonuce_sum decimal(16, 2) DEFAULT 0.00,
  level int(5) NOT NULL DEFAULT 0,
  version int(11) NOT NULL DEFAULT 0,
  deleted tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (card, balance_date),
  INDEX IDX_client_balance_version (version)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

CREATE TABLE balance_levels (
  level int(5) NOT NULL,
  bonuce_sum decimal(16, 2) NOT NULL DEFAULT 0.00,
  version int(11) NOT NULL DEFAULT 0,
  deleted tinyint(1) NOT NULL DEFAULT 0,
  comment varchar(50) DEFAULT NULL,
  PRIMARY KEY (level)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

INSERT INTO balance_levels(level, bonuce_sum, version, deleted, comment) VALUES(0, 0.00, 0, 0, 'Незарегистрированная карта');

CREATE TABLE pshdata.balance_levels_Log (
  log_date date NOT NULL,
  level int(5) NOT NULL,
  bonuce_sum decimal(16, 2) NOT NULL DEFAULT 0.00,
  user varchar(50) DEFAULT NULL,
  comment varchar(50) DEFAULT NULL,
  PRIMARY KEY (log_date, level)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

INSERT INTO cnv_version(source, table_name, latest_version, syncorder) VALUES('', 'client_activity', 0, 0);
INSERT INTO cnv_version(source, table_name, latest_version, syncorder) VALUES('00', 'programs', 0, 0);
INSERT INTO cnv_version(source, table_name, latest_version, syncorder) VALUES('00', 'program_cards', 1, 0);
INSERT INTO cnv_version(source, table_name, latest_version, syncorder) VALUES('00', 'clients', 2, 0);
INSERT INTO cnv_version(source, table_name, latest_version, syncorder) VALUES('00', 'client_balance', 3, 0);
