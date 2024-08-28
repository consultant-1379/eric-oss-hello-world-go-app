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
  INGRESS_LOGIN_PASSWORD,
  INGRESS_SAMPLE_APP_USER,
  INGRESS_K6_USER,
  SAMPLE_APP_HELLO_URI,
  SAMPLE_APP_ROUTE_ID
} from '../modules/constants.js'

export function getSampleAppRoutePayload () {
  const sampleAppRoutePayload = JSON.parse(open('../resources/route.json'))
  sampleAppRoutePayload.id = SAMPLE_APP_ROUTE_ID
  sampleAppRoutePayload.predicates[0].args._genkey_0 = SAMPLE_APP_HELLO_URI
  return [sampleAppRoutePayload]
}

export function getSampleAppUserPayload () {
  const sampleAppUserPayload = JSON.parse(open('../resources/user.json'))
  sampleAppUserPayload.password = INGRESS_LOGIN_PASSWORD
  sampleAppUserPayload.user.username = INGRESS_SAMPLE_APP_USER
  sampleAppUserPayload.user.privileges = ['SampleApp_Application_Administrator']
  return [sampleAppUserPayload]
}

export function getK6UserPayload () {
  const samplek6UserData = JSON.parse(open('../resources/k6-user.json'))
  samplek6UserData.password = INGRESS_LOGIN_PASSWORD
  samplek6UserData.user.username = INGRESS_K6_USER
  samplek6UserData.user.privileges = 
  ['MetricsViewer',
   'LogAPI_ExtApps_Application_ReadOnly',
   'Exposure_Application_Administrator',
   'UserAdministration_ExtAppRbac_Application_SecurityAdministrator']
  return [samplek6UserData]
}

export function getSampleAppRbacPayload () {
  const sampleAppRbacPayload = JSON.parse(
    open('../resources/sample_app_rbac.json')
  )
  sampleAppRbacPayload.authorization.resources[0].uris = [SAMPLE_APP_HELLO_URI]
  return [sampleAppRbacPayload]
}

export function getSampleAppLoggingPayload () {
  const sampleAppLoggingPayload = JSON.parse(
    open('../resources/log_query.json')
  )
  let lte = new Date(Date.now() + 5 * 60000)
  let gte = new Date()
  sampleAppLoggingPayload.query.bool.must[1].range.timestamp.gte =
    gte.toISOString().replace('Z', '') + getISOOffset(gte)
  sampleAppLoggingPayload.query.bool.must[1].range.timestamp.lte =
    lte.toISOString().replace('Z', '') + getISOOffset(lte)
  return [sampleAppLoggingPayload]
}

export function getISOOffset (date) {
  date = date
  function z (n) {
    return ('0' + n).slice(-2)
  }
  var offset = date.getTimezoneOffset()
  var sign = offset < 0 ? '+' : '-'
  if (offset == 0) {
    sign = '+'
  }
  offset = Math.abs(offset)
  return sign + z((offset / 60) | 0) + ':' + z(offset % 60)
}
