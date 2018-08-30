create database if not exists mafanr_admin;
USE  mafanr_admin;
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(255) DEFAULT '' COMMENT '用户名',
  `pw` varchar(255) DEFAULT '' COMMENT '用户密码',
  `priv` varchar(255) DEFAULT 'normal' COMMENT '用户权限',
  `create_date` timestamp NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';