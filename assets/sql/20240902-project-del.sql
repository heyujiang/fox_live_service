alter table project_person add column is_deleted tinyint(1) unsigned not null default 0 comment '是否删除，0否，1是' after `type`;
alter table project_record add column is_deleted tinyint(1) unsigned not null default 0 comment '是否删除，0否，1是' after `state`;
alter table project_attached add column is_deleted tinyint(1) unsigned not null default 0 comment '是否删除，0否，1是' after `size`;

update project_record set is_deleted = 1 where project_id not in (select id from project where project.is_deleted = 0);
update project_person set is_deleted = 1 where project_id not in (select id from project where project.is_deleted = 0);
update project_attached set is_deleted = 1 where project_id not in (select id from project where project.is_deleted = 0);