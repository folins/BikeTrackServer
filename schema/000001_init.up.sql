CREATE TABLE users (
    id serial not null unique,
    name varchar(255),
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE temp_users (
    id serial not null unique,
    email varchar(255) not null unique,
    confirm_code int not null
);

CREATE TABLE trips (
    id serial not null unique,
    date_start timestamp,
    pause_duration bigint,
    duration bigint,
    distance real,
    avg_speed real,
    max_speed real,
    user_id int references users(id) on delete cascade not null
);