// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import './style/index.css'
import './style/media.css'
import 'babel-polyfill'
import { Table, TableColumn, Loading } from 'element-ui'
import filter from './filters/filter'

Vue.config.productionTip = false
// 引入全局过滤器
Vue.filter('getFormatDate', filter.getFormatDate)
Vue.filter('getFormatNumber', filter.getFormatNumber)
Vue.filter('getPercent', filter.getPercent)
Vue.filter('getFormatDay', filter.getFormatDay)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Loading)

// 路由回到顶部
router.beforeEach((to, from, next) => {
  window.scrollTo(0, 0)
  if (to.meta === from.meta) {}
  next()
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
