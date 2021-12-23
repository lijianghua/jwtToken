CREATE TABLE `tbl_user` (
                            `id` INT(11) NOT NULL AUTO_INCREMENT,
                            `user_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名' COLLATE 'utf8mb4_general_ci',
                            `user_pwd` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码' COLLATE 'utf8mb4_general_ci',
                            `email` VARCHAR(64) NULL DEFAULT '' COMMENT '邮箱' COLLATE 'utf8mb4_general_ci',
                            `phone` VARCHAR(128) NULL DEFAULT '' COMMENT '手机号' COLLATE 'utf8mb4_general_ci',
                            `email_validated` TINYINT(1) NULL DEFAULT '0' COMMENT '邮箱是否已验证',
                            `phone_validated` TINYINT(1) NULL DEFAULT '0' COMMENT '手机号是否已验证',
                            `signup_at` DATETIME NULL DEFAULT current_timestamp() COMMENT '注册日期',
                            `last_active` DATETIME NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '最后活跃时间戳',
                            `profile` TEXT NULL DEFAULT NULL COMMENT '用户属性' COLLATE 'utf8mb4_general_ci',
                            `status` INT(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
                            PRIMARY KEY (`id`) USING BTREE,
                            UNIQUE INDEX `idx_username` (`user_name`) USING BTREE,
                            INDEX `idx_status` (`status`) USING BTREE
)
    COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=35
;