create extension if not exists citext;

create unlogged table if not exists users (
  id       bigserial         unique not null,
  nickname citext         collate "ucs_basic" not null primary key,
  fullname varchar(255)         not null,
  about    text           not null,
  email    citext         unique not null
);

create unlogged table if not exists forums (
  id      bigserial,
  title   text       not null,
  "user"  citext       not null,
  slug    citext primary key,
  posts   int    not null default 0,
  threads int        not null default 0,
  constraint fk_fr_user foreign key ("user") references users (nickname)
);

create unlogged table if not exists threads (
  id        bigserial primary key,
  title     text        not null,
  author    citext        not null,
  forum     citext        not null,
  message   text        not null,
  slug      citext,
  votes     int         not null default 0,
  created   timestamp with time zone default now(),
  constraint fk_th_user foreign key (author) references users (nickname),
  constraint fk_th_forum foreign key (forum) references forums (slug)
);

create unlogged table posts
(
    id        bigserial primary key,
    parent    integer default 0,
    author    citext not null,
    message   text   not null,
    is_edited boolean default false,
    forum     citext,
    thread    integer,
    created   timestamp with time zone default now(),
    path      bigint[] default ARRAY []::integer[],
    constraint fk_ps_user foreign key (author) references users (nickname),
    constraint fk_ps_thread foreign key (thread) references threads (id),
    constraint fk_ps_forum foreign key (forum) references forums (slug),
    constraint fk_ps_post foreign key (parent) references posts (id)
);

create unlogged table if not exists votes
(
    nickname  citext not null,
    thread_id int    not null,
    voice     int    not null,
    constraint fk_vt_user foreign key (nickname) references users (nickname),
    constraint fk_vt_thread foreign key (thread_id) references threads (id),
    unique (nickname, thread_id)
);

create unlogged table if not exists users_to_forums
(
    nickname citext collate "ucs_basic" not null,
    forum    citext not null,
    constraint fk_utf_user foreign key (nickname) references users (nickname),
    constraint fk_utf_forum foreign key (forum) references forums (slug),
    unique (nickname, forum)
);

-- trigger for increment threads count on forum
create or replace function add_thread()
    returns trigger as
$add_thread$
begin
    update forums
    set threads = forums.threads + 1
    where slug = new.forum;
    return new;
end;
$add_thread$ language plpgsql;

drop trigger if exists add_thread on threads;
create trigger add_thread
    after insert
    on threads
    for each row
execute procedure add_thread();

-- trigger for insert vote to forum
create or replace function insert_vote()
    returns trigger as
$insert_vote$
begin
    update threads
    set votes = votes + new.voice
    where id = new.thread_id;
    return new;
end;
$insert_vote$ language plpgsql;

drop trigger if exists insert_vote on votes;
create trigger insert_vote
    after insert
    on votes
    for each row
execute procedure insert_vote();

-- trigger for change vote to forum
create or replace function change_vote()
    returns trigger as
$change_vote$
begin
    update threads
    set votes=votes + 2 * new.voice
    where id = new.thread_id;
    return new;
end;
$change_vote$ language plpgsql;


drop trigger if exists change_vote on votes;
create trigger change_vote
    after update
    on votes
    for each row
execute procedure change_vote();

-- trigger for add users on forum on posts and threads
create or replace function add_users_to_forum()
    returns trigger as
$add_users_to_forum$
begin
    insert into users_to_forums (nickname, forum)
    values (new.author, new.forum)
    on conflict do nothing;
    return new;
end;
$add_users_to_forum$ language plpgsql;

drop trigger if exists add_users_to_forum_on_threads on threads;
create trigger add_users_to_forum_on_threads
    after insert
    on threads
    for each row
execute procedure add_users_to_forum();

drop trigger if exists add_users_to_forum_on_posts on posts;
create trigger add_users_to_forum_on_posts
    after insert
    on posts
    for each row
execute procedure add_users_to_forum();

-- trigger for path
create or replace function update_path() returns trigger as
$update_path$
declare
    parent_path         bigint[];
    first_parent_thread int;
begin
    if (new.parent is null) then
        new.path := array_append(new.path, new.id);
    else
        select path from posts where id = new.parent into parent_path;
        select thread from posts where id = parent_path[1] into first_parent_thread;
        if not found or first_parent_thread != new.thread then
            raise exception 'parent is from different thread' using errcode = '00409';
        end if;

        new.path := new.path || parent_path || new.id;
    end if;
    update forums set posts=posts + 1 where forums.slug = new.forum;
    return new;
end
$update_path$ language plpgsql;

drop trigger if exists update_path on posts;
create trigger update_path
    before insert
    on posts
    for each row
execute procedure update_path();

create index post_first_parent_thread_index on posts ((posts.path[1]), thread);
create index post_first_parent_id_index on posts ((posts.path[1]), id);
create index post_first_parent_index on posts ((posts.path[1]));
create index post_path_index on posts ((posts.path));
create index post_thread_index on posts (thread);
create index post_thread_id_index on posts (thread, id);
create index post_path_id_index on posts (id, (posts.path));
create index post_thread_path_id_index on posts (thread, (posts.parent), id);

create index forum_slug_lower_index on forums (slug);

create index users_nickname_lower_index on users (nickname);
create index users_nickname_index on users (nickname);
create index users_email_index on users (email);

create index users_forum_forum_user_index on users_to_forums (forum, nickname);
create index users_forum_user_index on users_to_forums (nickname);
create index users_forum_forum_index on users_to_forums (forum);

create index thread_slug_index on threads (slug);
create index thread_slug_id_index on threads (slug, id);
create index thread_forum_lower_index on threads (forum);
create index thread_id_forum_index on threads (id, forum);
create index thread_created_index on threads (created);

create index vote_nickname on votes (nickname, thread_id, voice);

vacuum analyse;