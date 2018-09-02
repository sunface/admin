create database if not exists mafanr_admin;
USE  mafanr_admin;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(255) DEFAULT '' COMMENT '用户名，不可修改',
  `pw` varchar(255) DEFAULT '' COMMENT '用户密码',
  `priv` varchar(255) DEFAULT 'normal' COMMENT '全局用户权限',
  `create_date` timestamp NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

CREATE TABLE IF NOT EXISTS `service` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) NOT NULL COMMENT 'service名，不可修改',
  `owner` varchar(255) DEFAULT NULL COMMENT 'owner',
  `description` text COMMENT '介绍',
  `create_date` timestamp NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_service_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'servie表';

CREATE TABLE  IF NOT EXISTS `service_priv` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(255) DEFAULT NULL COMMENT '用户名',
  `service` varchar(255) DEFAULT NULL COMMENT 'service名',
  `priv` varchar(255) DEFAULT 'normal' COMMENT 'Service用户权限',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_privilege_username` (`username`,`service`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT 'Service权限表';


CREATE TABLE IF NOT EXISTS `application` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) NOT NULL COMMENT 'service名，不可修改',
  `service` varchar(255) DEFAULT NULL COMMENT 'service',
  `description` text COMMENT '介绍',
  `create_date` timestamp NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_application_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'application表';