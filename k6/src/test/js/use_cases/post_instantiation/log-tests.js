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

import * as rbac from '../../modules/rbac.js'
import * as sampleApp from '../../modules/sampleApp.js'
import { isStatusOk } from '../../utils/validationUtils.js'

function verifyLogsStreamed (logData) {
  group('Verify Logs have been streamed', () => {
    check(sampleApp.getLogs(rbac.performK6UserLogin(), logData), {
      'THEN Expect 200': r => isStatusOk(r.status),
      'AND response body contains Sample App Logs from the last 5 minutes': r =>
        r.json('hits.hits').length > 0
    })
  })

  rbac.clearSession()
}

module.exports = {
  verifyLogsStreamed
}
