/*
 * COPYRIGHT Ericsson 2023
 *
 *
 *
 * The copyright to the computer program(s) herein is the property of
 *
 * Ericsson Inc. The programs may be used and/or copied only with written
 *
 * permission from Ericsson Inc. or in accordance with the terms and
 *
 * conditions stipulated in the agreement/contract under which the
 *
 * program(s) have been supplied.
 */
import { group, sleep } from 'k6'
import { SharedArray } from 'k6/data'

import * as metrics from './use_cases/post_instantiation/metrics-tests.js'
import * as gateway from './use_cases/pre_onboarding/gateway-tests.js'
import * as rbac from './use_cases/post_instantiation/rbac-tests.js'
import * as serviceExposure from './use_cases/post_instantiation/service-exposure-tests.js'
import * as sampleApp from './use_cases/post_instantiation/sample-app-tests.js'
import * as logging from './use_cases/post_instantiation/log-tests.js'

import { logData } from './modules/common.js'
import { DEFAULT_E2E_OPTIONS } from './modules/constants.js'
import { textSummary } from './modules/k6-summary.js'

import {
  getSampleAppRoutePayload,
  getSampleAppRbacPayload,
  getSampleAppUserPayload,
  getK6UserPayload,
  getSampleAppLoggingPayload
} from './utils/testDataUtils.js'
import { htmlReport } from 'https://arm1s11-eiffel004.eiffel.gic.ericsson.se:8443/nexus/content/sites/oss-sites/common/k6/eric-k6-static-report-plugin/latest/bundle/eric-k6-static-report-plugin.js'

export const options = DEFAULT_E2E_OPTIONS

const userData = new SharedArray(
  'Sample App user JSON',
  getSampleAppUserPayload
)
const k6UserData = new SharedArray(
  'K6 user JSON',
  getK6UserPayload
)
const rbacData = new SharedArray(
  'Sample App RBAC JSON',
  getSampleAppRbacPayload
)
const routeData = new SharedArray(
  'Sample App route JSON',
  getSampleAppRoutePayload
)
const loggingData = new SharedArray(
  'Sample App Logging JSON',
  getSampleAppLoggingPayload
)

export default function () {
  if (__ENV.TEST_PHASE === 'PRE_ONBOARDING') {
    logData('PRE_ONBOARDING')
    group('GIVEN The API Gateway is available', () => {
      gateway.verifySessionCreation()
    })
  } else {
    logData('POST_INSTANTIATION')
    group('GIVEN The API Gateway is available', () => {
      gateway.verifySessionCreation()
    })

    group('GIVEN The K6 User is created', () => {
      rbac.verifyK6User(k6UserData[0])
    })

    group('GIVEN The Service Exposure API is available', () => {
      serviceExposure.verifySampleAppRoute(routeData[0])
    })

    group('GIVEN The User Administration API is available', () => {
      rbac.verifySampleAppRbac(rbacData[0])
      rbac.verifySampleAppUser(userData[0])
    })

    logData('Sleeping for 60s to allow RBAC to take effect')
    sleep(60)

    group('GIVEN SampleApp RBAC is in place', () => {
      rbac.verifyRbacEnforced()
      sampleApp.verifyHelloEndpoint()
    })

    logData('Sleeping for 10s to allow Logs to be registered')
    sleep(10)

    group('GIVEN Log Consumer is in place', () => {
      logging.verifyLogsStreamed(loggingData[0])
    })

    group('GIVEN App Metrics is in place', () => {
      const initialCountValue = metrics.verifyMetricsAvailable()
      metrics.verifyHttpRequestCountMetrics(initialCountValue)
      metrics.verifyHttpRequestFailureCountMetrics()
    })

    group('GIVEN tests are finished', () => {
      rbac.verifyCleanUpSampleAppRbac(rbacData[0])
      serviceExposure.verifyCleanUpSampleAppRoute()
      rbac.verifyCleanUpK6User()
    })
  }
}

export function handleSummary (data) {
  const reportPath = '/k6/reports/'
  let result = { stdout: textSummary(data) }
  result[reportPath.concat('k6-test-results.html')] = htmlReport(data)
  result[reportPath.concat('summary.json')] = JSON.stringify(data)
  return result
}
