# REPLACE THIS WITH YOUR PROJECT_ID
PROJECT := feed-pruning-proxy

GCLOUD := gcloud --project $(PROJECT)

.PHONY: deploy
deploy:
	$(GCLOUD) app deploy


.PHONY: localup
localup:
	go build && echo 'open http://localhost:8080' && GOOGLE_CLOUD_PROJECT=$(PROJECT) ./feed-pruning-proxy

.PHONY: test
test:
	go test -cover -v ./...
