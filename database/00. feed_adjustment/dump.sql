-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema livestock
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema livestock
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `livestock` DEFAULT CHARACTER SET latin1 ;
USE `livestock` ;

-- -----------------------------------------------------
-- Table `livestock`.`feed_adjustment`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`feed_adjustment` ;

CREATE TABLE IF NOT EXISTS `livestock`.`feed_adjustment` (
  `id` CHAR(36) NOT NULL,
  `feed_type_id` CHAR(36) NULL,
  `qty` DECIMAL(20,2) NULL,
  `remarks` VARCHAR(255) NULL,
  `created` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_feed_adjustment_feed_type_idx` (`feed_type_id` ASC),
  CONSTRAINT `fk_feed_adjustment_feed_type`
    FOREIGN KEY (`feed_type_id`)
    REFERENCES `livestock`.`feed_type` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

ALTER TABLE `feed`
  DROP `deleted`,
  DROP `updated`;

ALTER TABLE `feed` ADD `origin` VARCHAR(45) NULL AFTER `remarks`;
