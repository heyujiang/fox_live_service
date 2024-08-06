create table file (
      id int unsigned not null auto_increment primary key comment 'id',
      type varchar(30) not null default '' comment '文件名称',
      url varchar(100)  not null  default '' comment '文件路径',
      mime varchar(100) not null  default '' comment '文件mime',
      size int not null default 0 comment '文件大小',
      filename varchar(100) not null  default '' comment '文件名称',
      created_id int unsigned not null default 0 comment '',
      created_at timestamp not null default current_timestamp comment '创建时间'
) engine = innodb comment 'file上传记录';


