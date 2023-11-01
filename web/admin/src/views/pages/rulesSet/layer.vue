<template>
  <Layer :layer="layer" @confirm="submit" ref="layerDom">
    <el-form :model="form" :rules="rules" ref="ruleForm" label-width="120px" style="margin-right:30px;"  v-if="layerSign==1||layerSign==3">
      <el-form-item label="任务名称：" prop="app_name">
        <el-input v-model="form.app_name" placeholder="请输入名称" ></el-input>
      </el-form-item>
      <el-form-item label="来源域名：" prop="from_domain">
       <el-input v-model="form.from_domain"  :autosize="{ minRows: 6, maxRows: 10}" type="textarea" placeholder="请输入来源域名"/>
      <span class="tipText">请用 "回车" 分隔</span>
      </el-form-item>
      <el-form-item label="IP地址：" prop="ip_blacks">
       <el-input v-model="form.ip_blacks"  :autosize="{ minRows:6, maxRows: 10 }" type="textarea" placeholder="请输入ip范围"/>
      <span class="tipText">请用 "回车" 分隔</span>
      </el-form-item>
      <el-form-item label="目标URL：" prop="rule_data">

            <el-row v-for="(item, index) in form.rule_data" :key="index" >
              <el-col :span="12">
                <el-form-item :prop="'rule_data.' + index + '.to_domain'" >
                  <el-input v-model="item.to_domain" placeholder="http://或https://"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3" align=right>权重：</el-col>
              <el-col :span="4" >
                
                <el-form-item  :prop="'rule_data.' + index + '.weight'" >
                  <el-input v-model.number="item.weight" placeholder="权重"></el-input>
                </el-form-item>
              </el-col>
            
              <el-col :span="5" class="iconBtn">
                <div v-if="index+1== form.rule_data.length" @click="addIpForm">
                  <i class="el-icon-plus"></i> 
                </div>
                <div v-if="index != 0" @click="delIpForm(item,index)">
                  <i class="el-icon-minus"></i>
                </div>
              </el-col>
          </el-row>
      </el-form-item>
      <el-form-item label="默认URL:" prop="default_url">
        <el-input v-model="form.default_url" placeholder="请输入"></el-input>
      </el-form-item>
       <el-form-item label="状态：" prop="status">
       <el-switch v-model="form.status" :active-value="1" :inactive-value="0" @change="statusN(form.status)" />
      </el-form-item>
        <el-form-item label="备注：" prop="remark">
       <el-input v-model="form.remark" placeholder="备注"></el-input>
      </el-form-item>
    </el-form>
    <!--更多-->
    <el-form :model="form" :rules="rules" ref="ruleForm" label-width="120px" style="margin-right:30px;"  v-if="layerSign==0">
      <el-form-item label="来源域名：" prop="from_domain">
        <el-input v-model="form.from_domain"  :autosize="{ minRows: 6, maxRows: 30}" type="textarea" placeholder="请输入来源域名"/>
        <span class="tipText">请用 "回车" 分隔</span>
      </el-form-item>
      
    </el-form>
     <el-form :model="form" :rules="rules" ref="ruleForm" label-width="120px" style="margin-right:30px;"  v-if="layerSign==2">
       <el-form-item label="IP地址：" prop="ip_blacks">
        <el-input v-model="form.ip_blacks"  :autosize="{ minRows: 6, maxRows: 30}" type="textarea" placeholder="请输入来源域名"/>
        <span class="tipText">请用 "回车" 分隔</span>
      </el-form-item>
    </el-form>
     <el-form :model="shareData" :rules="rules" ref="ruleForm" label-width="120px" style="margin-right:30px;"  v-if="layerSign==4">
        <el-form-item label="URL:" prop="share_url">
          <el-input v-model="shareData.share_url"  :disabled=true></el-input>
        </el-form-item>
        <el-form-item label="密码:" prop="password">
          <el-input v-model="shareData.password" placeholder="请输入"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="copyDomText(shareData)">复制</el-button>
          <el-button @click="close">取消</el-button>
       </el-form-item>
    
     </el-form>
  </Layer>

</template>

<script>
import { defineComponent, ref,inject  } from 'vue'
import { ElMessage } from 'element-plus'
import Layer from '@/components/layer/index.vue'
import { addRule, updateRule ,updatePassword} from '@/api/rules'
export default defineComponent({
  components: {
    Layer,
  },
  props: {
   layerSign:1,
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
   let shareData = inject("shareData");

    const ruleForm = ref(null)
    const layerDom=ref(null)

    let form = ref({
      id:'',
      app_name: '',
      from_domain:null,
      rule_data:[{
        to_domain:'',
        weight:100
      }],
      status:1,
      remark:'',
      ip_blacks:null,
      default_url:''
    })
   
    const rules = {
      app_name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
      from_domain: [{ required: true, message: '请输入来源域名', trigger: 'blur' }]

    }
    
    init()
    function init() { // 用于判断新增还是编辑功能
      if (props.layer.row){// 编辑，查看更多
         form.value = JSON.parse(JSON.stringify(props.layer.row)) // 数量量少的直接使用这个转  
            form.value.from_domain=form.value.from_domain.join("\n")  //数组转字符串
            if(form.value.ip_blacks){
              form.value.ip_blacks=form.value.ip_blacks.join("\n")
            }
          if(form.value.rule_data==null){
            form.value.rule_data=[{
              to_domain:'',
              weight:100
            }]
          }
  
      } 
      else {// 新增
        
      }
      
    }
      function addIpForm(){
       form.value.rule_data.push({
           to_domain:'',
           weight:100,
      });

      }
      function delIpForm(item, index){
         form.value.rule_data.splice(index, 1);
      }

    const statusN = (status) => {
       form.value.status  = status
    
    }


const copyDomText=(val)=> {
  
  // 获取需要复制的元素以及元素内的文本内容
  let text='地址：'+val.share_url+'，'+'密码：'+val.password
  // 添加一个input元素放置需要的文本内容
  const input = document.createElement("input");
  input.value = text;
  document.body.appendChild(input);
  // 选中并复制文本到剪切板
  input.select();
  document.execCommand("copy");
  // 移除input元素
  document.body.removeChild(input);
  let data={
    rule_id:props.layer.row.id,
    password:val.password
  }
  updatePassword(data).then(res=>{
      ElMessage({
    message: "复制成功",
    type: "success",
  });
  props.layer.show=false

  })

}
 const close=()=> {
      props.layer.show=false
    }
    return {
      form,
      rules,
      layerDom,
      ruleForm,
      addIpForm,
      delIpForm,
      statusN,
      shareData,
      close,
      copyDomText
 
  
 
    }
    
  },
  methods: {
  // 字符串转数组
    formatAddress(address) {
      let temp=[]
      if(address==null){
        temp=null
      }else{
        temp = address.split(/[\n]/g)
        for (var i = 0; i < temp.length; i++) {
            if (temp[i] == '') {
              temp.splice(i, 1) // 删除数组索引位置应保持不变
              i--
            }
          }   
           temp=[...new Set(temp)]
      }

      return temp
    },
    
  
   // 提交
    submit() {
      if (this.ruleForm) {
        this.ruleForm.validate((valid) => {
          if (valid) {
               let params={
                  id:this.form.id,
                  app_name:this.form.app_name,
                  from_domain:this.formatAddress(this.form.from_domain),
                  rule_data:this.form.rule_data,
                  status:this.form.status,
                  remark:this.form.remark,
                  ip_blacks:this.formatAddress(this.form.ip_blacks),
                  default_url:this.form.default_url,
              }  
            // let params = this.form
            if (this.layer.row) {
              this.updateForm(params)
            }else {
              this.addForm(params)
            }
          } else {

            return false
          }
        })
      }
    },
    // 新增提交事件
    addForm(params) {
      addRule(params)
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
      updateRule(params)
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
  .tipText{font-size:12px;color:#999;position: absolute;right:10px;top:-23px;}

  .iconBtn{
      padding-left:10px;
      i{cursor: pointer;display: inline-block;}
      div{
        display: inline-block !important;width:30px;
      }
    }
    // .moreFromDomainList{
    //   margin-bottom:30px;
    //     li{display: inline-block;width: 20%;line-height:30px;}
    // } 
</style>