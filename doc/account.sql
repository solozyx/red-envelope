-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP  TABLE IF EXISTS `account`;
CREATE TABLE `account`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '账户ID',
    `account_no` varchar(32) NOT NULL COMMENT '账户编号,账户唯一标识',
    `account_name` varchar(64) NOT NULL COMMENT '账户名称,用来说明账户的简短描述,账户对应的名称或命名,比如xxx积分,xxx零钱',
    `account_type` tinyint(2) NOT NULL COMMENT '账户类型，用来区分不同的账户：积分，会员，钱包，红包',
    `currency_code` char(3) not null default 'CNY' comment '货币类型：CNY人民币，EUR欧元，USD美元。。。',
    `user_id` varchar(40) not null comment '用户编号，账户所属用户',
    `username` varchar(64) default '' not null comment '用户名称',
    `balance` decimal(30,6) unsigned not null default '0.000000' comment '账户可用余额',
    `status` tinyint(2) not null comment '账户状态：0初始化，1启用，2停用',
    `created_at` datetime(3) not null default current_timestamp(3) comment '创建时间',
    `updated_at` datetime(3) not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`) using btree ,
    unique key `account_no_idx` (`account_no`) using btree,
    key `id_user_idx` (`user_id`) using btree
)engine = InnoDB AUTO_INCREMENT=171 DEFAULT charset=utf8 ROW_FORMAT=DYNAMIC;

INSERT INTO `account` (account_no, account_name, account_type, user_id, username, status)
    values ('10000020190101010000000000000001','系统红包账户',1,'10001','系统红包账户',1);

-- ----------------------------
-- Table structure for account_log
-- ----------------------------
DROP TABLE IF EXISTS `account_log`;
create table `account_log`
(
    `id` bigint(20) NOT NULL auto_increment,
    `trade_no` varchar(32) not null comment '交易单号 全局不重复字母或数字 唯一性标识',
    `log_no` varchar(32) not null comment '流水编号 全局不重复字母或数字 唯一性标识',
    `account_no` varchar(32) NOT NULL COMMENT '账户编号',
    `target_account_no` varchar(32) NOT NULL COMMENT '目标账户编号',
    `user_id` varchar(40) not null comment '用户编号，账户所属用户',
    `username` varchar(64) default '' not null comment '用户名称',
    `target_user_id` varchar(40) not null comment '目标用户编号，账户所属用户',
    `target_username` varchar(64) default '' not null comment '目标用户名称',
    `amount` decimal(30,6) unsigned not null default '0.000000' comment '交易金额',
    `balance` decimal(30,6) unsigned not null default '0.000000' comment '该交易后的余额',
    `change_type` tinyint(2) not null default '0' comment '流水交易类型，0创建账户，>0为收入类型，<0为支出类型',
    `change_flag` tinyint(2) not null default '0' comment '交易变化标识，-1出账，1进账，枚举',
    `status` tinyint(2) not null default '0' comment '交易状态',
    `desc` varchar(128) not null comment '交易描述',
    `created_at` datetime(3) not null default current_timestamp(3) comment '创建时间',
    primary key (`id`) using btree,
    unique key `id_log_no_idx` (`log_no`) using btree,
    key `id_user_idx` (`user_id`) using btree,
    key `id_account_idx` (`account_no`) using btree,
    key `id_trade_idx` (`trade_no`) using btree
)engine = InnoDB AUTO_INCREMENT=171 DEFAULT charset=utf8 ROW_FORMAT=DYNAMIC;

set foreign_key_checks = 1;