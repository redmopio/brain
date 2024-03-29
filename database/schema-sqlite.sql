CREATE TABLE user (
    id         varchar not null primary key,
    -- created_at           timestamp with time zone default now()             not null,
    -- updated_at           timestamp with time zone default now()             not null,
    phone_number         text                                               not null,

    jid                  varchar,
    telegram_id          text,
    context              text,
    -- conversation_buffer  text,
    -- conversation_summary text,
    user_name            text
    -- tools json
);


CREATE TABLE message (
    id         varchar not null primary key,

    -- created_at timestamp with time zone default now() not null,
    -- updated_at timestamp with time zone default now() not null,

    user_id    uuid references user(id),
    role       text,
    content    text,

    parent_id  uuid references message(id)
);