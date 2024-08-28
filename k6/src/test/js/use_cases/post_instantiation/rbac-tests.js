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
import { logData } from '../../modules/common.js'
import { INGRESS_SAMPLE_APP_USER, INGRESS_K6_USER } from '../../modules/constants.js'
import {
  isStatusOk,
  isStatusNoContent,
  isStatusNotFound,
  isSessionValid,
  isStatusForbidden
} from '../../utils/validationUtils.js'

function verifyK6User(userData) {
  const response = rbac.getK6User()
  const user = JSON.parse(response.body).find(
    user => user.username === INGRESS_K6_USER
  )

  if (!!user) {
    group(`Delete K6 User, if exists`, () => {
      check(rbac.deleteK6User(), {
        'THEN Expect 204': r => isStatusNoContent(r.status)
      })
    })
  }

  group(`Create K6 User`, () => {
    check(rbac.createUser(userData), {
      'THEN Expect 200': r => isStatusOk(r.status)
    })
  })

  rbac.clearSession()
}

function verifySampleAppRbac (rbacData) {
  group('Delete Sample App RBAC if exists', () => {
    check(rbac.deleteSampleAppRbac(rbacData), {
      'THEN Expect 404 OR 204': r => {
        logData(`Delete RBAC Actual returned status: ${r.status}`)
        return isStatusNotFound(r.status) || isStatusNoContent(r.status)
      }
    })
  })

  group('Create Sample App RBAC', () => {
    check(rbac.createSampleAppRbac(rbacData), {
      'THEN Expect 200': r => isStatusOk(r.status)
    })
  })

  rbac.clearSession()
}

function verifySampleAppUser (userData) {
  const response = rbac.getSampleAppUser()
  const user = JSON.parse(response.body).find(
    user => user.username === INGRESS_SAMPLE_APP_USER
  )

  if (!!user) {
    group(`Delete Sample App User, if exists`, () => {
      check(rbac.deleteSampleAppUser(), {
        'THEN Expect 204': r => isStatusNoContent(r.status)
      })
    })
  }

  group(`Create Sample App User`, () => {
    check(rbac.createUser(userData), {
      'THEN Expect 200': r => isStatusOk(r.status)
    })
  })

  group(`Login as Sample App User`, () => {
    check(sampleApp.performSampleAppUserLogin(), {
      'THEN JSESSIONID is created and valid': r => isSessionValid(r)
    })
  })

  rbac.clearSession()
  sampleApp.cleanUserSession()
}

function verifyRbacEnforced () {
  rbac.performK6UserLogin()

  group('GET hello, logged in as K6 Admin User', () => {
    check(sampleApp.getHello(), {
      'THEN Expect 403': r => isStatusForbidden(r.status)
    })
  })

  rbac.clearSession()

  sampleApp.performSampleAppUserLogin()
  group('GET hello, logged in as Sample App User', () => {
    check(sampleApp.getHello(), {
      'THEN Expect 200': r => isStatusOk(r.status)
    })
  })

  sampleApp.cleanUserSession()
}

function verifyCleanUpSampleAppRbac (rbacData) {
  group('Delete Sample App RBAC', () => {
    check(rbac.deleteSampleAppRbac(rbacData), {
      'THEN Expect 204': r => isStatusNoContent(r.status)
    })
  })

  rbac.clearSession()

  group('Delete Sample App User', () => {
    check(rbac.deleteSampleAppUser(), {
      'THEN Expect 204': r => isStatusNoContent(r.status)
    })
  })

  rbac.clearSession()
}

function verifyCleanUpK6User () {
  group('Delete K6 User', () => {
    check(rbac.deleteK6User(), {
      'THEN Expect 204': r => isStatusNoContent(r.status)
    })
  })

  rbac.clearSession ()
}

module.exports = {
  verifyK6User,
  verifySampleAppRbac,
  verifySampleAppUser,
  verifyCleanUpSampleAppRbac,
  verifyCleanUpK6User,
  verifyRbacEnforced
}
