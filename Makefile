version=latest

deploy:
	docker build -t mwaaas/aws_ssh:$(version) .
	docker push mwaaas/aws_ssh:$(version)
