# Use an official Go runtime as a parent image
FROM golang
#FROM golang:alpine #(lighter image)

# Make port 80 available to the world outside this container
EXPOSE 8080

# setup GOPATH and friends
RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# Define environment variable
ENV NAME world

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/app/
#RUN mkdir -p $GOPATH/src/app/vendor
ADD ./ $GOPATH/src/app/
#ADD ./vendor/ $GOPATH/src/app/vendor/
ADD ./vendor $GOPATH/src

# should be able to build now
WORKDIR $GOPATH/src/app 
RUN go build -o main .
CMD ["/go/src/app/main"]

#-------------------------------------------------------USAGE-------------------------------------------------------#
#
# to build: $ sudo docker build ./ -t mycontainer
#
# to run full node: $ sudo docker run -p 8084:8080 --ip=10.0.0.24 --net=mynet 
#		 	-e K=3 -e VIEW="10.0.0.21:8080,10.0.0.22:8080,10.0.0.23:8080,10.0.0.24:8080" 
#		 	-e IPPORT="10.0.0.24:8080" mycontainer
#
# to run half node: $ sudo docker run -p 8084:8080 --ip=10.0.0.24 --net=mynet -e IPPORT="10.0.0.24:8080" mycontainer