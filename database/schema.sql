create table conversations (
    id                   uuid                     default gen_random_uuid() not null
                                                  primary key,
    created_at           timestamp with time zone default now()             not null,
    updated_at           timestamp with time zone default now()             not null,
    phone_number         text                                               not null,
    jid                  varchar,
    context              text,
    conversation_buffer  text,
    conversation_summary text,
    user_name            text
);