FROM golang:1.23.2 as build

WORKDIR /analysis

COPY . .

RUN go build -tags lambda.norpc -o main ./cmd/analysis/main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /analysis/main ./main

ENTRYPOINT [ "./main" ]
