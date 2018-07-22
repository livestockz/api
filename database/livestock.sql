-- phpMyAdmin SQL Dump
-- version 4.5.1
-- http://www.phpmyadmin.net
--
-- Host: 127.0.0.1
-- Generation Time: Jul 19, 2018 at 10:44 AM
-- Server version: 10.1.9-MariaDB
-- PHP Version: 5.6.15

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `livestock`
--

-- --------------------------------------------------------

--
-- Table structure for table `feed`
--

CREATE TABLE `feed` (
  `id` char(36) NOT NULL,
  `feed_id` char(36) NOT NULL,
  `qty` decimal(10,2) NOT NULL,
  `remarks` varchar(255) NOT NULL,
  `reference` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `feed_type`
--

CREATE TABLE `feed_type` (
  `id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `unit` varchar(50) NOT NULL,
  `status` tinyint(1) NOT NULL,
  `deleted` tinyint(1) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_batch`
--

CREATE TABLE `growth_batch` (
  `id` char(36) NOT NULL,
  `name` varchar(45) NOT NULL,
  `status` tinyint(1) NOT NULL,
  `deleted` tinyint(1) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_batch_cycle`
--

CREATE TABLE `growth_batch_cycle` (
  `id` char(36) NOT NULL,
  `growth_batch_id` char(36) NOT NULL,
  `growth_pool_id` char(36) NOT NULL,
  `cycle_start` date NOT NULL,
  `cycle_finish` date DEFAULT NULL,
  `weight` decimal(10,2) NOT NULL,
  `amount` decimal(20,0) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_death`
--

CREATE TABLE `growth_death` (
  `id` char(36) NOT NULL,
  `growth_batch_cycle_id` char(36) NOT NULL,
  `death_date` date NOT NULL,
  `amount` decimal(20,0) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_feeding`
--

CREATE TABLE `growth_feeding` (
  `id` char(36) NOT NULL,
  `growth_batch_cycle_id` char(36) NOT NULL,
  `feed_date` date NOT NULL,
  `amount` decimal(5,2) NOT NULL,
  `feed_id` char(36) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_pool`
--

CREATE TABLE `growth_pool` (
  `id` char(36) NOT NULL,
  `name` varchar(45) NOT NULL,
  `status` char(32) NOT NULL,
  `deleted` tinyint(1) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_sales`
--

CREATE TABLE `growth_sales` (
  `id` char(36) NOT NULL,
  `sales_date` date NOT NULL,
  `weight` decimal(10,2) NOT NULL,
  `user_id` int(11) NOT NULL,
  `timestamp` datetime NOT NULL,
  `reference` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_sales_detail`
--

CREATE TABLE `growth_sales_detail` (
  `id` char(36) NOT NULL,
  `sales_id` char(36) NOT NULL,
  `growth_batch_cycle_id` char(36) NOT NULL,
  `amount` decimal(20,0) NOT NULL,
  `weight` decimal(10,2) NOT NULL,
  `detail_date` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `growth_summary`
--

CREATE TABLE `growth_summary` (
  `id` char(36) NOT NULL,
  `growth_cycle_batch_id` char(36) NOT NULL,
  `summary_date` date NOT NULL,
  `adg` decimal(5,2) NOT NULL,
  `fcr` int(4) NOT NULL,
  `sr` decimal(5,2) NOT NULL,
  `created` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `id` char(36) NOT NULL,
  `username` varchar(45) NOT NULL,
  `password` varchar(255) NOT NULL,
  `fullname` varchar(100) NOT NULL,
  `role` varchar(255) NOT NULL,
  `status` tinyint(1) NOT NULL,
  `deleted` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `feed`
--
ALTER TABLE `feed`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_feed_feed_type_idx` (`feed_id`);

--
-- Indexes for table `feed_type`
--
ALTER TABLE `feed_type`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `growth_batch`
--
ALTER TABLE `growth_batch`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `growth_batch_cycle`
--
ALTER TABLE `growth_batch_cycle`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_placement_batch_idx` (`growth_batch_id`),
  ADD KEY `fk_placement_pool_idx` (`growth_pool_id`);

--
-- Indexes for table `growth_death`
--
ALTER TABLE `growth_death`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_death_batch_cycle_idx` (`growth_batch_cycle_id`);

--
-- Indexes for table `growth_feeding`
--
ALTER TABLE `growth_feeding`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_feeding_batch_cycle_idx` (`growth_batch_cycle_id`),
  ADD KEY `fk_feeding_feed_idx` (`feed_id`);

--
-- Indexes for table `growth_pool`
--
ALTER TABLE `growth_pool`
  ADD PRIMARY KEY (`id`),
  ADD KEY `status` (`status`);

--
-- Indexes for table `growth_sales`
--
ALTER TABLE `growth_sales`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `growth_sales_detail`
--
ALTER TABLE `growth_sales_detail`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_sales_batch_cycle_idx` (`growth_batch_cycle_id`),
  ADD KEY `fk_sales_detail_sales_idx` (`sales_id`);

--
-- Indexes for table `growth_summary`
--
ALTER TABLE `growth_summary`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_summary_cycle_batch_idx` (`growth_cycle_batch_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);

--
-- Constraints for dumped tables
--

--
-- Constraints for table `feed`
--
ALTER TABLE `feed`
  ADD CONSTRAINT `fk_feed_feed_type` FOREIGN KEY (`feed_id`) REFERENCES `feed_type` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `growth_batch_cycle`
--
ALTER TABLE `growth_batch_cycle`
  ADD CONSTRAINT `fk_batch_cycle_batch` FOREIGN KEY (`growth_batch_id`) REFERENCES `growth_batch` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_batch_cycle_pool` FOREIGN KEY (`growth_pool_id`) REFERENCES `growth_pool` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `growth_death`
--
ALTER TABLE `growth_death`
  ADD CONSTRAINT `fk_death_batch_cycle` FOREIGN KEY (`growth_batch_cycle_id`) REFERENCES `growth_batch_cycle` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `growth_feeding`
--
ALTER TABLE `growth_feeding`
  ADD CONSTRAINT `fk_feeding_batch_cycle` FOREIGN KEY (`growth_batch_cycle_id`) REFERENCES `growth_batch_cycle` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_feeding_feed` FOREIGN KEY (`feed_id`) REFERENCES `feed` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `growth_sales_detail`
--
ALTER TABLE `growth_sales_detail`
  ADD CONSTRAINT `fk_sales_batch_cycle` FOREIGN KEY (`growth_batch_cycle_id`) REFERENCES `growth_batch_cycle` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_sales_detail_sales` FOREIGN KEY (`sales_id`) REFERENCES `growth_sales` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `growth_summary`
--
ALTER TABLE `growth_summary`
  ADD CONSTRAINT `fk_summary_cycle_batch` FOREIGN KEY (`growth_cycle_batch_id`) REFERENCES `growth_batch_cycle` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
