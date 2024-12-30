FROM golang:1.23.2 as build

WORKDIR /scraper

COPY . .

RUN go build -tags lambda.norpc -o main ./cmd/scraper/main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /scraper/main ./main

ENTRYPOINT [ "./main" ]