easy:
	cd internal/pkg/responses && easyjson --all responses.go
	cd .. && cd .. && cd ..
	cd internal/app/forum/models && easyjson --all forum.go
	cd .. && cd .. && cd .. && cd ..
	cd internal/app/user/models && easyjson --all user.go
	cd .. && cd .. && cd .. && cd ..
	cd internal/app/thread/models && easyjson --all thread.go
	cd .. && cd .. && cd .. && cd ..
	cd internal/app/post/models && easyjson --all post.go