import Vue from 'vue'
import Router from 'vue-router'
import linkChain from '@/components/linkchain'
import nodeList from '@/components/nodelist'
import nodeInfo from '@/components/nodeinfo'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/',
      name: 'linkchain',
      component: linkChain,
      meta: {
        title: '链状态'
      }
    },
    {
      path: '/node/list',
      name: 'node-list',
      component: nodeList,
      meta: {
        title: '节点列表'
      }
    },
    {
      path: '/node/info/:id',
      name: 'node-status',
      component: nodeInfo,
      meta: {
        title: '兑豆豆数据监控系统'
      }
    }
  ]
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title
  next()
})
export default router
