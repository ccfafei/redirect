<template>
  <div class="layout-container">
    <div class="layout-container-form flex space-between">
      <div class="layout-container-form-handle">
        <el-button type="primary" icon="el-icon-circle-plus-outline" @click="handleAdd">新增</el-button>
        <el-popconfirm title="批量删除" @confirm="handleDel(chooseData)">
          <template #reference>
            <el-button type="danger" icon="el-icon-delete" :disabled="chooseData.length === 0">批量删除</el-button>
          </template>
        </el-popconfirm>
      </div>
      <div class="layout-container-form-search">
        <el-input v-model="query.search" placeholder="请输入关键词进行检索" @change="getTableData(true)"></el-input>
        <el-button type="primary" icon="el-icon-search" class="search-btn" @click="getTableData(true)">搜索</el-button>
      </div>
    </div>
    <div class="layout-container-table">
      <Table
        ref="table"
        v-model:page="page"
        v-loading="loading"
        :showIndex="true"
        :showSelection="true"
        :data="userListData"
        @getTableData="getTableData"
        @selection-change="handleSelectionChange"
      >
      <el-table-column prop="account" label="账号" align="center" />
        <el-table-column prop="name" label="名称" align="center" />
       
        <el-table-column label="操作" align="center"  width="120">
          <template #default="scope">
            <el-button @click="handleEdit(scope.row)" icon="el-icon-edit"></el-button>
            <el-popconfirm title="删除" @confirm="handleDel([scope.row])">
              <template #reference>
                <el-button type="danger" :disabled="adminHas==1?true:false" icon="el-icon-delete"></el-button><!---->
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </Table>
      <Layer :layer="layer" @getTableData="getTableData" v-if="layer.show" />
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, reactive } from 'vue'
import Table from '@/components/table/index.vue'
import { getUserList ,delUser } from '@/api/user'
import Layer from './layer.vue'
import { ElMessage } from 'element-plus'

export default defineComponent({
  name: 'crudTable',
  components: {
    Table,
    Layer
  },
  setup() {
   const adminHas = ref(null) 
    // 存储搜索用的数据
    const query = reactive({
      search: ''
    })
    // 弹窗控制器
    const layer = reactive({
      show: false,
      title: '新增',
      showButton: true
    })
    // 分页参数, 供table使用
    const page = reactive({
      index: 1,
      size: 20,
      total: 0
    })
    const loading = ref(true)
    const userListData = ref([])
    const chooseData = ref([])
    const handleSelectionChange = (val) => {
      chooseData.value = val
    }
    // 获取表格数据
    // params <init> Boolean ，默认为false，用于判断是否需要初始化分页
    const getTableData = (init) => {
      loading.value = true
      if (init) {
        page.index = 1
      }
      let params = {
        page: page.index,
        size: page.size,
        ...query
      }
      getUserList(params)
      .then(res => {
        userListData.value = res.result.data
        page.total = Number(res.result.total)
      })
      .catch(error => {
        userListData.value = []
        page.index = 1
        page.total = 0
      })
      .finally(() => {
        loading.value = false
      })
    }
    // 删除功能
    const handleDel = (data) => {
      
      let params = {
        ids: data.map((e)=> {
          if(e.account=='admin'){
              adminHas.value=1
          }
          return e.id
        }).join(',')
      }
       if(adminHas.value==1){// 判断管理员不能删除
             ElMessage({
          type: 'warning',
          message: '管理员不能删除'
        })
            return false
          }
      delUser(params)
      .then(res => {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        getTableData(userListData.value.length === 1 ? true : false)
      })
      .catch(error => {
        console.log(error)
      })
    }
    // 新增弹窗功能
    const handleAdd = () => {
      layer.title = '新增'
      layer.show = true
      delete layer.row
    }
    // 编辑弹窗功能
    const handleEdit = (row) => {
      layer.title = '编辑'
      layer.row = row
      layer.show = true
    }
    getTableData(true)
    return {
      query,
      userListData,
      chooseData,
      loading,
      page,
      layer,
      handleSelectionChange,
      handleAdd,
      handleEdit,
      handleDel,
      getTableData,
      adminHas
    }
  }
})
</script>

<style lang="scss" scoped>
  
</style>