pipeline {
    agent any

    // tools {
    //     go "Go 1.14.6"
    // }

    // set jenkins environment
    environment {
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        // GOBIN = "${GOPATH}/bin"
        CGO_ENABLED = 0
        GO1136MODULE = "off"
    }

    stages {
        /*
        stage("Pre Test") {
            steps {
                bat "echo 'Installing dependencies'"
                bat "go version"
                bat "go get -u golang.org/x/lint/golint"
            }
        }
        */
        
        stage("Build") {
            steps {
                bat "echo 'Compiling and building'"
                bat "go build"
            }
        }

        stage("Test") {
            steps {
                // withEnv(["PATH+GO=${GOPATH}/bin"]){
                    // bat "echo 'Vetting'"
                    // bat 'go vet .'

                    // bat "echo 'Linting'"
                    // bat 'golint ./...'

                    bat "echo 'Testing'"
                    bat 'go test -v ./test'
                // }
            }
        }
    }
    post {
        always {
            bat "echo 'Finished! Jenkins Build ${currentBuild.currentResult}: Job ${env.JOB_NAME}'"
        }
        failure {
            bat "echo 'should send a message saying the pipeline has FAILED status'"
        }
        unstable  {
            bat "echo 'should send a message saying the pipeline has UNSTABLE status'"
        }
    }
}