CREATE TABLE users (
    id                   uuid                     default gen_random_uuid() not null primary key,
    created_at           timestamp with time zone default now()             not null,
    updated_at           timestamp with time zone default now()             not null,
    phone_number         text                                               not null,

    jid                  varchar,
    telegram_id          text,
    context              text,
    -- conversation_buffer  text,
    -- conversation_summary text,
    user_name            text
    -- tools json
);


CREATE TABLE messages (
    id         uuid default gen_random_uuid() not null primary key,

    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,

    user_id    uuid references users(id),
    role       text,
    content    text,

    parent_id  uuid NULL default NULL references messages(id) ON DELETE SET NULL,
    agent_id   uuid NULL default NULL references agents(id) ON DELETE SET NULL
);


CREATE TABLE agents (
    id         uuid default gen_random_uuid() not null primary key,

    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,

    name         text not null,
    constitution text not null
);

CREATE UNIQUE INDEX agents_name_idx ON agents (name);
