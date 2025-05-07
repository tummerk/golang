CREATE TABLE schedules(
    ID SERIAL PRIMARY KEY,
    medicament_name varchar(100) not null,
    user_id int not null,
    receptions_per_day int not null,
    date_start DATE,
    date_end DATE
);
