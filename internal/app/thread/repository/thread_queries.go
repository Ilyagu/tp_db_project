package repository

const (
	CreateThreadQuery = `insert into threads(title, author, forum, message, slug, created)
	values ($1,$2,$3,$4,$5,$6)
	returning id, title, author, forum, message, slug, votes, created;`

	CreateThreadNoCreatedQuery = `insert into threads(title, author, forum, message, slug)
	values ($1,$2,$3,$4,$5)
	returning id, title, author, forum, message, slug, votes, created;`

	GetThreadBySlugOrIdQuery = `select t.id, t.title, u.nickname as author, f.slug as forum, t.message, t.slug, t.votes, t.created
    from threads t join forums f on t.forum = f.slug join users u on t.author = u.nickname
	where t.slug=$1 or t.id=$2;`

	UpdateThreadQuery = `update threads set title=$1, message=$2 where slug=$3 or id=$4
	returning id, title, author, forum, message, slug, votes, created;`

	CreateVoteQuery = "insert into votes (nickname, thread_id, voice) values ($1, $2, $3)"

	UpdateVoteQuery = "update votes set voice = $1 where thread_id = $2 and nickname = $3 and voice != $1"
)
