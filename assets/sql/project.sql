create table project (
    id int unsigned not null auto_increment primary key comment 'id',
    name varchar(255) not null default '' comment '节点名称',
    node_id int unsigned not null default 0 comment '当前项目节点',
    description varchar(2000) not null default '' comment '项目概览',
    attr tinyint unsigned not null default 0 comment '项目属性，1-集中式，2-分布式，3-分散式',
    state tinyint unsigned not null default 0 comment '状态，1-待定，2-推荐，3-终止，4-已完成',
    type tinyint unsigned not null default 0 comment '类型，1-风电，2-光伏，3-风电+光伏，4-储能，5-风电+储能，6-光伏+储能，7-风光储一体',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment 'project';


create table project_node (
     id int unsigned not null auto_increment primary key comment 'id',
     project_id int unsigned not null default 0 comment '项目id',
     node_id int unsigned not null default 0 comment '节点id',
     state tinyint unsigned not null default 0 comment '状态，1-未开始，2-进行中，3-已完成',
     created_id int unsigned not null default 0 comment '',
     updated_id int unsigned not null default 0 comment '',
     created_at timestamp not null default current_timestamp comment '创建时间',
     updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment '项目节点';

create table project_person (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    user_id int unsigned not null default 0 comment '用户id',
    type tinyint unsigned not null default 0 comment '成员类型:1-第一负责人，2-第二负责人，3-普通成员',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment '项目成员负责人';

create table project_record (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    node_id int unsigned not null default 0 comment '节点id',
    user_id int unsigned not null default 0 comment '用户id',
    overview varchar(1000) not null default '' comment '概况',
    state tinyint unsigned not null default 0 comment '状态',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment '项目记录';


create table project_attached (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    node_id int unsigned not null default 0 comment '节点id',
    record_id int unsigned not null default 0 comment '节点记录id',
    user_id int unsigned not null default 0 comment '用户id',
    file_url varchar(255) not null default '' comment '附件地址',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment '项目附件';