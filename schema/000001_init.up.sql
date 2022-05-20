CREATE TABLE users (
    id serial not null unique,
    name varchar(255),
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE trips (
    id serial not null unique,
    date_start timestamp with time zone,
    date_end timestamp with time zone,
    pause_duration bigint,
    active_duration bigint,
    distance real,
    avg_speed real,
    max_speed real,
    user_id int references users(id) on delete cascade not null
);

CREATE TABLE trip_points (
	id serial not null unique,
    latitude double precision,
    longitude double precision,
    point_date timestamp with time zone,
    speed real,
    trip_id int references trips(id) on delete cascade not null
);