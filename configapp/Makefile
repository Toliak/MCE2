test_e2e:
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it archlinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it alpine
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it debian
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it fedora
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it oraclelinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build  --rm -it rocky

test_pkg:
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce archlinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce alpine
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce debian
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce fedora
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce oraclelinux
	BUILDKIT_PROGRESS=plain docker compose -f e2e/docker-compose.yaml run --build --rm -it --entrypoint /test/mce rocky
