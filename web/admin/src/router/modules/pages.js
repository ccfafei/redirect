import Layout from '@/layout/index.vue'
import { createNameComponent } from '../createNode'
const route = [                                                                                         
  {  
    path: '/',
    component: Layout,
    redirect: '/pages/rulesSet',
    alwayShow: false,
    children: [
      {
            path: '',
            component: createNameComponent(() => import('@/views/pages/rulesSet/index.vue')),
            meta: { title: '规则设置', icon: 'el-icon-s-operation', cache: true, roles: ['admin'] },
          },
    ]
  },

  {
    path: '/userList',
    component: Layout,
    redirect: '/pages/userList',
    alwayShow: false,
    children: [
      {
        path: '',
        component: createNameComponent(() => import('@/views/pages/userList/index.vue')),
        meta: { title: '用户管理', icon: 'el-icon-user-solid', hideClose: true , roles: ['admin'] }
      }
    ]
  },
  {
    path: '/logs',
    component: Layout,
    redirect: '/pages/logs',
    alwayShow: false,
    children: [
      {
        path: '',
        component: createNameComponent(() => import('@/views/pages/logs/index.vue')),
        meta: { title: '日志管理', icon: 'el-icon-s-order', hideClose: true , roles: ['admin'] }
      }
    ]
  }
    
 
  
]

export default route