all:
	echo "this is the counter for Chinese characters in Latex files"

test:
	go test ./...

deploy:
	echo "should commit latest changes before deploying"
	git tag v0.2.0
	git push origin v0.2.0