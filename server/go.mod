module example.com/server

go 1.21.6

replace example.com/deck => ../deck

require example.com/deck v0.0.0-00010101000000-000000000000

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
)
