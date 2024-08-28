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
import * as rbac from '../../modules/rbac.js'
import { logData } from '../../modules/common.js'
import {
  isStatusCreated,
  isStatusNoContent,
  isStatusNotFound
} from '../../utils/validationUtils.js'

function verifySampleAppRoute (routeData) {
  const k6sessionId = rbac.performK6UserLogin()
  group('Delete Sample App Route if exists', () => {
    check(sampleApp.deleteSampleAppRoute(), {
      'THEN Expect 404 OR 204': r => {
        logData(`Delete Route Actual returned status: ${r.status}`)
        return isStatusNotFound(r.status) || isStatusNoContent(r.status)
      }
    })
  })

  group('Create Sample App Route', () => {
    check(
      sampleApp.createSampleAppRoute(k6sessionId, routeData),
      {
        'THEN Expect 201': r => isStatusCreated(r.status)
      }
    )
  })
  rbac.clearSession()
}

function verifyCleanUpSampleAppRoute () {
  rbac.performK6UserLogin()

  group('Delete Sample App Route', () => {
    check(sampleApp.deleteSampleAppRoute(), {
      'THEN Expect 204': r => isStatusNoContent(r.status)
    })
  })

  rbac.clearSession()
}

module.exports = {
  verifySampleAppRoute,
  verifyCleanUpSampleAppRoute
}
