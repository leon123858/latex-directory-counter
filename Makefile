all:
	echo "this is the counter for Chinese characters in Latex files"

test:
	go test ./...

deploy:
	git tag v0.1.0
	git push origin v0.1.0