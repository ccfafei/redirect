<template>
  <header>
    <div class="left-box">
      <!-- 收缩按钮 -->
      <div class="menu-icon" @click="opendStateChange">
        <i :class="isCollapse ? 'el-icon-s-unfold' : 'el-icon-s-fold'"></i>
      </div>
      <Breadcrumb />
    </div>
    <div class="right-box">
      <!-- 快捷功能按钮 -->
      <div class="function-list">
        <!-- <div class="function-list-item hidden-sm-and-down"><Full-screen /></div>
        <div class="function-list-item"><SizeChange /></div> -->
        <div class="function-list-item hidden-sm-and-down"><Theme /></div>
      </div>
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
import FullScreen from './functionList/fullscreen.vue'
import SizeChange from './functionList/sizeChange.vue'
import Theme from './functionList/theme.vue'
import Breadcrumb from './Breadcrumb.vue'

export default defineComponent({
  components: {
    FullScreen,
    Breadcrumb,
    SizeChange,
    Theme,
   
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
    const isCollapse = computed(() => store.state.app.isCollapse)
    // isCollapse change to hide/show the sidebar
    const opendStateChange = () => {
      store.commit('app/isCollapseChange', !isCollapse.value)
    }

    // login out the system
    const loginOut = () => {
      store.dispatch('user/loginOut')
    }
    
   
    return {
      isCollapse,
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