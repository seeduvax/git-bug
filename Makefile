GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_LAST_TAG:=$(shell git describe --abbrev=0 --tags)
GIT_EXACT_TAG:=$(shell git name-rev --name-only --tags HEAD)
COMMANDS_PATH:=github.com/MichaelMure/git-bug/commands
LDFLAGS:=-X ${COMMANDS_PATH}.GitCommit=${GIT_COMMIT} \
	-X ${COMMANDS_PATH}.GitLastTag=${GIT_LAST_TAG} \
	-X ${COMMANDS_PATH}.GitExactTag=${GIT_EXACT_TAG}

git-bug: $(shell find . -name "*.go")
	go build -ldflags "$(LDFLAGS)" .

doc: 
	go generate

# produce a build debugger friendly
debug-build:
	go build -ldflags "$(LDFLAGS)" -gcflags=all="-N -l" .

install: git-bug doc
	go install -ldflags "$(LDFLAGS)" .

test: git-bug
	go test -v -bench=. ./...

pack-webui: git-bug
	npm run --prefix webui build
	go run webui/pack_webui.go

# produce a build that will fetch the web UI from the filesystem instead of from the binary
debug-webui:
	go build -ldflags "$(LDFLAGS)" -tags=debugwebui

clean-local-bugs:
	git for-each-ref refs/bugs/ | cut -f 2 | xargs -r -n 1 git update-ref -d
	git for-each-ref refs/remotes/origin/bugs/ | cut -f 2 | xargs -r -n 1 git update-ref -d
	rm -f .git/git-bug/bug-cache

clean-remote-bugs:
	git ls-remote origin "refs/bugs/*" | cut -f 2 | xargs -r git push origin -d

clean-local-identities:
	git for-each-ref refs/identities/ | cut -f 2 | xargs -r -n 1 git update-ref -d
	git for-each-ref refs/remotes/origin/identities/ | cut -f 2 | xargs -r -n 1 git update-ref -d
	rm -f .git/git-bug/identity-cache

.PHONY: build install test pack-webui debug-webui clean-local-bugs clean-remote-bugs doc
