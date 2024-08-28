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

export const INGRESS_SCHEMA = __ENV.INGRESS_SCHEMA
  ? __ENV.INGRESS_SCHEMA
  : 'https'
export const INGRESS_HOST = __ENV.INGRESS_HOST
  ? __ENV.INGRESS_HOST
  : 'gas.stsvp1eic26.stsoss.sero.gic.ericsson.se'
export const INGRESS_URL = INGRESS_SCHEMA.concat('://').concat(INGRESS_HOST)
// EIC URIs
export const INGRESS_LOGIN_URI = '/auth/v1/login'
export const INGRESS_ROUTES_URI = '/v1/routes'
export const INGRESS_ROUTES_RBAC_URI = '/idm/rolemgmt/v1/extapp/rbac'
export const INGRESS_ROUTES_USER_URI = '/idm/usermgmt/v1/users'
export const INGRESS_LOGGING_URI = '/_search'
export const INGRESS_METRIC_QUERY_URI = '/metrics/viewer/api/v1/query?query='
// EIC Users
export const INGRESS_LOGIN_USER = __ENV.INGRESS_LOGIN_USER
  ? __ENV.INGRESS_LOGIN_USER
  : 'sys-user'
export const INGRESS_GAS_USER = __ENV.INGRESS_GAS_USER
  ? __ENV.INGRESS_GAS_USER
  : 'gas-user'
export const INGRESS_K6_USER = __ENV.INGRESS_K6_USER
  ? __ENV.INGRESS_K6_USER
  : 'k6-go-app-admin-user'
export const INGRESS_SAMPLE_APP_USER = __ENV.INGRESS_SAMPLE_APP_USER
  ? __ENV.INGRESS_SAMPLE_APP_USER
  : 'sample-app-user'
export const INGRESS_LOGIN_PASSWORD = __ENV.INGRESS_LOGIN_PASSWORD
  ? __ENV.INGRESS_LOGIN_PASSWORD
  : 'idunEr!css0n'
//Test Params
export const X_TENANT = 'master'
export const APPLICATION_FORM_URL_ENCODED = 'application/x-www-form-urlencoded'
export const APPLICATION_JSON = 'application/json'
export const SAMPLE_APP_HELLO_URI = '/hello'
export const SAMPLE_APP_ROUTE_ID = 'hello-route-001'
export const SAMPLE_APP_HELLO_REQUEST_COUNT_PROMQL_QUERY =
  'hello_world_requests_total'
export const SAMPLE_APP_HELLO_REQUEST_FAILURE_COUNT_PROMQL_QUERY =
  'hello_world_requests_failed_total'

export const DEFAULT_TIMEOUT = 60
export const MAX_RETRY = 10

export const DEFAULT_E2E_OPTIONS = {
  duration: '30m',
  vus: 1,
  iterations: 1,
  thresholds: {
    checks: ['rate == 1.0']
  }
}

const getLoginParams = (user, password = INGRESS_LOGIN_PASSWORD) => ({
  headers: {
    'X-Login': user,
    'X-password': INGRESS_LOGIN_PASSWORD,
    'X-tenant': X_TENANT,
    'Content-Type': APPLICATION_FORM_URL_ENCODED
  }
})

export const INGRESS_LOGIN_PARAMS = getLoginParams(INGRESS_LOGIN_USER)
export const INGRESS_GAS_LOGIN_PARAMS = getLoginParams(INGRESS_GAS_USER)
export const INGRESS_K6_LOGIN_PARAMS = getLoginParams(INGRESS_K6_USER)
export const INGRESS_SAMPLE_APP_USER_PARAMS = getLoginParams(
  INGRESS_SAMPLE_APP_USER
)
