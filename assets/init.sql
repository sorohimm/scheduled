CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS chats (
    chat_id bigint,
    group_name text,
    group_uuid UUID,
    notif_time text
);