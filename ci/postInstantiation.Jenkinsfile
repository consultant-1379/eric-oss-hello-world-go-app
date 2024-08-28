#!/usr/bin/env groovy

pipeline {
    agent {
        label env.SLAVE_LABEL
    }

    parameters {
        string(name: 'GERRIT_REFSPEC',
                defaultValue: 'refs/heads/master',
                description: 'Referencing to a commit by Gerrit RefSpec')
        string(name: 'SLAVE_LABEL',
                defaultValue: 'evo_docker_engine_gic_IDUN',
                description: 'Specify the slave label that you want the job to run on')
        string(name: 'INGRESS_PREFIX',
                defaultValue: '',
                description: 'The prefix to the ingress URL')
        string(name: 'INGRESS_HOST',
                defaultValue: '',
                description: 'The EIC APIGW Host')
        string(name: 'INGRESS_LOGIN_USER',
                defaultValue: '',
                description: 'The user name to use for login')
        string(name: 'INGRESS_LOGIN_PASSWORD',
                defaultValue: '',
                description: 'The password to use')
        string(name: 'INGRESS_SAMPLE_APP_USER',
                defaultValue: '',
                description: 'The user name to use for eacc')
        string(name: 'INGRESS_GAS_USER',
                defaultValue: '',
                description: 'The user name to use for gas')

    }

    options {
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        buildDiscarder(logRotator(daysToKeepStr: '14', numToKeepStr: '40', artifactNumToKeepStr: '40', artifactDaysToKeepStr: '14'))
    }

     environment {
        INGRESS_SCHEMA = "${params.INGRESS_PREFIX}"
        INGRESS_HOST = "${params.INGRESS_HOST}"
        INGRESS_LOGIN_USER = "${params.INGRESS_LOGIN_USER}"
        INGRESS_LOGIN_PASSWORD = "${params.INGRESS_LOGIN_PASSWORD}"
        TEST_PHASE = "POST_INSTANTIATION"
        INGRESS_SAMPLE_APP_USER = "${params.INGRESS_SAMPLE_APP_USER}"
        INGRESS_GAS_USER = "${params.INGRESS_GAS_USER}"
    }

    // Stage names (with descriptions) taken from ADP Microservice CI Pipeline Step Naming Guideline: https://confluence.lmera.ericsson.se/pages/viewpage.action?pageId=122564754
    stages {
        stage('Clean') {
            steps {
                sh "rm -rf ./.aws ./.kube/ ./.cache/"
                archiveArtifacts allowEmptyArchive: true, artifacts: 'ci/postInstantiation.Jenkinsfile'
            }
        }
        stage('K6 Post Instantiation E2E Tests') {
            steps {
                sh "chmod 777 ./k6/scripts/run_k6_end2end_staging.sh"
                sh "./k6/scripts/run_k6_end2end_staging.sh"
            }
            post {
                always {
                    archiveArtifacts allowEmptyArchive: true, artifacts: 'k6/reports/k6-test-results.html'
                    archiveArtifacts allowEmptyArchive: true, artifacts: 'k6/reports/summary.json'
                    publishHTML([allowMissing: true,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: '',
                        reportFiles: 'k6/reports/k6-test-results.html',
                        reportName: 'K6 Test Results',
                        reportTitles: ''])
                }
            }
        }
    }
}

