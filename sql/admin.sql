create database if not exists mafanr_admin;
USE  mafanr_admin;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(255) DEFAULT '' COMMENT '用户名，不可修改',
  `pw` varchar(255) DEFAULT '' COMMENT '用户密码',
  `priv` varchar(255) DEFAULT 'normal' COMMENT '全局用户权限',
  `create_date` varchar(255) NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

CREATE TABLE IF NOT EXISTS `service` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) NOT NULL COMMENT 'service名，不可修改',
  `owner` varchar(255) DEFAULT NULL COMMENT 'owner',
  `description` text COMMENT '介绍',
  `create_date` varchar(255) NOT NULL COMMENT '记录创建时间',
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
  `create_date` varchar(255) NOT NULL COMMENT '记录创建时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_application_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'application表';


CREATE TABLE IF NOT EXISTS `cloud_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `public_ip` varchar(255) COMMENT '公网IP',
  `private_ip` varchar(255) COMMENT '内网IP',
  `cloud_provider` varchar(255) COMMENT '云服务提供商，1:私有云 2:阿里云 ',
  `ssh_user` varchar(20) DEFAULT 'root' COMMENT 'ssh登录账号',
  `ssh_pw` varchar(40) DEFAULT 'golang' COMMENT 'ssh登录密码',
  `region` varchar(100) DEFAULT '' COMMENT '服务器所属可用区', 
  `configure` varchar(255) DEFAULT ''  COMMENT '机器配置、网络配置',
  `expire` varchar(255) DEFAULT '' COMMENT '机器过期时间',
  `modify_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Index_cloud_server_private_ip` (`private_ip`,`cloud_provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '云服务器表';