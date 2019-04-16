import request from './config'

export function getLinkChain () {
  return request({
    url: '/linkchian/overview',
    method: 'get'
  })
}
export function getBusiness () {
  return request({
    url: '/business/overview',
    method: 'get'
  })
}
export function getBusinessList () {
  return request({
    url: '/business/list',
    method: 'get'
  })
}
export function getBusinessInfo (id) {
  return request({
    url: `/business/info/${id}`,
    method: 'get'
  })
}
export function getNode () {
  return request({
    url: '/linkchian/nodes',
    method: 'get'
  })
}
export function getNodeInfo (ip) {
  return request({
    url: `/linkchian/authnodes/${ip}`,
    method: 'get'
  })
}
