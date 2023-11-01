<template>
  <div class="container">
    <div class="box">
      <!-- <h3>{{ systemTitle }}</h3> -->
      <el-form class="form">
      
        <el-input
          size="large"
          ref="password"
          v-model="form.password"
          placeholder="密码"
          name="password"
          maxlength="50"
          @keyup.enter="submit"
       
        >
         
        </el-input>
      <el-input hidden/>  
        <el-button type="primary" @click="submit"  size="medium" style="position:absolute;right:2px;top:2px;">进入</el-button>
      </el-form>
    </div>
  </div>
</template>

<script>
import { systemTitle } from '@/config'
import { defineComponent, ref, reactive  } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import { addRoutes } from '@/router'
import { ElMessage } from 'element-plus'

export default defineComponent({
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
   const form = ref({
      rule_id: '',
      password: '',
   
    })

    const GetRequest=() =>{
 
    const url = decodeURI(window.location.href)
    if (url.indexOf("?") != -1)//url中存在问号，也就说有参数。  
    {
       const str = url.split("?");  //得到?后面的字符串
       form.value.rule_id=parseInt(str[1].split("=")[1])
       
    }

   
  }
    GetRequest()
	
    const checkForm = () => {
      
      return new Promise((resolve, reject) => {
       
        if (form.value.password === '') {
          ElMessage.warning({
            message: '密码不能为空',
            type: 'warning'
          })
          return;
        }
        
        resolve(true)
      })
    }
    
    const submit = () => {
      
      checkForm()
      .then(() => {
        store.dispatch('user/login', form.value)
        .then(() => {
          ElMessage.success({
            message: '登录成功',
            type: 'success',
            showClose: true,
            duration: 1000
          })
          
         addRoutes() 
         
        router.push(route.query.redirect || '/')
  
        })

      })
    }
   
    return {
      systemTitle,
      form,
   
      submit,
    
    }
  }
})
</script>

<style lang="scss" scoped>
  .container {
    position: relative;
    width: 100vw;
    height: 100vh;
    background-color: #eef0f3;
    .box {
      width: 500px;
      position: absolute;
      left: 50%;
      top: 50%;
      background: white;
      border-radius: 8px;
      transform: translate(-50%, -50%);
      height: 140px;
      overflow: hidden;
      // padding:20px 0;
      box-shadow: 0 6px 20px 5px rgba(152, 152, 152, 0.1), 0 16px 24px 2px rgba(117, 117, 117, 0.14);
      h1 {
        margin-top:60px;
        text-align: center;
        font-size:1.8em;
      }
      .form {
         position: relative;
        width: 80%;
        margin: 50px auto 15px;
        .el-input {
          margin-bottom: 20px;
        }
        .password-icon {
          cursor: pointer;
          color: #409EFF;
        }
      }
      .fixed-top-right {
        position: absolute;
        top: 10px;
        right: 10px;
      }
    }
  }
  @media screen and ( max-width: 750px ) {
    .container .box {
      width: 100vw;
      height: 100vh;
      box-shadow: none;
      left: 0;
      top: 0;
      transform: none;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      h1 {
        margin-top: 0;
      }
      .form {
         margin: 30px auto 15px;
      }
    }
  }
</style>