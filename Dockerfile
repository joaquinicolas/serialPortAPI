FROM resin/raspberry-pi2-golang
ENV WORKSPACE $GOPATH/src/github.com/joaquinicolas/Elca
#RUN mkdir -p $WORKSPACE
WORKDIR $WORKSPACE
ADD ./ $WORKSPACE
RUN go build
CMD ["./Elca"]

