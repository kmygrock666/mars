DROP DATABASE IF EXISTS `myjapan`;

CREATE DATABASE IF NOT EXISTS `myjapan` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `myjapan`;

CREATE TABLE IF NOT EXISTS `games`(
    `game_id` bigint(20) NOT NULL COMMENT '遊戲編號',
    `state` tinyint(1) NOT NULL COMMENT '狀態，0:新局, 1:未完成, 2:完成',
    `user_id` int(10) UNSIGNED NOT NULL COMMENT '會員',
    `updated_at` timestamp NOT NULL COMMENT '更新時間戳記',
    `created_at` int(10) NOT NULL COMMENT '建立時間戳記'
)   ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `orders`(
    `order_id` bigint(20) NOT NULL COMMENT '注單編號',
    `result` int(3) NOT NULL COMMENT '賽果編號',
    `game_id` bigint(20) NOT NULL COMMENT '遊戲編號',
    `state` tinyint(1) NOT NULL COMMENT '狀態，0:未作答, 1:錯誤, 2:正確',
    `updated_at` timestamp NOT NULL COMMENT '更新時間戳記',
    `created_at` int(10) NOT NULL COMMENT '建立時間戳記'
)   ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `word_list` (
    `id` bigint(20) NOT NULL ,
    `foreign` text NOT NULL COMMENT '外語',
    `native` text NOT NULL COMMENT '母語',
    `pinyin` text NOT NULL COMMENT '拼音',
    `extra_info` text NOT NULL COMMENT '額外資訊',
    `updated_at` timestamp NOT NULL COMMENT '更新時間戳記',
    `created_at` int(10) NOT NULL COMMENT '建立時間戳記'
)   ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `games` ADD PRIMARY KEY (`game_id`);
ALTER TABLE `games` ADD UNIQUE `user_id` (`user_id`);

ALTER TABLE `orders` ADD PRIMARY KEY (`order_id`);
ALTER TABLE `orders` ADD UNIQUE `game_id` (`game_id`);

ALTER TABLE `word_list` ADD PRIMARY KEY (`id`);