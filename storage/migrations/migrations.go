package migrations

//go:generate go get github.com/go-bindata/go-bindata
//go:generate go get -u github.com/go-bindata/go-bindata/...
//go:generate go-bindata -pkg migrations -o migrations_gen.go sqlite3/
