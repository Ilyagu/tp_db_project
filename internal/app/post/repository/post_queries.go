package repository

const (
	GetPostQuery = `select id, case when parent is null then 0 else parent end as parent,
	author, message, is_edited, forum, thread, created, path
	from posts where id=$1;`

	UpdatePostQuery = `update posts
	set message = (case
			when ltrim($2) = '' then message
			else $2 end),
	is_edited = (case
			when ltrim($2) = '' or $2=message then false
			else true end)
	where id=$1
	returning id, case when parent is null then 0 else parent end as parent,
	author, message, is_edited, forum, thread, created, path;`

	GetStatusQuery = `select (select count(*) from users) as "user",
	(select count(*) from forums) as forum,
	(select count(*) from threads) as thread,
	(select count(*) from posts) as post;`

	ClearAllQuery = `truncate table users cascade;
	truncate table forums cascade;
	truncate table threads cascade;
	truncate table posts cascade;
	truncate table votes cascade;`
)
