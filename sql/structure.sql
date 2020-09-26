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

create type goal_direction_type as enum ('min', 'max', 'none');

create table if not exists measures (
	id serial primary key,
	project_id int not null references projects(id),
	code varchar(42) not null,
	name varchar(140) not null,
	unit varchar(42) not null,
	goal_direction goal_direction_type not null,
	goal int not null,
	unique (project_id, code)
);

create table if not exists notes (
	id serial primary key,
	project_id int not null references projects(id),
	user_id int not null references users(id),
	creation_date timestamp not null,
	comment varchar(1000) not null
);

create table if not exists notes_measures_values (
	id serial primary key,
	note_id int references notes(id),
	measure_id int references measures(id),
	value int not null,
	unique (note_id, measure_id)
);
