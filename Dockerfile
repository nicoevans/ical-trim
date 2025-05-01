FROM golang:1.18 as build
WORKDIR /ical-trim
COPY go.mod go.sum ./
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -tags lambda.norpc cmd/ical-trim-lambda/ical-trim-lambda.go
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /ical-trim/ical-trim-lambda ./ical-trim-lambda
COPY res/ res/
ENTRYPOINT [ "./ical-trim-lambda" ]
