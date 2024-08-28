#!/usr/bin/env groovy
def bob() {
    def defaultBobImage = 'armdocker.rnd.ericsson.se/proj-adp-cicd-drop/bob.2.0:1.7.0-55'

    return new BobCommand()
            .bobImage(defaultBobImage)
            .envVars([
                    GERRIT_REVIEW_PASS_USER   : '${GERRIT_REVIEW_PASS_USR}',
                    GERRIT_REVIEW_PASS_PASS   : java.net.URLEncoder.encode("${env.GERRIT_REVIEW_PASS_PSW}", "UTF-8"),
                    GERRIT_USERNAME           : '${GERRIT_USERNAME}',
                    GERRIT_PASSWORD           : '${GERRIT_PASSWORD}',
                    HOME                      : '${HOME}',
                    ISO_VERSION               : '${ISO_VERSION}',
                    RELEASE                   : '${RELEASE}',
                    SONAR_HOST_URL            : '${SONAR_HOST_URL}',
                    SONAR_AUTH_TOKEN          : '${SONAR_AUTH_TOKEN}',
                    GERRIT_CHANGE_NUMBER      : '${GERRIT_CHANGE_NUMBER}',
                    KUBECONFIG                : '${KUBECONFIG}',
                    USER                      : '${USER}',
                    SELI_ARTIFACTORY_REPO_USER: '${CREDENTIALS_SELI_ARTIFACTORY_USR}',
                    SELI_ARTIFACTORY_REPO_PASS: '${CREDENTIALS_SELI_ARTIFACTORY_PSW}',
                    MAVEN_CLI_OPTS            : '${MAVEN_CLI_OPTS}',
                    OPEN_API_SPEC_DIRECTORY   : '${OPEN_API_SPEC_DIRECTORY}'
            ])
            .needDockerSocket(true)
            .toString()
}

pipeline {
    agent {
        node {
            label 'GridEngine'
        }
    }

    environment {
        FOSSA_CLI = true
        GERRIT_REVIEW_PASS = credentials('GERRIT_PASSWORD')
        CREDENTIALS_SELI_ARTIFACTORY = credentials('SELI_ARTIFACTORY')

        HADOLINT_ENABLED = "true"
        KUBEAUDIT_ENABLED = "true"
        KUBESEC_ENABLED = "true"
    }

    stages {
        stage('Clean') {
            steps {
                echo 'Inject settings.xml into workspace:'
                configFileProvider([configFile(fileId: "${env.SETTINGS_CONFIG_FILE_NAME}", targetLocation: "${env.WORKSPACE}")]) {}
                archiveArtifacts allowEmptyArchive: true, artifacts: 'ruleset2.0.yaml, precodereview.Jenkinsfile'
                sh "${bob()} clean"
            }
        }

        stage('Generate Variables') {
            steps {
                sh "${bob()} generate-variables"
            }
        }

        stage('Download Golang dependencies') {
            steps {
                sh "${bob()} download-go-dependencies"
            }
        }

        stage('Build Golang') {
            steps {
                sh "${bob()} build-go"
            }
        }

        stage('Test') {
            steps {
                sh "${bob()} test"
            }
        }

        stage('SonarQube analysis') {
            steps {
                withSonarQubeEnv("${env.SQ_SERVER}") {
                    sh "${bob()} sonar-enterprise-pcr"
                }
                timeout(time: 10, unit: 'MINUTES') {
                    // Parameter indicates whether to set pipeline to UNSTABLE if Quality Gate fails
                    // true = set pipeline to UNSTABLE, false = don't
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Lint') {
            steps {
                parallel(
                    "lint markdown": {
                        sh "${bob()} lint:markdownlint lint:vale"
                    },
                    "lint helm": {
                        sh "${bob()} lint:helm"
                    },
                    "lint helm design rule checker": {
                        sh "${bob()} lint:helm-chart-check"
                    },
                    "lint golang": {
                        sh "${bob()} lint:golang"
                    }
                )
            }
            post {
                always {
                    archiveArtifacts allowEmptyArchive: true, artifacts: 'zally-api-lint-report.txt, .bob/design-rule-check-report.*'
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "${bob()} create-image-build-name-internal"
                sh "${bob()} build-docker-image"

            }
        }

        stage('Verify Docker Image') {
            steps {
                sh "${bob()} test-docker-image:run-sample-app-container"
                script {
                    def response = sh(script: 'curl -s -w " Response code: %{response_code}" http://localhost:8050/hello', returnStdout: true)
                    echo "Response: " + response
                    if (response.contains("Hello World!!") && response.contains("Response code: 200")) {
                        echo "Hello Endpoint works as intended"
                    } else {
                        error('Build failed due to Hello Endpoint error')
                    }
                    response = sh(script: 'curl -s -w " Response code: %{response_code}" http://localhost:8050/health', returnStdout: true)
                    echo "Response: " + response
                    if (response.contains("Ok") && response.contains("Response code: 200")) {
                        echo "Health Endpoint works as intended"
                    } else {
                        error('Build failed due to Health Endpoint error')
                    }
                    response = sh(script: 'curl -s -w " Response code: %{response_code}" http://localhost:8050/metrics', returnStdout: true)
                    echo "Response: " + response
                    if (response.contains("hello_world_requests_total 1") && response.contains("Response code: 200")) {
                        echo "Metrics Endpoint works as intended"
                    } else {
                        error('Build failed due to Metrics Endpoint error')
                    }
                    sh(script: 'curl http://localhost:8050/hello')
                    response = sh(script: 'curl -s -w " Response code: %{response_code}" http://localhost:8050/metrics', returnStdout: true)
                    echo "Response: " + response
                    if (response.contains("hello_world_requests_total 2")) {
                        echo "The 'hello_world_requests_total' metric iterated as intended"
                    } else {
                        error('Build failed due to Metrics Endpoint error')
                    }
                }
            }
        }

        // bob rule 'create-image-build-name-internal' needs to run before the 'Contracts Tests' stage to
        // populate the 'image-build-name' variable
        stage('Contract Tests') {
            steps {
                script {
                    try {
                        sh "${bob()} contract-testing"
                    }
                    finally {
                        sh "${bob()} copy-contract-testing-output"
                    }
                }
            }
        }
        stage('Vulnerability Analysis') {
            steps {
                parallel(
                    "Hadolint": {
                        script {
                            if (env.HADOLINT_ENABLED == "true") {
                                sh "${bob()} hadolint"
                                archiveArtifacts "**/va-reports/hadolint-scan/**.*"
                            }
                        }
                    },
                    "Kubeaudit": {
                        script {
                            if (env.KUBEAUDIT_ENABLED == "true") {
                                sh "${bob()} kubeaudit"
                                archiveArtifacts "**/va-reports/kube-audit-report/**/*"
                            }
                        }
                    },
                    "Kubsec": {
                        script {
                            if (env.KUBESEC_ENABLED == "true") {
                                sh "${bob()} kubesec"
                                archiveArtifacts "**/va-reports/kubesec-reports/*"
                            }
                        }
                    }
                )
            }
            post {
                always {
                    sh "${bob()} generate-VA-report:no-upload"
                    archiveArtifacts allowEmptyArchive: true, artifacts: '**/va-reports/Vulnerability_Report_2.0.md'
                }
            }
        }

        stage('Create Helm Package') {
            steps {
                sh "${bob()} create-helm-build-vars-internal"
                sh "${bob()} publish-helm:package-helm"
            }
        }

        stage('Publish Helm Package to Internal Repo') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'SELI_ARTIFACTORY', usernameVariable: 'SELI_ARTIFACTORY_REPO_USER', passwordVariable: 'SELI_ARTIFACTORY_REPO_PASS')]) {
                    sh "${bob()} publish-helm:upload-helm-to-repo"
                }
            }
        }

    }
    post {
        always {
            archiveArtifacts artifacts: '**/*.log,**/*.ear,**/*.rpm,**/design-rule-check-report.html,**/spring-cloud-contract-output/reports/tests/**', allowEmptyArchive: true
            sh "${bob()} archive-artifacts"
            archiveArtifacts 'artifact.properties'
            sh "${bob()} cleanup-images"
        }
        failure {
            mail to: "${env.GERRIT_PATCHSET_UPLOADER_EMAIL}",
                    from: 'lord.vader@ericsson.com',
                    subject: "Failed Pipeline: ${currentBuild.fullDisplayName}",
                    body: "Failure on ${env.BUILD_URL}"
        }
        success {
            cleanWs()
        }
    }
}

// More about @Builder: http://mrhaki.blogspot.com/2014/05/groovy-goodness-use-builder-ast.html
import groovy.transform.builder.Builder
import groovy.transform.builder.SimpleStrategy

@Builder(builderStrategy = SimpleStrategy, prefix = '')
class BobCommand {
    def bobImage = 'bob.2.0:latest'
    def envVars = [:]
    def needDockerSocket = false
    String toString() {
        def env = envVars
                .collect({ entry -> "-e ${entry.key}=\"${entry.value}\"" })
                .join(' ')
        def cmd = """\
            |docker run
            |--init
            |--rm
            |--workdir \${PWD}
            |--user \$(id -u):\$(id -g)
            |-v \${PWD}:\${PWD}
            |-v /etc/group:/etc/group:ro
            |-v /etc/passwd:/etc/passwd:ro
            |-v \${HOME}:\${HOME}
            |-v /proj/mvn/:/proj/mvn
            |${needDockerSocket ? '-v /var/run/docker.sock:/var/run/docker.sock' : ''}
            |${env}
            |\$(for group in \$(id -G); do printf ' --group-add %s' "\$group"; done)
            |--group-add \$(stat -c '%g' /var/run/docker.sock)
            |${bobImage}
            |"""

        return cmd
                .stripMargin()           // remove indentation
                .replace('\n', ' ')      // join lines
                .replaceAll(/[ ]+/, ' ') // replace multiple spaces by one
    }
}