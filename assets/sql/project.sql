create table project (
    id int unsigned not null auto_increment primary key comment 'id',
    name varchar(255) not null default '' comment '节点名称',
    description varchar(2000) not null default '' comment '项目概览',
    attr tinyint unsigned not null default 0 comment '项目属性，1-集中式，2-分布式，3-分散式',
    state tinyint unsigned not null default 0 comment '状态，1-待定，2-推荐，3-终止，4-已完成',
    type tinyint unsigned not null default 0 comment '类型，1-风电，2-光伏，3-风电+光伏，4-储能，5-风电+储能，6-光伏+储能，7-风光储一体',
    node_id int unsigned not null default 0 comment '当前项目节点',
    node_name varchar(255) not null default '' comment '当前节点名称',
    schedule decimal(4,2) not null default 0 comment '项目当前进度',
    capacity decimal(8,2) not null default 0 comment '容量大小',
    properties varchar(255) not null default '' comment '土地性质',
    area decimal(8,2) not null default  0 comment '土地面积',
    address varchar(255) not null default '' comment '项目地址',
    connect varchar(1000) not null default '' comment '电网接入情况',
    star tinyint unsigned not null default 0 comment '星级',
    user_id int unsigned not null default 0 comment '第一负责人id',
    username varchar(100)  not null default '' comment '第一负责人',
    investment_agreement varchar(1000) not null default '' comment '投资协议',
    business_condition varchar(1000) not null default '' comment '商务条件',
    begin_time timestamp null comment '开始时间',
    is_deleted tinyint unsigned not null  default 0 comment '是否删除，0否，1是',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间'
) engine = innodb comment 'project';


create table project_node (
     id int unsigned not null auto_increment primary key comment 'id',
     project_id int unsigned not null default 0 comment '项目id',
     p_id int unsigned not null default 0 comment '父节点id',
     node_id int unsigned not null default 0 comment '节点id',
     name varchar(255) not null default '' comment '当前节点名称',
     is_leaf tinyint unsigned not null default 0 comment '是否是叶子节点，0-否，1-是',
     sort int unsigned not null default 0 comment '排序',
     state tinyint unsigned not null default 0 comment '状态，1-未开始，2-进行中，3-已完成',
     created_id int unsigned not null default 0 comment '',
     updated_id int unsigned not null default 0 comment '',
     created_at timestamp not null default current_timestamp comment '创建时间',
     updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    unique key `uk_project_id_node_id` (`project_id`,`node_id`)
) engine = innodb comment '项目节点';

create table project_person (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    user_id int unsigned not null default 0 comment '用户id',
    name varchar(100) not null default 0 comment '用户',
    phone_number varchar(11) not null default 0 comment '手机号',
    type tinyint unsigned not null default 0 comment '成员类型:1-第一负责人，2-第二负责人，3-普通成员',
    created_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    UNIQUE KEY `uk_project_id_user_id` (`project_id`,`user_id`)
) engine = innodb comment '项目成员负责人';

create table project_record (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    project_name varchar(200) not null default '' comment '项目名称',
    node_id int unsigned not null default 0 comment '节点id',
    node_name varchar(200) not null default '' comment '节点名称',
    user_id int unsigned not null default 0 comment '用户id',
    username varchar(200) not null default '' comment '用户名称',
    overview varchar(1000) not null default '' comment '概况',
    state tinyint unsigned not null default 0 comment '提交的节点状态',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    key `idx_project_id_node_id` (`project_id`,`node_id`)
) engine = innodb comment '项目节点记录';

create table project_attached (
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    node_id int unsigned not null default 0 comment '节点id',
    record_id int unsigned not null default 0 comment '节点记录id',
    user_id int unsigned not null default 0 comment '用户id',
    url varchar(255) not null default '' comment '附件地址',
    filename varchar(255) not null default '' comment '附件名称',
    mime varchar(100) not null default  '' comment '附件格式',
    size int unsigned not null default 0 comment '文件大小',
    created_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    KEY  `idx_record_id` (`record_id`)
) engine = innodb comment '项目节点记录附件';


-- 电话
create table project_contact(
    id int unsigned not null auto_increment primary key comment 'id',
    project_id int unsigned not null default 0 comment '项目id',
    name varchar(100) not null default 0 comment '姓名',
    phone_number varchar(20) not null default 0 comment '电话',
    type tinyint unsigned not null default 0 comment '联系人类型:1-第一负责人，2-第二负责人，3-普通成员',
    description varchar(500) not null default '' comment '描述备注',
    created_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    KEY `idk_project_id` (`project_id`)
) engine = innodb comment '项目联系人';
