version=latest

deploy:
	docker build -t mwaaas/awsSSh:$(latest) .
	docker push mwaaas/awsSSh:$(latest)