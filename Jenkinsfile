pipeline {
    agent { label 'dapper '}
    triggers { pollSCM('0 * * * *') }
    options { timestamps() }

    stages {
        // Build first to catch compilation errors - there's no point running
        // unit tests against code that doesn't compile.
        stage("Build") {
            steps {
                sh 'make -B build-docker'
            }

            // Always publish c2 and other binaries to Jenkins' filestore.
            post { success { archiveArtifacts artifacts: 'bin/*', fingerprint: true } }
        }

        stage("Quality checks") {
            parallel {
                stage("Unit test only") {
                    steps {
						sh 'make -B test-junit-docker'
					}
					post {
						success { junit 'test_junit.xml' }
					}
                }

                stage("Lints") {
                    steps {
                        sh 'make lint-docker'
                    }
                }
            }
        }
    }

    post {
        always {
            sh 'make clean-docker'
            sh 'sudo chown -R $USER *'
        }
    }
}
