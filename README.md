# go-micro-learning

Learning microservices with golang. The course is https://www.udemy.com/course/working-with-microservices-in-go/


# GRPC
For generating grpc launch the following command in the /logger-service/logs directory 
`bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\logs.proto
`