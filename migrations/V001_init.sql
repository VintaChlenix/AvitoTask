CREATE TABLE segments(slug text PRIMARY KEY);

CREATE TABLE users_segments(
    user_id  integer,
    slug text,
    PRIMARY KEY(user_id, slug)
);
