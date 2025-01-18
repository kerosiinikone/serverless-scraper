all: api analysis scraper

api:
	docker buildx build --platform linux/amd64 --provenance false -t api -f cmd/api/Dockerfile .
	docker tag api:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):api-req-lambda
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):api-req-lambda

analysis:
	docker buildx build --platform linux/amd64 --provenance false -t analysis -f cmd/analysis/Dockerfile .
	docker tag analysis:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):analysis-lambda
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):analysis-lambda

scraper:
	docker buildx build --platform linux/amd64 --provenance false -t scraper -f cmd/scraper/Dockerfile .
	docker tag scraper:latest $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):scraper-lambda
	docker push $(AWS_ACCOUNT).dkr.ecr.eu-north-1.amazonaws.com/$(PROJECT_NAME):scraper-lambda