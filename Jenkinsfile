pipeline{
    agent any
    environment {
            PROJECT_PATH = "/go/src/github.com/cjburchell/loki-go"
            VERSION = "v1.1.${env.BUILD_NUMBER}"
            repository = "github.com/cjburchell/loki-go.git"
    }

    parameters {
            booleanParam(name: 'Release', defaultValue: false, description: 'Should tag release?')
        }

    stages {
        stage('Setup') {
            steps {
                script{
                    slackSend color: "good", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} started"
                }
                
                /* Let's make sure we have the repository cloned to our workspace */
                checkout scm
             }
         }
        stage('Static Analysis') {
            parallel {
                stage('Vet') {
                    agent {
                        docker {
                            image 'cjburchell/goci:1.14'
                            args '-v $WORKSPACE:$PROJECT_PATH'
                        }
                    }
                    steps {
                        script{
                                sh """go vet ./..."""

                                def checkVet = scanForIssues tool: [$class: 'GoVet']
                                publishIssues issues:[checkVet]
                        }
                    }
                }

                stage('Lint') {
                    agent {
                        docker {
                            image 'cjburchell/goci:1.14'
                            args '-v $WORKSPACE:$PROJECT_PATH'
                        }
                    }
                    steps {
                        script{
                            sh """golint ./..."""

                            def checkLint = scanForIssues tool: [$class: 'GoLint']
                            publishIssues issues:[checkLint]
                        }
                    }
                }
            }
        }
        stage('Release') {
            when { expression { params.Release } }
            steps {
                 script {
                    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'github', usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD']]) {
                        sh """git tag ${VERSION}"""
                        sh """git push https://${env.GIT_USERNAME}:${env.GIT_PASSWORD}@${env.repository} ${VERSION}"""
                    }

                }
            }
        }
    }

    post {
        always {
              script{
                  if ( currentBuild.currentResult == "SUCCESS" ) {
                    slackSend color: "good", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was successful"
                  }
                  else if( currentBuild.currentResult == "FAILURE" ) {
                    slackSend color: "danger", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was failed"
                  }
                  else if( currentBuild.currentResult == "UNSTABLE" ) {
                    slackSend color: "warning", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was unstable"
                  }
                  else {
                    slackSend color: "danger", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} its result (${currentBuild.currentResult}) was unclear"
                  }
              }
        }
    }
}