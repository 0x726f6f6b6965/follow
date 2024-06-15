CREATE TABLE IF NOT EXISTS t_user (
    id SERIAL NOT NULL,
    username character varying(128) UNIQUE NOT NULL,
    salt character varying(64) NOT NULL,
    password character varying(128) NOT NULL,
    create_time timestamp NOT NULL DEFAULT (now())::timestamp,
    update_time timestamp NOT NULL DEFAULT (now())::timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS t_followers (
    id SERIAL NOT NULL,
    follower_id integer NOT NULL,
    following_id integer NOT NULL,
    create_time timestamp NOT NULL DEFAULT (now())::timestamp,
    update_time timestamp NOT NULL DEFAULT (now())::timestamp,
    FOREIGN KEY (follower_id) REFERENCES t_user (id),
    FOREIGN KEY (following_id) REFERENCES t_user (id),
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_follow_relation ON t_followers (follower_id, following_id);

CREATE INDEX IF NOT EXISTS idx_followers ON t_followers (follower_id);

CREATE INDEX IF NOT EXISTS idx_following ON t_followers (following_id);