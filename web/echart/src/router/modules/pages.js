import Layout from '@/layout/index.vue'
import { createNameComponent } from '../createNode'
const route = [                                                                                         
  {  
    path: '/',
    component: Layout,
    redirect: '/pages/index',
    alwayShow: false,
    children: [
      {
            path: '',
            component: createNameComponent(() => import('@/views/pages/index/index.vue')),
            meta: { title: '首页', icon: 'el-icon-s-operation'},
          },
    ]
  },
  
]

export default route