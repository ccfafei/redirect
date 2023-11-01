<template>
  <div class="container">
    <div class="box">
      <h1>{{ systemTitle }}</h1>
      <el-form class="form">
        <el-input
          size="large"
          v-model="form.account"
          placeholder="用户名"
          type="text"
          maxlength="50"
        >
          <template #prepend>
            <i class="sfont system-xingmingyonghumingnicheng"></i>
          </template>
        </el-input>
        <el-input
          size="large"
          ref="password"
          v-model="form.password"
          :type="passwordType"
          placeholder="密码"
          name="password"
          maxlength="50"
        >
          <template #prepend>
            <i class="sfont system-mima"></i>
          </template>
          <template #append>
            <i class="sfont password-icon" :class="passwordType ? 'system-yanjing-guan': 'system-yanjing'" @click="passwordTypeChange"></i>
          </template>
        </el-input>
         <el-input
          size="large"
          v-model="form.captcha_text"
          placeholder="验证码"
          type="text"
          maxlength="50"
           @keyup.enter="submit"
        >
          <template #append >
         <img :src="captchaDataImg" style="width:100px;height:35px;cursor: pointer;" @click="captcha()">
        </template>

        </el-input>
        <el-button type="primary" @click="submit" style="width: 100%;" size="medium">登录</el-button>
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
import { getCaptchaId } from "@/api/user"
export default defineComponent({
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    const form = reactive({
      account: '',
      password: '',
      captcha_text :'',
      captcha_id :''
    })
  const captchaDataImg = ref(null) 
  const captcha = () => {
      getCaptchaId()
      .then(res => { 
       captchaDataImg.value = res.result.captcha_data
       form.captcha_id = res.result.captcha_id
      })
      .catch(error => {
       
      })
      .finally(() => {
        
      })
    }
    captcha()
    const passwordType = ref('password')
    const passwordTypeChange = () => {
      passwordType.value === '' ? passwordType.value = 'password' : passwordType.value = ''
    }
    const checkForm = () => {
      return new Promise((resolve, reject) => {
        if (form.account === '') {
          ElMessage.warning({
            message: '用户名不能为空',
            type: 'warning'
          });
          return;
        }
        if (form.password === '') {
          ElMessage.warning({
            message: '密码不能为空',
            type: 'warning'
          })
          return;
        }
         if (form.captcha_text === '') {
          ElMessage.warning({
            message: '验证码不能为空',
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
        let params = {
          account: form.account,
          password: form.password,
          captcha_text :form.captcha_text,
          captcha_id :form.captcha_id
        }
       
        store.dispatch('user/login', params)
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
      passwordType,
      passwordTypeChange,
      submit,
      captchaDataImg,
      captcha
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
      height: 400px;
      overflow: hidden;
      box-shadow: 0 6px 20px 5px rgba(152, 152, 152, 0.1), 0 16px 24px 2px rgba(117, 117, 117, 0.14);
      h1 {
        margin-top:40px;
        text-align: center;
        font-size:1.4em;
      }
      .form {
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