
import VueCookies from 'vue-cookies'

function routingGuard (to, from, next) {
  // console.log(to)
  // console.log('=====================')
  // console.log(from)
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

  // if route found for 'to'; if it is myaccount, check role
  if (to.path.indexOf('/myaccount') !== -1) {
    // if login,
    if (VueCookies.isKey('userprofile')) {
      let userprofile = JSON.parse(VueCookies.get('userprofile'))
      // if login, and user type correct, next()
      console.log(to.meta.role)
      if (to.meta.role.indexOf(userprofile.user_type) !== -1) {
        next()
      } else {
        // else back to 'from'
        if (from.matched.length === 0) {
          // if route for 'from' not found, back to HOME page
          next({ path: '/' })
        } else {
          // else back to 'from'
          next({ path: from.path })
        }
      }
    } else {
      // if not login, cookie has no 'userprofile', back to 'login'
      next({path: '/login'})
    }
  } else {
    // if it not myaccount, no need to check role, next()
    next()
  }
}

export default routingGuard
