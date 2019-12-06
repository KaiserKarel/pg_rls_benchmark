run:
	dropdb rls_bench
	createdb rls_bench
	go run *.go