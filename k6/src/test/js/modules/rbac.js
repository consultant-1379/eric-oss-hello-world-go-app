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
  createSession,
  httpDelete,
  httpGet,
  httpPostJson,
  logData
} from './common.js'
import {
  INGRESS_URL,
  INGRESS_LOGIN_URI,
  INGRESS_ROUTES_USER_URI,
  INGRESS_GAS_LOGIN_PARAMS,
  INGRESS_ROUTES_RBAC_URI,
  APPLICATION_JSON,
  INGRESS_SAMPLE_APP_USER,
  INGRESS_K6_LOGIN_PARAMS,
  INGRESS_K6_USER
} from './constants.js'

let sessionId

function performGasUserLogin (options = {}) {
  if (!sessionId) {
    sessionId = createSession(
      INGRESS_URL,
      INGRESS_LOGIN_URI,
      '',
      INGRESS_GAS_LOGIN_PARAMS,
      options
    )
    logData('CREATE SESSION', sessionId)
  }
  return sessionId
}

function performK6UserLogin (options = {}) {
  if (!sessionId) {
    sessionId = createSession(
      INGRESS_URL,
      INGRESS_LOGIN_URI,
      '',
      INGRESS_K6_LOGIN_PARAMS,
      options
    )
    logData('CREATE SESSION', sessionId)
  }
  return sessionId
}

function getSessionParams (relevantSessionId) {
  return {
    headers: {
      'Content-Type': APPLICATION_JSON,
      Accept: APPLICATION_JSON,
      Cookie: relevantSessionId
    }
  }
}

function clearSession () {
  sessionId = undefined
}

function getSearchQuery (user) {
  return `?\&search=(username==*${user}*;tenantname==master)`
}

function getDeleteUri (user) {
  return `/${user}?tenantname=master`
}

function deleteSampleAppRbac (rbacPayload, options = {}) {
  return httpDelete(
    INGRESS_URL,
    INGRESS_ROUTES_RBAC_URI,
    rbacPayload,
    getSessionParams(performK6UserLogin()),
    options
  )
}

function createSampleAppRbac (rbacPayload, options = {}) {
  return httpPostJson(
    INGRESS_URL,
    INGRESS_ROUTES_RBAC_URI,
    rbacPayload,
    getSessionParams(performK6UserLogin()),
    options
  )
}

function getSampleAppUser (options = {}) {
  return httpGet(
    INGRESS_URL,
    INGRESS_ROUTES_USER_URI.concat(getSearchQuery(INGRESS_SAMPLE_APP_USER)),
    getSessionParams(performGasUserLogin()),
    options
  )
}

function getK6User (options = {}) {
  return httpGet(
    INGRESS_URL,
    INGRESS_ROUTES_USER_URI.concat(getSearchQuery(INGRESS_K6_USER)),
    getSessionParams(performGasUserLogin()),
    options
  )
}

function createUser (userPayload, options = {}) {
  return httpPostJson(
    INGRESS_URL,
    INGRESS_ROUTES_USER_URI,
    userPayload,
    getSessionParams(performGasUserLogin()),
    options
  )
}

function deleteSampleAppUser (options = {}) {
  return httpDelete(
    INGRESS_URL,
    INGRESS_ROUTES_USER_URI.concat(getDeleteUri(INGRESS_SAMPLE_APP_USER)),
    undefined,
    getSessionParams(performGasUserLogin()),
    options
  )
}

function deleteK6User (options = {}) {
  return httpDelete(
    INGRESS_URL,
    INGRESS_ROUTES_USER_URI.concat(getDeleteUri(INGRESS_K6_USER)),
    undefined,
    getSessionParams(performGasUserLogin()),
    options
  )
}

module.exports = {
  performGasUserLogin,
  performK6UserLogin,
  deleteSampleAppRbac,
  deleteK6User,
  createSampleAppRbac,
  getSampleAppUser,
  getK6User,
  createUser,
  deleteSampleAppUser,
  clearSession
}
