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

import { group, check } from 'k6'

import * as sampleApp from '../../modules/sampleApp.js'
import { isStatusOk } from '../../utils/validationUtils.js'

function verifyHelloEndpoint () {
  sampleApp.performSampleAppUserLogin()

  group('GET hello and verify response, logged in as Sample App User', () => {
    check(sampleApp.getHello(), {
      'THEN Expect 200': r => isStatusOk(r.status),
      'AND response body contains text "Hello World!!"': r =>
        r.body.includes('Hello World!!')
    })
  })

  sampleApp.cleanUserSession()
}

module.exports = {
  verifyHelloEndpoint
}
