CREATE TABLE `user_profile`
(
    `id`            bigint                                                       NOT NULL AUTO_INCREMENT,
    `nick`          varchar(255)                                                 NOT NULL DEFAULT '',
    `pwd`           varchar(128)                                                 NOT NULL DEFAULT '',
    `email`         varchar(128)                                                 NOT NULL DEFAULT '',
    `create_time`   datetime                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`   datetime                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `user_id`       varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
    `vip_time`      datetime                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_status`   tinyint                                                      NOT NULL DEFAULT '0' COMMENT '0正常 1禁用',
    `register_from` int                                                          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_email` (`email`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;