//This is where all the api related config resides on
// i.e: Cookies management, axios global config
import axios from 'axios'
import TypeChecker from './type-checker.js'
import router from '../router'
import StoreManager from './store-manager.js'

let CancelToken = axios.CancelToken

let axiosHelper = {
  createAxios: function (obj) {
    if (obj.baseURL) {
      let instance = axios.create({
        baseURL: obj.baseURL,
        timeout: obj.timeout ? obj.timeout : 30000
      });
      instance.interceptors.request.use(function (config) {
        return config;
      }, function (err) {
        return Promise.reject(err);
      });
      instance.interceptors.response.use(function (res) {
        return res;
      }, function (err) {
        if (TypeChecker.isObject(err.response) && TypeChecker.isObject(err.response.data) && TypeChecker.isString(err.response.data.message)) {
          console.error(err.response.data.message);
        } else {
          console.error(err.message);
        }

        if (err.response && (err.response.status == 401 || err.response.status == 403)) {
          // UserServices.clearCurrentUser();
          const storageKey = 'logined-user';
          const tokenCookieOpt = {
            path: '/services/api/v1',
            secure: false
          };
          this.$cookies.remove(storageKey);
          this.$cookies.remove(tokenKey, { path: tokenCookieOpt.path });
          delete axios.defaults.headers.common['Authorization'];

          // ProcessSelectionService.clearProcessSelection();
          const storageKey3 = 'logined-process';
          this.$cookies.remove(storageKey3);

          // router redirrect
          router.replace('/passport/login');

        } else {
          return Promise.reject(err);
        }
      });
      return instance;
    }
  }
}


export { axios, cookies, common, axiosHelper }
