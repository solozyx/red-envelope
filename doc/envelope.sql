-- ----------------------------
-- Table structure for envelope_goods
-- ----------------------------

DROP TABLE IF EXISTS `red_envelope_goods`;
CREATE TABLE `red_envelope_goods`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `envelope_no` varchar(32) not null comment '红包编号',
    `envelope_type` tinyint(2) not null comment '包含普通红包，碰运气红包',
    `user_id` varchar(40) not null comment '用户编号，红包所属用户',
    `username` varchar(64) default '' not null comment '用户名称',
    `blessing` varchar(64) default '恭喜发财' not null comment '红包祝福语',
    `amount` decimal(30,6) unsigned not null default '0.000000' comment '红包总金额',
    `amount_one` decimal(30,6) unsigned not null default '0.000000' comment '单个红包金额，只属于普通平均分红包',
    `quantity` int(10) unsigned not null comment '红包总数量',
    `remain_amount` decimal(30,6) unsigned not null default '0.000000' comment '红包剩余金额',
    `remain_quantity` int(10) unsigned not null comment '红包剩余数量',
    `expired_at` datetime(3) not null comment '过期时间',
    `status` tinyint(2) not null comment '红包/订单状态：0创建，1发布启用，2过期，3失效',
    `order_type` tinyint(2) not null comment '订单类型，发布单，退款单',
    `pay_status` tinyint(2) not null comment '支付状态:未支付，支付中，已支付',
    `created_at` datetime(3) not null default current_timestamp(3) comment '创建时间',
    `updated_at` datetime(3) not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    `origin_envelope_no` varchar(32) not null default '' comment '原红包编号',
    primary key (`id`) using btree ,
    unique key `envelope_no_idx` (`envelope_no`) using btree ,
    key `id_user_idx` (`user_id`) using btree
) engine = InnoDB AUTO_INCREMENT=171 DEFAULT charset=utf8 ROW_FORMAT=DYNAMIC;


-- ----------------------------
-- Table structure for envelope_item
-- ----------------------------

DROP TABLE IF EXISTS `red_envelope_item`;
CREATE TABLE `red_envelope_item`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `item_no` varchar(32) not null comment '红包订单详情编号',
    `envelope_no` varchar(32) not null comment '红包编号',
    `recv_user_id` varchar(40) not null comment '红包接受者用户编号',
    `recv_username` varchar(64) default '' not null comment '红包接受者用户名称',
    `amount` decimal(30,6) unsigned not null default '0.000000' comment '收到红包金额',
    `quantity` int(10) unsigned not null comment '大吼道红包数量',
    `remain_amount` decimal(30,6) unsigned not null default '0.000000' comment '收到后原红包剩余金额',
    `account_no` varchar(32) not null comment '红包接受者账户编号',
    `pay_status` tinyint(2) not null comment '支付状态:未支付，支付中，已支付',
    `desc` varchar(128) not null comment '交易描述',
    `created_at` datetime(3) not null default current_timestamp(3) comment '创建时间',
    `updated_at` datetime(3) not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`) using btree ,
    unique key `item_no_idx` (`item_no`) using btree
) engine = InnoDB AUTO_INCREMENT=171 DEFAULT charset=utf8 ROW_FORMAT=DYNAMIC;