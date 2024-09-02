/*
 Navicat Premium Data Transfer

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80300 (8.3.0)
 Source Host           : localhost:3306
 Source Schema         : yuka

 Target Server Type    : MySQL
 Target Server Version : 80300 (8.3.0)
 File Encoding         : 65001

 Date: 02/09/2024 10:20:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for yu_user
-- ----------------------------
DROP TABLE IF EXISTS `yu_user`;
CREATE TABLE `yu_user` (
                           `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
                           `wallet` varchar(42) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'eth 钱包地址',
                           `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '账号状态',
                           `type` int NOT NULL DEFAULT '0' COMMENT '账号类型 0-用户 1-商户 2-管理员',
                           `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
                           `login_time` TIMESTAMP  COMMENT '登录时间',
                           `login_ip` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录ip',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `unique_wallet` (`wallet`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

SET FOREIGN_KEY_CHECKS = 1;
