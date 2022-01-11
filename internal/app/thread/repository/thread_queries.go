package repository

const (
	CreateThreadQuery = `insert into threads(title, author, forum, message, slug, created)
	values ($1,(select nickname from users where nickname=$2),(select slug from forums where forums.slug=$3),$4,$5,$6)
	returning id, title, author, forum, message, slug, votes, created;`

	CreateThreadNoCreatedQuery = `insert into threads(title, author, forum, message, slug)
	values ($1,(select nickname from users where nickname=$2),(select slug from forums where forums.slug=$3),$4,$5)
	returning id, title, author, forum, message, slug, votes, created;`

	GetThreadBySlugOrIdQuery = `select id, title, (select nickname from users where nickname=threads.author) as author,
		(select forums.slug from forums where slug=threads.forum) as forum, message, slug, votes, created
	from threads where id=$1 or slug=$2;`

	GetThreadByIdQuery = `select id, title, (select nickname from users where nickname=threads.author) as author,
		(select forums.slug from forums where slug=threads.forum) as forum, message, slug, votes, created
	from threads where id=$1;`

	UpdateThreadQuery = `update threads set title=$1, message=$2 where slug=$3 or id=$4
	returning id, title, author, forum, message, slug, votes, created;`

	CreateVoteQuery = "insert into votes (nickname, thread_id, voice) values ($1, $2, $3)"

	UpdateVoteQuery = "update votes set voice = $1 where thread_id = $2 and nickname = $3 and voice != $1"
)
