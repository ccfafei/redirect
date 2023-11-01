<template>
  <header>
    <div class="right-box">
  
      <!-- 用户信息 -->
      <div class="user-info">
        <span><i class="el-icon-user"></i> 您好，{{name}}</span> <el-button type='text' size="small" icon="el-icon-switch-button"  @click="loginOut">退出</el-button>
       
      </div>
     
    </div>
  </header>
</template>

<script>
import { defineComponent, computed, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
export default defineComponent({
  components: {
    
 
   
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    const layer = reactive({
      show: false,
      showButton: true
    })
  
    const name = computed(() =>  store.state.user.info) // 获取用户账号

    // const opendStateChange = () => {
    //   store.commit('app/isCollapseChange', !isCollapse.value)
    // }

    // login out the system
    const loginOut = () => {
      store.dispatch('user/loginOut')
    }
  
    return {
      
      layer,
      opendStateChange,
      loginOut,
      name
    }
  }
})
</script>

<style lang="scss" scoped>
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 60px;
    background-color: var(--system-header-background);
    padding-right: 22px;
  }
  .left-box {
    height: 100%;
    display: flex;
    align-items: center;
    .menu-icon {
      width: 60px;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 25px;
      font-weight: 100;
      cursor: pointer;
      margin-right: 10px;
      &:hover {
        background-color: var(--system-header-item-hover-color);
      }
      i {
        color: var(--system-header-text-color);
      }
    }
  }
  .right-box {
    display: flex;
    justify-content: center;
    align-items: center;
    .function-list{
      display: flex;
      .function-list-item {
        width: 30px;
        display: flex;
        justify-content: center;
        align-items: center;
        :deep(i) {
          color: var(--system-header-text-color);
        }
      }
    }
    .user-info {
      margin-left: 20px;
      font-size:14px;
      span{display: inline-block;margin-right:20px;}
      .el-dropdown-link {
        cursor: pointer;
        color: var(--system-header-breadcrumb-text-color);
      }
    }
  }
</style>