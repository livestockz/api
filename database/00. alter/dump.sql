ALTER TABLE `growth_death` ADD `weight` DECIMAL(10,2) NOT NULL AFTER `death_date`;
ALTER TABLE `growth_feeding` ADD `feeding_date` DATE NOT NULL AFTER `feed_type_id`;
ALTER TABLE `growth_summary` ADD `weight` DECIMAL(10,2) NOT NULL AFTER `summary_date`, ADD `amount` DECIMAL(20,0) NOT NULL AFTER `weight`;
ALTER TABLE `growth_summary` CHANGE `growth_cycle_batch_id` `growth_batch_cycle_id` CHAR(36) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL;