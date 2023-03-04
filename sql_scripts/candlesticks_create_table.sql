CREATE TABLE `candlesticks` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `instrument` VARCHAR(20) NOT NULL,
  `time` DATETIME NOT NULL,
  `open` DECIMAL(18, 8) NOT NULL,
  `high` DECIMAL(18, 8) NOT NULL,
  `low` DECIMAL(18, 8) NOT NULL,
  `close` DECIMAL(18, 8) NOT NULL,
  `volume` INT NOT NULL,
  `complete` TINYINT(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_candlesticks` (`instrument`, `time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

