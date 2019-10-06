DROP  TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '账户ID',
    `username` varchar(64) NOT NULL COMMENT '账户名称,用来说明账户的简短描述,账户对应的名称或命名,比如xxx积分,xxx零钱',
    `mobile` varchar(13) NOT NULL COMMENT '手机号, unique key',
    `created_at` datetime(3) not null default current_timestamp(3) comment '创建时间',
    `updated_at` datetime(3) not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`) using btree ,
    unique key `mobile_idx` (`mobile`) using btree
)engine = InnoDB AUTO_INCREMENT=171 DEFAULT charset=utf8 ROW_FORMAT=DYNAMIC;