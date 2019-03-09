package migrations

//go:generate go get -v github.com/shuLhan/go-bindata/go-bindata
//go:generate go-bindata -pkg migrations -o migrations_gen.go sqlite3/
