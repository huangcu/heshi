
import VueCookies from 'vue-cookies'
// import TypeChecker from '../util/type-checker.js'

function routingGuard (to, from, next) {
  if (to.matched.length === 0) {
    // if route not found for 'to', back to 'from'
    if (from.matched.length === 0) {
      // if route for 'from' not found, back to HOME page
      next({ path: '/' })
    } else {
      // else back to 'from'
      next({ path: from.path })
    }
    console.warn('No Matched route found for: ')
    console.warn(to)
    return
  }

  let aRolesRequired = (to.meta.role)
  let isValid = validateAccessByRoles(to, aRolesRequired)
  if (!isValid) {
    if (to.path.indexOf('users') !== -1) {
      next({ path: '/users/login' })
    } else {
      next({path: '/manage/login'})
    }
    return
  }
  next()
}

function validateAccessByRoles (route, aRolesRequired) {
  let valid = true

  let userRoleType = null
  if (VueCookies.isKey('userprofile')) {
    let userprofile = JSON.parse(VueCookies.get('userprofile'))
    userRoleType = userprofile.user_type
  }
  // if (TypeChecker.isObject(store.getters.userProfile) &&
  // TypeChecker.isObject(store.getters.userProfile.userFe) &&
  // TypeChecker.isString(store.getters.userProfile.userFe.roleType)) {
  //   userRoleType = store.getters.userProfile.userFe.roleType.toLowerCase()
  // }
  // that will be valid if any one role is in the user's role list
  // Emily make the change 2018.JAN.30th
  let foundRoleInRequiredRoles = false
  if (aRolesRequired === undefined) {
    return valid
  }

  aRolesRequired.every(function (role) {
    if (userRoleType === role) {
      foundRoleInRequiredRoles = true
      return false
    }
    return true
  })

  return valid && foundRoleInRequiredRoles
}

export default routingGuard
