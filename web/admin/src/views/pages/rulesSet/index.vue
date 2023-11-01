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
        :data="rulesData"
        @getTableData="getTableData"
        @selection-change="handleSelectionChange"
      >
        <el-table-column prop="app_name" label="任务名称" align="center" min-width="100" />
        <el-table-column prop="from_domain" label="来源域名" align="center"  min-width="100"  class-name="myCell" >
           <template #default="scope"> 
              <p v-for="(item,index) in scope.row.from_domain.slice(0, 5)">{{item}}</p>
            <el-button type="text" v-if="scope.row.from_domain.length>5" @click="moreFromDomain(scope.row)"  icon="el-icon-more"></el-button>
           </template>
        </el-table-column>
        <el-table-column prop="rule_data" label="目标url" align="center"  min-width="160"  class-name="myCell">
          <template #default="scope" >
            <p v-for="(item2,index) in scope.row.rule_data" class="myCellP">
             {{item2.to_domain}}
              <i class="vlineIcon"></i> {{item2.weight}}
            </p> 
          </template>

        </el-table-column>
         <el-table-column prop="ip_blacks" label="IP地址" align="left"  min-width="80"  class-name="myCell">
          <template #default="scope" >
            <div v-if="scope.row.ip_blacks">
              <p v-for="(item2,index) in scope.row.ip_blacks.slice(0, 5)" >
            {{item2}}
              </p> 
              
            <el-button type="text" v-if="scope.row.ip_blacks.length>5" @click="moreIp(scope.row)" icon="el-icon-more"></el-button>
            </div>
            <div v-else></div>
          </template>
         

        </el-table-column>

        <el-table-column prop="default_url" label="默认目标URL" align="center" min-width="100" />
     
        <el-table-column prop="status" label="状态" align="center"  min-width="60"  >
           <template #default="scope" >
          <span v-if="scope.row.status==1">启用</span>
          <span v-else>禁用</span>
           </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" align="center" />
        <el-table-column prop="updated_at" :formatter="filterTime" label="时间" align="center"  min-width="110" />
        <el-table-column label="操作" align="center" min-width="150" >
          <template #default="scope">
            <el-button @click="handleShare(scope.row)" icon="el-icon-share"></el-button>
            <el-button @click="handleEdit(scope.row)" icon="el-icon-edit"></el-button>
            <el-popconfirm title="删除" @confirm="handleDel([scope.row])">
              <template #reference>
                <el-button type="danger" icon="el-icon-delete"></el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </Table>
      <Layer :layer="layer" @getTableData="getTableData" v-if="layer.show"  :layerSign="layerSign"/>
     
     
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, reactive, inject, watch,provide  } from 'vue'
import Table from '@/components/table/index.vue'
import { getRulesList,delRule,urlPassword} from '@/api/rules'
import Layer from './layer.vue'
import { ElMessage } from 'element-plus'
export default defineComponent({
  components: {
    Table,
    Layer
  },
  setup() {
    const layerSign=ref(1)
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
    let shareForm=ref({
    share_url:'',
    password:''
   })
      
    // 分页参数, 供table使用
    const page = reactive({
      index: 1,
      size: 20,
      total: 0
    })
   
    const loading = ref(true)
    const rulesData = ref([])
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
      getRulesList(params)
      .then(res => {
        rulesData.value = res.result.data
        page.total = Number(res.result.total)
      })
      .catch(error => {
        rulesData.value = []
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
      delRule(params)
      .then(res => {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        getTableData(rulesData.value.length === 1 ? true : false)
      })
    }
    // 新增弹窗功能
    const handleAdd = () => {
      layer.title = '新增数据'
      layer.show = true
      layerSign.value=3
      delete layer.row
    }
    // 编辑弹窗功能
    const handleEdit = (row) => {
      layer.title = '编辑'
      layer.row = row
      layer.show = true
      layerSign.value=1
      layer.showButton=true
  
    }
    const handleShare=(row)=>{
        layer.title = '分享'
        layer.row = row
        layer.show = true
        layerSign.value=4
        layer.showButton=false
        getUrlPassword(row.id)
    }
  // 获取分享地址密码
  const getUrlPassword = (id) => {
      let params = {
       rule_id:id
      }
      urlPassword(params)
      .then(res => {
        shareForm.value.share_url=res.result.share_url
        shareForm.value.password=res.result.password

      })
    }

 provide("shareData", shareForm)
   const filterTime=(row, column, cellValue)=>{
      return cellValue.slice(0,10)
    }
  //查看更多
  const moreFromDomain=(row)=>{
     layer.title = '来源域名'
     layer.row = row
     layer.show = true
    layerSign.value=0
  }
   const moreIp=(row)=>{
     layer.title = 'Ip地址'
     layer.row = row
     layer.show = true
    layerSign.value=2
  }
  
 getTableData(true)
      
    return {
      query,
      rulesData,
      chooseData,
      loading,
      page,
      layer,
      handleSelectionChange,
      handleAdd,
      handleEdit,
      handleDel,
      getTableData,
      filterTime,
      moreFromDomain,
      moreIp,
      layerSign,
      handleShare,
shareForm,


    }
  }
})
</script>

<style lang="scss" scoped>
  
// .myCell:not(.table-header),.myCell:not(.table-header) .cell{
//   padding: unset!important;
// }
.el-table .cell{padding:0 !important;}
.myCell .myCellP{
  border-bottom: var(--el-table-border);
  padding: 4px 0;
}
.myCell .cell .myCellP:last-child{
  border-bottom: unset;
}

</style>