import request from '@/utils/request'

// const req = {
//   index: (query) => {
//     return request({
//       url: '/x/img/index',
//       method: 'get',
//       params: query
//     })
//   },
//   detail: () => {
//     return request({
//       url: '/x/img/detail',
//       method: 'get'
//     })
//   },
//   create: (data) => {
//     return request({
//       url: '/x/img/create',
//       method: 'post',
//       data
//     })
//   },
//   update: (data) => {
//     return request({
//       url: '/x/img/update',
//       method: 'put',
//       data
//     })
//   }
// }

// export default req

export function fetchIndex(query) {
  return request({
    url: '/x/img/index',
    method: 'get',
    params: query
  })
}

export function fetchDetail() {
  return request({
    url: '/x/img/detail',
    method: 'get'
  })
}

export function createImg(data) {
  return request({
    url: '/x/img/create',
    method: 'post',
    data
  })
}

export function updateImg(data) {
  return request({
    url: '/x/img/update',
    method: 'put',
    data
  })
}
