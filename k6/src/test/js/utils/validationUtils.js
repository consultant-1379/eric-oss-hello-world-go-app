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
const JSESSION_ID_PATTERN = /JSESSIONID=[\d\w-]+/
export function isStatusOk (status) {
  return status === 200
}
export function isStatusCreated (status) {
  return status === 201
}
export function isStatusNoContent (status) {
  return status === 204
}
export function isStatusNotFound (status) {
  return status === 404
}
export function isStatusForbidden (status) {
  return status === 403
}

export function isSessionValid (session) {
  return JSESSION_ID_PATTERN.test(session)
}
