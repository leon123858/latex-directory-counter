all:
	echo "this is the counter for Chinese characters in Latex files"

test:
	go test ./...

deploy:
	git ls-remote --tags origin
	git tag v0.2.1
	git push origin v0.2.1
	git tag lastest -f
	git push origin lastest -f