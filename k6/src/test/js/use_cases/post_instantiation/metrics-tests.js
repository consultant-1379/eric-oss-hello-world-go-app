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

import { group, check, sleep } from 'k6'
import {
  SAMPLE_APP_HELLO_REQUEST_FAILURE_COUNT_PROMQL_QUERY,
  SAMPLE_APP_HELLO_REQUEST_COUNT_PROMQL_QUERY
} from '../../modules/constants.js'

import { logData } from '../../modules/common.js'
import * as sampleAppTest from './sample-app-tests.js'
import * as rbac from '../../modules/rbac.js'
import * as sampleApp from '../../modules/sampleApp.js'
import { isStatusOk } from '../../utils/validationUtils.js'

function verifyMetricsAvailable () {
  const k6sessionId = rbac.performK6UserLogin()
  let initialCountValue = 0
  group('Verify Sample App metrics available in pm-server', () => {
    check(
      sampleApp.getMetrics(
        k6sessionId,
        SAMPLE_APP_HELLO_REQUEST_COUNT_PROMQL_QUERY
      ),
      {
        'THEN Expect 200': r => isStatusOk(r.status),
        'AND Http Request Count Metric found': r => {
          const resultFound = r.json('data.result').length == 1
          if (resultFound) {
            initialCountValue = parseInt(
              r.json('data.result')[0].value[1].trim()
            )
          }
          return resultFound
        }
      }
    )
    check(
      sampleApp.getMetrics(
        k6sessionId,
        SAMPLE_APP_HELLO_REQUEST_FAILURE_COUNT_PROMQL_QUERY
      ),
      {
        'THEN Expect 200': r => isStatusOk(r.status),
        'AND Http Request Failure Count Metric found': r =>
          r.json('data.result').length == 1
      }
    )
  })
  rbac.clearSession()
  return initialCountValue
}

function awaitInitialScrape (k6sessionId, initialCountValue) {
  logData('Awaiting Initial Scrape')
  const end = new Date(Date.now() + 1 * 60000)
  let response = {}

  while (Date.now() < end.getTime()) {
    response = sampleApp.getMetrics(
      k6sessionId,
      SAMPLE_APP_HELLO_REQUEST_COUNT_PROMQL_QUERY
    )

    if (
      isStatusOk(response.status) &&
      parseInt(response.json('data.result')[0].value[1].trim()) >
        initialCountValue
    ) {
      break
    } else {
      sleep(5)
    }
  }
  return parseInt(response.json('data.result')[0].value[1].trim())
}

function verifyHttpRequestCountMetrics (initialCountValue) {
  const originalHttpRequestCount = awaitInitialScrape(
    rbac.performK6UserLogin(),
    initialCountValue
  )
  rbac.clearSession()

  group('Verify Sample App hello http request count', () => {
    sampleAppTest.verifyHelloEndpoint()
    logData('Sleeping for 60s to allow metrics to be scraped')
    sleep(60)

    check(
      sampleApp.getMetrics(
        rbac.performK6UserLogin(),
        SAMPLE_APP_HELLO_REQUEST_COUNT_PROMQL_QUERY
      ),
      {
        'THEN Expect 200': r => isStatusOk(r.status),
        'AND Http Request Count incremented correctly': r =>
          parseInt(r.json('data.result')[0].value[1].trim()) ==
          originalHttpRequestCount + 1
      }
    )
    rbac.clearSession()
  })
}

function verifyHttpRequestFailureCountMetrics () {
  group('Verify Sample App hello http request failure count', () => {
    check(
      sampleApp.getMetrics(
        rbac.performK6UserLogin(),
        SAMPLE_APP_HELLO_REQUEST_FAILURE_COUNT_PROMQL_QUERY
      ),
      {
        'THEN Expect 200': r => isStatusOk(r.status),
        'AND Http Request Failure Count is 0': r =>
          parseInt(r.json('data.result')[0].value[1].trim()) == 0
      }
    )
    rbac.clearSession()
  })
}

module.exports = {
  verifyMetricsAvailable,
  verifyHttpRequestCountMetrics,
  verifyHttpRequestFailureCountMetrics
}
