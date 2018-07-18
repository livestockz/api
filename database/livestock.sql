-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema livestock
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `livestock` ;

-- -----------------------------------------------------
-- Schema livestock
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `livestock` DEFAULT CHARACTER SET latin1 ;
USE `livestock` ;

-- -----------------------------------------------------
-- Table `livestock`.`feed_type`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`feed_type` ;

CREATE TABLE IF NOT EXISTS `livestock`.`feed_type` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `unit` VARCHAR(50) NOT NULL,
  `status` TINYINT(1) NOT NULL,
  `deleted` TINYINT(1) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`feed`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`feed` ;

CREATE TABLE IF NOT EXISTS `livestock`.`feed` (
  `id` CHAR(36) NOT NULL,
  `feed_id` CHAR(36) NOT NULL,
  `qty` DECIMAL(10,2) NOT NULL,
  `remarks` VARCHAR(255) NOT NULL,
  `reference` VARCHAR(255) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_feed_feed_type_idx` (`feed_id` ASC),
  CONSTRAINT `fk_feed_feed_type`
    FOREIGN KEY (`feed_id`)
    REFERENCES `livestock`.`feed_type` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_batch`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_batch` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_batch` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `status` INT(4) NOT NULL,
  `deleted` TINYINT(1) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_pool`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_pool` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_pool` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `status` INT(4) NOT NULL,
  `deleted` TINYINT(1) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_batch_cycle`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_batch_cycle` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_batch_cycle` (
  `id` CHAR(36) NOT NULL,
  `growth_batch_id` CHAR(36) NOT NULL,
  `growth_pool_id` CHAR(36) NOT NULL,
  `cycle_start` DATE NOT NULL,
  `cycle_finish` DATE NULL DEFAULT NULL,
  `weight` DECIMAL(10,2) NOT NULL,
  `amount` DECIMAL(20,0) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_placement_batch_idx` (`growth_batch_id` ASC),
  INDEX `fk_placement_pool_idx` (`growth_pool_id` ASC),
  CONSTRAINT `fk_batch_cycle_batch`
    FOREIGN KEY (`growth_batch_id`)
    REFERENCES `livestock`.`growth_batch` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_batch_cycle_pool`
    FOREIGN KEY (`growth_pool_id`)
    REFERENCES `livestock`.`growth_pool` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_death`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_death` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_death` (
  `id` CHAR(36) NOT NULL,
  `growth_batch_cycle_id` CHAR(36) NOT NULL,
  `death_date` DATE NOT NULL,
  `amount` DECIMAL(20,0) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_death_batch_cycle_idx` (`growth_batch_cycle_id` ASC),
  CONSTRAINT `fk_death_batch_cycle`
    FOREIGN KEY (`growth_batch_cycle_id`)
    REFERENCES `livestock`.`growth_batch_cycle` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_feeding`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_feeding` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_feeding` (
  `id` CHAR(36) NOT NULL,
  `growth_batch_cycle_id` CHAR(36) NOT NULL,
  `feed_date` DATE NOT NULL,
  `amount` DECIMAL(5,2) NOT NULL,
  `feed_id` CHAR(36) NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_feeding_batch_cycle_idx` (`growth_batch_cycle_id` ASC),
  INDEX `fk_feeding_feed_idx` (`feed_id` ASC),
  CONSTRAINT `fk_feeding_batch_cycle`
    FOREIGN KEY (`growth_batch_cycle_id`)
    REFERENCES `livestock`.`growth_batch_cycle` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_feeding_feed`
    FOREIGN KEY (`feed_id`)
    REFERENCES `livestock`.`feed` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_sales`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_sales` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_sales` (
  `id` CHAR(36) NOT NULL,
  `sales_date` DATE NOT NULL,
  `weight` DECIMAL(10,2) NOT NULL,
  `user_id` INT(11) NOT NULL,
  `timestamp` DATETIME NOT NULL,
  `reference` VARCHAR(255) NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_sales_detail`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_sales_detail` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_sales_detail` (
  `id` CHAR(36) NOT NULL,
  `sales_id` CHAR(36) NOT NULL,
  `growth_batch_cycle_id` CHAR(36) NOT NULL,
  `amount` DECIMAL(20,0) NOT NULL,
  `weight` DECIMAL(10,2) NOT NULL,
  `detail_date` DATE NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_sales_batch_cycle_idx` (`growth_batch_cycle_id` ASC),
  INDEX `fk_sales_detail_sales_idx` (`sales_id` ASC),
  CONSTRAINT `fk_sales_batch_cycle`
    FOREIGN KEY (`growth_batch_cycle_id`)
    REFERENCES `livestock`.`growth_batch_cycle` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_sales_detail_sales`
    FOREIGN KEY (`sales_id`)
    REFERENCES `livestock`.`growth_sales` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`growth_summary`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`growth_summary` ;

CREATE TABLE IF NOT EXISTS `livestock`.`growth_summary` (
  `id` CHAR(36) NOT NULL,
  `growth_cycle_batch_id` CHAR(36) NOT NULL,
  `summary_date` DATE NOT NULL,
  `adg` DECIMAL(5,2) NOT NULL,
  `fcr` INT(4) NOT NULL,
  `sr` DECIMAL(5,2) NOT NULL,
  `created` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_summary_cycle_batch_idx` (`growth_cycle_batch_id` ASC),
  CONSTRAINT `fk_summary_cycle_batch`
    FOREIGN KEY (`growth_cycle_batch_id`)
    REFERENCES `livestock`.`growth_batch_cycle` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `livestock`.`user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `livestock`.`user` ;

CREATE TABLE IF NOT EXISTS `livestock`.`user` (
  `id` CHAR(36) NOT NULL,
  `username` VARCHAR(45) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `fullname` VARCHAR(100) NOT NULL,
  `role` VARCHAR(255) NOT NULL,
  `status` TINYINT(1) NOT NULL,
  `deleted` TINYINT(1) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
