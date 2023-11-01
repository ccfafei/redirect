<template>
  
  <el-row class="p10" :gutter="10">
    
    <el-col  :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
      <el-row >
        <el-col :xs="24" :sm="8" :md="8" :lg="8" :xl="8">
          <div class="grid-content ep-bg-purple-dark statisticsTotal" >
            <dl>
              <dt class="ipborder"><strong>IP</strong>数量</dt>
              <dd>
              <p><span class="blue">{{totalData.today_ip_num}}</span>/ 今日</p>
              <p><span class="black">{{totalData.yesterday_ip_num}}</span>/ 昨日</p>
              </dd>
            </dl>
          </div> 
        </el-col> 
        <el-col :xs="24" :sm="8" :md="8" :lg="8" :xl="8">
          <div class="grid-content ep-bg-purple-dark statisticsTotal" >
              <dl>
              <dt class="pvborder"><strong>PV</strong>浏览量</dt>
              <dd>
              <p><span class="blue">{{totalData.today_pv_num}}</span>/ 今日</p>
              <p><span class="black">{{totalData.yesterday_pv_num}}</span>/ 昨日</p>
              </dd>
            </dl>
            
          </div> 
        </el-col>
        <el-col :xs="24" :sm="8" :md="8" :lg="8" :xl="8">
          <div class="grid-content ep-bg-purple-dark statisticsTotal" >
              <dl>
              <dt class="uvborder"><strong>UV</strong>访客数</dt>
              <dd>
              <p><span class="blue">{{totalData.today_uv_num}}</span>/ 今日</p>
              <p><span class="black">{{totalData.yesterday_uv_num}}</span>/ 昨日</p>
              </dd>
            </dl>
          
          </div> 
        </el-col>
      </el-row>
       <el-row ></el-row>
      <!---->
      <el-row>
         <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" >
           <div class="grid-content ep-bg-purple" ><lineChart /></div>
         </el-col>
      </el-row>
    </el-col>
    <el-col  :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
      <div class="grid-content ep-bg-purple-light visitorDomain" >
        <ul>
          <li>来源域名<span>访客数</span></li>
          <li v-for="item in domainData"> {{item.from_domain}}<span>{{item.num}}</span></li>
        </ul>
      </div>
    </el-col>
  </el-row>
 
</template>

<script>

import { defineComponent, ref, reactive, inject, watch  } from 'vue'
import lineChart from './lineChart.vue'
import { totalApi,dataDomainApi} from '@/api/index'
import { useStore } from 'vuex'
export default defineComponent({  
  components: {
 
   lineChart
  },
  setup() {
    const store = useStore()
   const totalData=ref({})
    const domainData=ref({})
 
   const getTotalData = () => {
    let params={
      rule_id:store.state.user.ruleId
    }
      totalApi(params).then(res=>{
        totalData.value=res.result
      }).catch(error => {
      
      })
   }
    getTotalData()

    const getDataDomain = () => {
      let params={
        rule_id:store.state.user.ruleId,
          top:17
      }
        dataDomainApi(params).then(res=>{
          domainData.value=res.result
        }).catch(error => {
       
        })
   }
   getDataDomain()

    return {
    totalData,
    getDataDomain,

    domainData,



    }
  }
})
</script>

<style lang="scss" scoped>
.blue{color:#409EFF;}
.black{color:#333;}
.pvborder{ border:2px solid #ffd96c; }
.ipborder{ border:2px solid #6cffd9;}
.uvborder{ border:2px solid #6cbdff;}
.p10{ padding:10px;}
.el-row {
  min-height: 10px;
  text-align: left;
   background: #efefef;
}
.el-col {

}
.grid-content {
  background: #fff;
  min-height: 100%;
  padding:20px;
}
.statisticsTotal{
  
  dl{margin:0;padding:10px 0; }

 p{padding:10px;color: #999;}
 strong{display: block;font-size:18px;color: #333;font-weight: 500;}
  dt{color: #999;border-radius: 50%;width:85px;height:60px;padding-top:25px;line-height:20px; text-align: center;}
  dd,dt,span{display: inline-block;margin-left:20px; vertical-align: middle;font-size:12px;}
  span{ font-weight:500;font-size:18px;margin-right:5px;}
}
.visitorDomain{

   li{
  position: relative;
  padding:10px ;
  &:first-child{
    font-weight: 700;
    padding:10px 0;
    // margin-bottom: 10px;;
    span{right:0;}
  }
  span{display:inline-block;position: absolute;right:10px;top:10px; text-align: right;}
 border-bottom: 1px solid #efefef;
}
}

 @media screen and ( max-width: 750px ) {

  .grid-content {

  padding:10px 10px 0 10px;
}
  .statisticsTotal dl{padding:0 0 5px 0;border-bottom:1px solid #efefef;}
   
 }
</style>