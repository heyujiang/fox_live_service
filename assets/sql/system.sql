create table rule (
    id int unsigned primary key auto_increment ,
    title varchar(100) not null default '' comment '菜单名称',
    route_name varchar(100) not null default '' comment '路由名称',
    route_path varchar(100) not null default '' comment '路由地址',
    component varchar(100) not null default '' comment '组件路径',
    redirect varchar(100) not null default '' comment '重定向地址',
    locale varchar(100) not null default '' comment '多语言标识',
    icon varchar(20) not   null default  '' comment '显示图标',
    permission varchar(50) not null default '' comment '权限标识',
    spacer varchar(100) not null default '' comment 'spacer',
    hide_in_menu tinyint unsigned not null default 0 comment '菜单隐藏，是否在左侧菜单栏中隐藏该项',
    is_ext tinyint unsigned not null default 0 comment '是否外链',
    no_affix tinyint unsigned not null default 0 comment '不添加tab中 （如果设置为true标签将不会添加到tab-bar中）',
    keepalive tinyint unsigned not null default 0 comment '是否缓存（在页签模式生效,页面name和路由name保持一致）',
    requires_auth tinyint unsigned not null default 0 comment '是否登录鉴权（是否需要登录鉴权）',
    only_page tinyint unsigned not null default 0 comment '独立页面 （不需layout和登录，如登录页、数据大屏）',
    active_menu tinyint unsigned not null default 0 comment '高亮菜单 （高亮设置的菜单项）',
    hide_children_in_menu tinyint unsigned not null default 0 comment '显示单项 （hideChildrenInMenu强制在左侧菜单中显示单项）',
    status tinyint unsigned not null default 0 comment '状态',
    `order` tinyint unsigned not null default 0 comment '排序',
    type tinyint unsigned not null default 0 comment '菜单类型 （菜单目录 菜单 按钮）',
    pid int unsigned not null default 0 comment '上级菜单id',
    created_id int unsigned not null default 0 comment '',
    updated_id int unsigned not null default 0 comment '',
    created_at timestamp not null default current_timestamp comment '创建时间',
    updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    key `idx_pid` (`pid`)
)  engine = innodb comment 'rule';


drop table role;

create table role (
      id int unsigned primary key auto_increment ,
      title varchar(100) not null default '' comment '菜单名称',
      remark varchar(255) not null default '' comment '备注',
      rule_ids varchar(255) not null default '' comment '菜单id',
      status tinyint unsigned not null default 0 comment '状态',
      pid int unsigned not null default 0 comment 'pid',
      created_id int unsigned not null default 0 comment '',
      updated_id int unsigned not null default 0 comment '',
      created_at timestamp not null default current_timestamp comment '创建时间',
      updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
      key `idx_pid` (`pid`)
)  engine = innodb comment 'rule';

create table dept (
      id int unsigned primary key auto_increment ,
      title varchar(100) not null default '' comment '菜单名称',
      remark varchar(255) not null default '' comment '备注',
      status tinyint unsigned not null default 0 comment '状态',
      pid int unsigned not null default 0 comment 'pid',
     `order` int unsigned not null default 0 comment '排序',
      created_id int unsigned not null default 0 comment '',
      updated_id int unsigned not null default 0 comment '',
      created_at timestamp not null default current_timestamp comment '创建时间',
      updated_at timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
      key `idx_pid` (`pid`)
)  engine = innodb comment '部门';

alter table dept add column title varchar(100) not null default '' comment '菜单名称';