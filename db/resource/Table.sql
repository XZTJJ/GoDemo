CREATE TABLE `finance_pay_info` (
	`pay_id` bigint NOT NULL AUTO_INCREMENT COMMENT '支付ID',
	`uid` INT DEFAULT NULL COMMENT '用户ID',
	`open_id` VARCHAR ( 32 ) DEFAULT NULL COMMENT '用户openID',
	`pay_amount` DECIMAL ( 10, 2 ) DEFAULT NULL COMMENT '支付金额',
	`charge_channel` tinyint DEFAULT NULL COMMENT '充值渠道：APP支付宝、APP微信、微信H5、支付宝小程序（charge_channel）',
	`status` INT DEFAULT NULL COMMENT '状态（pay_status）',
	`pay_time` datetime DEFAULT NULL COMMENT '支付时间',
	`refunded_amount` DECIMAL ( 10, 2 ) DEFAULT NULL COMMENT '已退款金额',
	`create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY ( `pay_id` ) USING BTREE 
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COMMENT = '支付单信息';