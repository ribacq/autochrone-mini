create table if not exists projects (
	id serial primary key,
	name varchar(140) not null,
	slug varchar(140) unique not null,
	description varchar(1000) not null,
	creation_date timestamp not null,
	target_date timestamp not null
);

create table if not exists users (
	id serial primary key,
	project_id int not null references projects(id),
	is_admin boolean not null default(false),
	slug varchar(140) unique not null,
	name varchar(140) not null,
	creation_date timestamp not null
);
