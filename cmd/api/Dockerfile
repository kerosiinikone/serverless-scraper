FROM golang:1.23.2 as build

WORKDIR /api

COPY . .

RUN go build -tags lambda.norpc -o main ./cmd/api/main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /api/main ./main

ENTRYPOINT [ "./main" ]
