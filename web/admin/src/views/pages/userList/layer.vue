<template>
  <Layer :layer="layer" @confirm="submit" ref="layerDom">
    <el-form :model="form" :rules="rules" ref="userForm" label-width="120px" style="margin-right:30px;">
      <el-form-item label="账号：" prop="account">
        <el-input v-model="form.account" placeholder="请输入账号"></el-input>
      </el-form-item>
      <el-form-item label="密码：" prop="password">
        <el-input v-model="form.password" placeholder="请输入密码"></el-input>
      </el-form-item>
       <el-form-item label="名称：" prop="name">
        <el-input v-model="form.name" placeholder="请输入名称"></el-input>
      </el-form-item>
     
    
    </el-form>
  </Layer>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { addUser, updateUser } from '@/api/user'

import Layer from '@/components/layer/index.vue'
export default defineComponent({
  components: {
    Layer
  },
  props: {
    layer: {
      type: Object,
      default: () => {
        return {
          show: false,
          title: '',
          showButton: true
        }
      }
    }
  },
  setup(props, ctx) {
    const userForm= ref(null)
    const layerDom = ref(null)
     let form = ref({
      name: '',
      account:'',
      password:''
    })
    const rules = {
      password: [{ required: true, message: '请输入密码，不少于6位', trigger: 'blur' }],
      account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
    
    }
    init()
    function init() { // 用于判断新增还是编辑功能
      if (props.layer.row) {
        form.value = JSON.parse(JSON.stringify(props.layer.row)) // 数量量少的直接使用这个转
      } else {

      }
    }
    return {
      form,
      rules,
      layerDom,
      userForm,
   
    }
  },
  methods: {
    submit() {
      if (this.userForm) {
        this.userForm.validate((valid) => {
          if (valid) {
            let params = this.form
            if (this.layer.row) {
              this.updateForm(params)
            } else {
              this.addForm(params)
            }
          } else {
            return false;
          }
        });
      }
    },
    // 新增提交事件
    addForm(params) {
      addUser(params)
      .then(res => {
        this.$message({
          type: 'success',
          message: '新增成功'
        })
        this.$emit('getTableData', true)
        this.layerDom && this.layerDom.close()
      })
    },
    // 编辑提交事件
    updateForm(params) {
      updateUser(params)
      .then(res => {
        this.$message({
          type: 'success',
          message: '编辑成功'
        })
        this.$emit('getTableData', false)
        this.layerDom && this.layerDom.close()
      })
    }
  }
})
</script>

<style lang="scss" scoped>
  
</style>