create database `test` default character set utf8mb4;

use `test`;

create table user_info (
 id int not null auto_increment,
 `name` varchar(20),
 `title` varchar(20),
 money int not null,
 primary key(id)

) engine = InnoDB charset = utf8mb4;

insert into user_info (`name`, `title`,`money`)values('milo', 'title1','5');

select * from user_info;

update user_info set `name`='milo2' where id = 1;


create view view_user_info
as

select id,`name`, `title` from user_info;

select * from view_user_info;



delimiter //
create procedure in_param(in p_in int)
begin
	select p_in;
    set p_in=2;
    select p_in;
end
//
delimiter ;

set @p_in = 1;
call in_param(@p_in);


delimiter //
create procedure out_param(out p_out int)
begin
	select p_out;
    set p_out=2;
    select p_out;
end
//
delimiter ;

set @p_out = 1;
call out_param(@p_out);
select @p_out;



create table user_attr (
	id int not null auto_increment,
    user_id int,
    score int,
    primary key(id)
) engine = InnoDB;


delimiter //
create trigger tri_user after insert
on user_info for each row
begin
	declare userId int;
    set userId = New.id;
    insert into user_attr(user_id,score)values(userId, 0);
end
//
delimiter ;


select * from user_attr;
select * from user_info;


insert into user_info(`name`,title,money)values('milo', '', 0);

SET SESSION tx_isolation='READ-UNCOMMITTED';



#定期备份

set global event_scheduler = ON;
show variables like 'event%';


create table t2(timeline timestamp);


create event event_test_insert
on schedule every 1 second
do insert into t2 values(CURRENT_timestamp);

select * from t2;



# select user,host from mysql.user;


CREATE USER 'milo'@'localhost' IDENTIFIED BY 'milo';

set SQL_SAFE_UPDATES = 1;
# update mysql.user set host='%' where user = 'milo';



GRANT ALL PRIVILEGES ON *.* TO 'milo'@'%' ;
FLUSH PRIVILEGES;


#是做同步  mysql -> redis  
#          mysql  -> es
