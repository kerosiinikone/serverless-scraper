api:
	docker build -t api cmd/api
	docker tag api:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):api
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):api

analysis:
	docker build -t analysis cmd/analysis
	docker tag analysis:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):analysis
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):analysis

scraper:
	docker build -t scraper cmd/scraper
	docker tag scraper:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):scraper
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):scraper