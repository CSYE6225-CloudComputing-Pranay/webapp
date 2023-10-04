
# webapp

A project for the course "Network structures and cloud computing" that is written in Go language and build to be deployable in cloud platforms


## Prerequisites to build and deploy the application locally

- Use the following link to install go on your machine, version: 1.21.x
    - https://go.dev/doc/install
- Install mysql database, version: 5.7
- Checkout the repository using the following link:
    - https://github.com/CSYE6225-CloudComputing-Pranay/webapp.git
- Create a .env or .env.local file in the project directory and Configure the following parameters:
    - DB_HOST: database host
    - DB_USER: database user
    - DB_PASSWORD: database password
    - DB_NAME: database name
    - DB_PORT: database port
    - PORT: Port to run the application
- Install required dependencies using the following command:
    - go get {DEPENDECY_NAME}
- Following are the required dependencies :
    - github.com/gin-gonic/gin v1.9.1
    - github.com/google/uuid v1.3.1
    - github.com/joho/godotenv v1.5.1
    - github.com/stretchr/testify v1.8.3
    - golang.org/x/crypto v0.13.0
    - gorm.io/driver/mysql v1.5.1
    - gorm.io/gorm v1.25.4
- Run "go mod tidy" command from project directory to make sure that the dependencies are properly installed
- Run "go build webapp/cmd/main" command from project directory to build the project
- Run "cd ./test ; go test -run TestHealthTestSuite ; cd .." command from project directory to run the integration test
- Now, run "go run webapp/cmd/main" command from project directory to deploy the project

