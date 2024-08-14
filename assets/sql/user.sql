create table user (
  id int unsigned not null auto_increment primary key ,
  username varchar(100) not null default '' comment 'username',
  password char(32) not null default '' comment 'pwd',
  phone_number char(11) not null default '' comment '',
  email varchar(30) not null default '' comment '',
  name varchar(200) not null default '' comment '名字',
  nick_name varchar(30) not null default '' comment '昵称',
  avatar varchar(200) not null default '' comment '头像',
  state tinyint not null default 0 comment '状态,1-启用，2-禁用',
  role_ids varchar(255) not null default '' comment '用户角色 id 用,分割' ,
  dept_id int unsigned not null default 0 comment '部门id',
  created_id int unsigned not null default 0 comment '',
  updated_id int unsigned not null default 0 comment '',
  created_at timestamp not null default current_timestamp comment '创建时间',
  updated_at timestamp not null default current_timestamp on update current_timestamp comment '编辑时间',
  unique key `idx_username` (username),
  unique key `uk_phone_number` (`phone_number`)
) engine = innodb comment 'USER';
