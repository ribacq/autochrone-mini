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

create table if not exists measures (
	id serial primary key,
	project_id int not null references projects(id),
	code varchar(42) not null,
	name varchar(42) not null,
	unit varchar(42) not null,
	unique (project_id, code)
);

create table if not exists measures_goals (
	measure_id int primary key references measures(id),
	goal_low int not null,
	goal_high int not null
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
