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

import {
  httpGet,
  httpPostJson,
  logData,
  httpDelete,
  createSession
} from './common.js'
import {
  INGRESS_URL,
  INGRESS_SAMPLE_APP_USER_PARAMS,
  INGRESS_LOGIN_URI,
  APPLICATION_JSON,
  SAMPLE_APP_ROUTE_ID,
  SAMPLE_APP_HELLO_URI,
  INGRESS_ROUTES_URI,
  INGRESS_LOGGING_URI,
  INGRESS_METRIC_QUERY_URI
} from './constants.js'

let userSession

function performSampleAppUserLogin (options = {}) {
  if (!userSession) {
    userSession = createSession(
      INGRESS_URL,
      INGRESS_LOGIN_URI,
      '',
      INGRESS_SAMPLE_APP_USER_PARAMS,
      options
    )
    logData('CREATE USER SESSION', INGRESS_SAMPLE_APP_USER_PARAMS)
    logData('CREATE USER SESSION', userSession)
  }
  return userSession
}

function getRouteSessionParams (k6sessionId) {
  return {
    headers: {
      'Content-Type': APPLICATION_JSON,
      Cookie: k6sessionId
    }
  }
}

function cleanUserSession () {
  userSession = undefined
}

function getHello () {
  return httpGet(INGRESS_URL, SAMPLE_APP_HELLO_URI)
}

function createSampleAppRoute (k6sessionId, routePayload) {
  return httpPostJson(
    INGRESS_URL,
    INGRESS_ROUTES_URI,
    routePayload,
    getRouteSessionParams(k6sessionId)
  )
}

function getLogs (k6sessionId, loggingPayload) {
  return httpPostJson(
    INGRESS_URL,
    INGRESS_LOGGING_URI,
    loggingPayload,
    getRouteSessionParams(k6sessionId)
  )
}

function getMetrics (k6sessionId, query) {
  return httpGet(
    INGRESS_URL,
    INGRESS_METRIC_QUERY_URI + query,
    getRouteSessionParams(k6sessionId)
  )
}

function deleteSampleAppRoute () {
  return httpDelete(INGRESS_URL, INGRESS_ROUTES_URI + '/' + SAMPLE_APP_ROUTE_ID)
}

module.exports = {
  performSampleAppUserLogin,
  cleanUserSession,
  getHello,
  createSampleAppRoute,
  deleteSampleAppRoute,
  getLogs,
  getMetrics
}
