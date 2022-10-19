pipeline {
    agent any

    tools {
        go '1.17.8'
    }

    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
    }
}
