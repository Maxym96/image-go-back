 CREATE TABLE IF NOT EXISTS images
(
    id           serial PRIMARY KEY,
    obj_id       int NULL,
    name         varchar(250) NOT NULL,
    quality      smallint NOT NULL,
    created_date timestamp    NOT NULL,
    updated_date timestamp    NOT NULL,
    deleted_date timestamp    NULL
);
