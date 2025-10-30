import service from '@/utils/request'

// 获取公告列表
export const getGMAnnouncementList = (data) => {
  return service({
    url: '/gm/announcement/list',
    method: 'post',
    data
  })
}


// 添加公告
export const addGMAnnouncement = (data) => {
  return service({
    url: '/gm/announcement/add',
    method: 'post',
    data
  })
}

// 删除公告
export const deleteGMAnnouncement = (data) => {
  return service({
    url: '/gm/announcement/delete',
    method: 'post',
    data
  })
}

// 更新公告
export const updateGMAnnouncement = (data) => {
  return service({
    url: '/gm/announcement/alter',
    method: 'post',
    data
  })
}

// 置顶公告
export const toppingGMAnnouncement = (data) => {
  return service({
    url: '/gm/announcement/topping',
    method: 'post',
    data
  })
}


