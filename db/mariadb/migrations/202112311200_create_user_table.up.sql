
CREATE TABLE `user` (
`id` CHAR(36) NOT NULL COLLATE 'utf8mb4_general_ci',
`user_name` VARCHAR(64) NOT NULL COMMENT '用户名' COLLATE 'utf8mb4_general_ci',
`user_pwd` VARCHAR(255) NOT NULL COMMENT '用户密码hash值' COLLATE 'utf8mb4_general_ci',
`created_at` DATETIME NOT NULL DEFAULT current_timestamp() COMMENT '注册日期',
`updated_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '修改日期',
`last_active` DATETIME NULL DEFAULT NULL COMMENT '最后活跃时间',
`email` VARCHAR(255) NULL DEFAULT NULL COMMENT '邮箱' COLLATE 'utf8mb4_general_ci',
`phone` VARCHAR(100) NULL DEFAULT NULL COMMENT '电话' COLLATE 'utf8mb4_general_ci',
`status` INT(4) NOT NULL DEFAULT '1' COMMENT '账户状态(启用/禁用)',
PRIMARY KEY (`id`) USING BTREE,
UNIQUE INDEX `idx_username` (`user_name`) USING BTREE,
INDEX `idx_status` (`status`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB;