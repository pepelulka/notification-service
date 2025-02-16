BEGIN;

CREATE TABLE IF NOT EXISTS persons
(
    person_id SERIAL PRIMARY KEY,
    email character varying(128),
    telegram_id character varying(128),
    phone_number character varying(12)
);

CREATE INDEX ON persons (email);
CREATE INDEX ON persons (telegram_id);
CREATE INDEX ON persons (phone_number);

CREATE TABLE IF NOT EXISTS groups
(
    group_id SERIAL PRIMARY KEY,
    name character varying(128)
);

CREATE INDEX ON groups(name);

CREATE TABLE IF NOT EXISTS person_to_group
(
    person_id INTEGER,
    group_id INTEGER,
    CONSTRAINT fk_person_id FOREIGN KEY (person_id) REFERENCES persons (person_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES groups (group_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (person_id, group_id)
);

END;
