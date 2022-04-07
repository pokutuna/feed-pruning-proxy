# REPLACE THIS WITH YOUR PROJECT_ID
PROJECT := slack-feed-proxy

GCLOUD := gcloud --project $(PROJECT)

.PHONY: deploy
deploy:
	$(GCLOUD) app deploy


.PHONY: localup
localup:
	go build && echo 'open http://localhost:8080' && ./slack-feed-proxy

.PHONY: test
test:
	go test -cover -v ./...
