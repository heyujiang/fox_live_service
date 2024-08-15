create table node (
    id int unsigned not null auto_increment primary key comment 'id',
    name varchar(255) not null default '' comment '节点名称',
    pid int unsigned not null  default 0 comment '父节点id',
    is_leaf tinyint unsigned not null default 0 comment '是否是叶子节点，0-否，1-是',
    sort int unsigned not null default 0 comment '排序',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment 'node';