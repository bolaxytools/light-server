/*
 Navicat Premium Data Transfer

 Source Server         : hsy
 Source Server Type    : MySQL
 Source Server Version : 50724
 Source Host           : 122.51.106.29:3306
 Source Schema         : wallet

 Target Server Type    : MySQL
 Target Server Version : 50724
 File Encoding         : 65001

 Date: 02/12/2019 18:33:25
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address` (
  `addr` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '地址',
  `add_time` bigint(20) NOT NULL COMMENT '添加时间',
  `update_time` bigint(20) NOT NULL COMMENT '最后更新时间',
  PRIMARY KEY (`addr`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for block
-- ----------------------------
DROP TABLE IF EXISTS `block`;
CREATE TABLE `block` (
  `height` bigint(20) NOT NULL COMMENT '区块高度',
  `hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '区块hash',
  `tx_count` int(11) NOT NULL COMMENT '交易数',
  `block_time` bigint(20) NOT NULL COMMENT '区块时间',
  PRIMARY KEY (`height`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for block_signer
-- ----------------------------
DROP TABLE IF EXISTS `block_signer`;
CREATE TABLE `block_signer` (
  `b_height` bigint(20) NOT NULL COMMENT '区块高度',
  `signer_address` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '签名者地址',
  PRIMARY KEY (`b_height`,`signer_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
  `contract` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合约地址',
  `wallet` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '钱包地址',
  `balance` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '余额',
  PRIMARY KEY (`contract`,`wallet`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
  `contract` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '智能合约地址',
  `symbol` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币名称',
  `logo` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币标志',
  `desc` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '简述',
  PRIMARY KEY (`contract`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for tx
-- ----------------------------
DROP TABLE IF EXISTS `tx`;
CREATE TABLE `tx` (
  `tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交易hash',
  `tx_type` smallint(255) NOT NULL DEFAULT '0' COMMENT '交易类型',
  `addr_from` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '转账者',
  `addr_to` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收款者',
  `amount` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '金额',
  `miner_fee` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '手续费',
  `block_height` bigint(255) NOT NULL COMMENT '块高',
  `tx_time` bigint(20) NOT NULL COMMENT '交易时间',
  `memo` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交易备注',
  PRIMARY KEY (`tx_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
