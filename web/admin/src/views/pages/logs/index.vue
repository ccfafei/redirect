<template>
  <div class="layout-container">
    <div class="layout-container-form flex space-between">
      <div class="layout-container-form-handle">
    
        <el-popconfirm title="批量删除" @confirm="handleDel(chooseData)">
          <template #reference>
            <el-button type="danger" icon="el-icon-delete" :disabled="chooseData.length === 0">批量删除</el-button>
          </template>
        </el-popconfirm>
      </div>
      <div class="layout-container-form-search">
         <el-date-picker
              v-model="date"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期" 
             
              :default-time="defaultTime"
              @change="searchTime"   
              style="width:400px"      
            />
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
        :data="userLogsData"
        @getTableData="getTableData"
        @selection-change="handleSelectionChange"
      >
      
        <el-table-column prop="from_domain" label="来源域名" align="center" />
        <el-table-column prop="to_domain" label="目标域名" align="center" />
        <el-table-column prop="ip" label="IP地址" align="center" />
        <el-table-column prop="user_agent" label="UA标识" align="center" width="400" />
        <el-table-column prop="access_time" label="访问时间" align="center">
            <template #default="scope">
              {{getTime(scope.row.access_time)}}
            </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="100">
          <template #default="scope">
            <el-popconfirm title="删除" @confirm="handleDel([scope.row])">
              <template #reference>
                <el-button type="danger" icon="el-icon-delete" ></el-button><!---->
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </Table>
    
    </div>
  </div>
</template>

<script>
import moment from 'moment'
import { defineComponent, ref, reactive } from 'vue'
import Table from '@/components/table/index.vue'
import { getLogsList,delLogs } from '@/api/user'

import { ElMessage } from 'element-plus'

export default defineComponent({
  name: 'crudTable',
  components: {
    Table,
  },
  setup() {
   
    // 存储搜索用的数据
    const date=ref([])
    const query = reactive({
      search: '',
      start_time: '',
      end_time:'',
    })
 


  const getTime = (time) => { 
      return moment(time).format("YYYY-MM-DD HH:mm:ss")
    }
  const defaultTime = ref([
    new Date(2000, 1, 1, 0, 0, 0),
    new Date(2000, 2, 1, 23, 59, 59),
  ])
    const searchTime=(dates)=> { 
      date.value = dates
      query.start_time = getTime(date.value[0])
      query.end_time = getTime(date.value[1])
	  if (dates === null || dates.length === 0) {
	    date.value = null
	  }

    }
    // 分页参数, 供table使用
    const page = reactive({
      index: 1,
      size: 20,
      total: 0
    })
    const loading = ref(true)
    const userLogsData = ref([])
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
      getLogsList(params)
      .then(res => {
        userLogsData.value = res.result.data
        page.total = Number(res.result.total)
      })
      .catch(error => {
        userLogsData.value = []
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
          return e.id
        }).join(',')
      }
      delLogs(params)
      .then(res => {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        getTableData(userLogsData.value.length === 1 ? true : false)
      })
      .catch(error => {
        console.log(error)
      })
    }
    getTableData(true)
    return {
      query,
      userLogsData,
      chooseData,
      loading,
      page,
      handleSelectionChange,
      handleDel,
      getTableData,
      searchTime,
      getTime,
      defaultTime,
      date
    }
  }
})
</script>

<style lang="scss" scoped>
  
</style>