# Telling to use Docker's golang ready image
FROM golang

RUN mkdir app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD [ "/app/main" ]

# docker tag local-image:tagname new-repo:tagname
# docker push new-repo:tagname