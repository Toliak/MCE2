build:
	go build -o bin/mce ./cmd/mce
	go build -o bin/mce-e2e ./cmd/mce-e2e

test_e2e:
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it archlinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it alpine
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it debian
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it fedora
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it oraclelinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it rocky

test_pkg:
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e archlinux		repo-install -to-install python3:zsh
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e alpine		repo-install -to-install python3:zsh
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e debian		repo-install -to-install python3:zsh
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e fedora		repo-install -to-install python3:zsh
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e oraclelinux	repo-install -to-install python3:zsh
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/bin/mce-e2e rocky			repo-install -to-install python3:zsh
