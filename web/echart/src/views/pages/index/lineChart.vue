<template>
  <h3 style="padding-left:20px;font-size:16px;padding-bottom:20px;">统计分析</h3>
   <el-radio-group v-model="timeCurrent"  size="small" @change="changeTime" class="timebar"  >
      <el-radio-button label="today" >今日</el-radio-button>
      <el-radio-button label="yesterday" >昨日</el-radio-button>
      <el-radio-button label="week" >最近七日</el-radio-button>
      <el-radio-button label="month" >最近30天</el-radio-button>
  </el-radio-group>
  <div id="lineChart"  />
</template>

<script >
import * as echarts from 'echarts'
import moment from 'moment'
import { useStore } from 'vuex'
import { defineComponent, reactive ,onMounted ,toRefs,ref,nextTick  } from 'vue'
import { useEventListener } from '@vueuse/core' //引入监听函数，监听在vue实例中可自动销毁，无须手动销毁
import { dataEchartApi} from '@/api/index'
export default defineComponent({
  setup() {
    const xAxisData=ref(null)// x轴
    const ipData=ref(null)//ip
    const pvData=ref(null)//pv
    const uvData=ref(null)//uv
    const timeCurrent = ref('today')
    

 const store = useStore()
    const echartsData=ref(null)
    const arrCon=(data)=>{
      let arr=[]
      for (let j = 0; j < data.length; j++) {
           arr.push(data[j].num)
        }
         return arr
    }
     const arrTime=(data)=>{
      let arr=[]
      for (let j = 0; j < data.length; j++) {
         arr.push(data[j].ts)
        }
         return arr
      
    }
   const getDataEchart = () => {
      let params={
        rule_id:store.state.user.ruleId,
        date_type:timeCurrent.value
      }
        dataEchartApi(params).then(res=>{
          echartsData.value=res.result   
          const timeArr=echartsData.value.ip_data
          const ipD=echartsData.value.ip_data
          const pvD=echartsData.value.pv_data
          const uvD=echartsData.value.uv_data
           ipData.value=arrCon(ipD)
           pvData.value=arrCon(pvD)
            uvData.value=arrCon(uvD)
            xAxisData.value=arrTime(timeArr)
      
          initeCharts()
        }).catch(error => {
        console.log(error)
        })
   }
    onMounted(() => {
        initeCharts()
      })
 
    const dateTime=ref({
      startTime:'',
      endTime:''
    })
 const timeF=(params)=>{
   let str
   if (timeCurrent.value == 'today'||timeCurrent.value == 'yesterday') {
      str = params.substr(10,6)// 
    }
    if ( timeCurrent.value =='week' || timeCurrent.value == 'month') {
      str = params.substr(5, 5)
    }
    return str
 }
     
     const initeCharts = () => {
        let myChart = echarts.init(document.getElementById('lineChart'))   
        const  option={
            title: {
              text: ''
            },
            tooltip: {
              trigger: 'axis',
              axisPointer: {
                type: 'none', //去除移入垂直指示线
              },
                formatter : function(data){ 
                  let str1,str2,str3,time
                     for(let j=0;j<data.length;j++){
                        str1= '<span style=display:inline-block;width:110px;margin:2px 0;>'+data[0].marker+data[0].seriesName+'</span>'+'<strong>'+data[0].value+'</strong>'+' <br/>'
                        str2= '<span style=display:inline-block;width:110px;margin:2px 0;>'+data[1].marker+data[1].seriesName+'</span>'+'<strong>'+data[1].value+'</strong>'+' <br/>'
                        str3= '<span style=display:inline-block;width:110px;margin:2px 0;>'+data[2].marker+data[2].seriesName+'</span>'+'<strong>'+data[2].value+'</strong>'+' <br/>'
                        time= timeF(data[0].name)
                      return '<i class=el-icon-time style=color:#000;margin-bottom:5px;></i> ' + time+' <br/>'+str1+str2+str3   
                   } 
                }
            },
            legend: {
               right: '2%',
                orient:"horizontal",
              
               formatter: function (name){
               
                 return name
               }

            },
          
            grid: {
              left: '2%',
              right: '2%',
              bottom: '2%',
              containLabel: true
            },
            xAxis: [
              {
                type: 'category',
                boundaryGap: false,
                data: xAxisData.value,
                axisLine: {
                show: true,
                lineStyle: {
                  width: 0 // x轴横线
                }
               },
               axisLabel : {
                        interval : 0,
                         rotate:"30",
                        formatter : function(params){ 
                           return timeF(params)
                        }

                    }
                
              }
              
            ],
            yAxis: [
              {
                
                type: 'value',
                 scale: true, // y轴的起始数据不是0，而是后台数据最小值
                  axisLabel: {
                  formatter: '{value}'
                }
              }
            ],
            series: [
              {
                name: 'IP数',
                type: 'line',
                stack: 'Total',
             
                emphasis: {
                  focus: 'series',
                   disabled: false // 鼠标移动区域高亮 去除
                },
                lineStyle: {
                  width: 1,
                  color: '#6cffd9'
                },
                  itemStyle: { // 曲线和点颜色
                    color: '#6cffd9'
                 },
                showSymbol: false, // 鼠标经过才显示线上圆点
                 areaStyle: {
                  color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 0,
                    y2: 1,
                    colorStops: [{
                      offset: 0,
                      color: '#6cffd9' ,// 0% 处的颜色
                      option: 1
                    }, {
                      offset: 1,
                      color: '#ffffff' // 100% 处的颜色
                    }],
                    global: false // 缺省为 false
                  }
                },
                data: ipData.value
              },
             {
                name: '浏览量PV',
                type: 'line',
                stack: 'Total',
               
                emphasis: {
                  focus: 'series',
                   disabled: false // 鼠标移动区域高亮 去除
                },
                 lineStyle: {
                  width: 1,
                  color: '#ffd96c'
                },
                itemStyle: { // 曲线和点颜色
                    color: '#ffd96c'
                 },
                 showSymbol: false, // 鼠标经过才显示线上圆点
                 areaStyle: {
                  color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 0,
                    y2: 1,
                    colorStops: [{
                      offset: 0,
                      color: '#ffd96c' ,// 0% 处的颜色
                      option: 1
                    }, {
                      offset: 1,
                      color: '#ffffff' // 100% 处的颜色
                    }],
                    global: false // 缺省为 false
                  }
                },
                data: pvData.value
              },
               {
                name: '访客数UV',
                type: 'line',
                stack: 'Total',
         
                emphasis: {
                  focus: 'series',
                   disabled: false // 鼠标移动区域高亮 去除
                },
                 lineStyle: {
                  width: 1,
                  color: '#6cbdff'
                },
                itemStyle: { // 曲线和点颜色
                    color: '#6cbdff'
                 },
                 showSymbol: false, // 鼠标经过才显示线上圆点
                 areaStyle: {
                  color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 0,
                    y2: 1,
                    colorStops: [{
                      offset: 0,
                      color: '#6cbdff' ,// 0% 处的颜色
                      option: 1
                    }, {
                      offset: 1,
                      color: '#ffffff' // 100% 处的颜色
                    }],
                    global: false // 缺省为 false
                  }
                },
                data: uvData.value
              },
              
             
            ]
          }
        myChart.setOption(option)
        useEventListener("resize", () => myChart.resize())
      }
      
    
    
   
   function changeTime(label){
  
      if (label == "today") {// 今日
          dateTime.value.startTime = moment().startOf('day').format('YYYY/MM/DD HH:mm:ss')
          dateTime.value.endTime = moment().format('YYYY/MM/DD HH:mm:ss')//当前时间
      }else if(label =="yesterday"){//昨日
          dateTime.value.startTime = moment().day(moment().day() - 1).startOf('day').format('YYYY/MM/DD HH:mm:ss')
          dateTime.value.endTime = moment().day(moment().day() - 1).endOf('day').format('YYYY/MM/DD HH:mm:ss')
      }else if(label == "week"){// 最近7天
          dateTime.value.startTime = moment().subtract(7,"days").startOf('day').format('YYYY/MM/DD HH:mm:ss')
          dateTime.value.endTime = moment().format('YYYY/MM/DD HH:mm:ss')//当前时间
      }else if(label == "month"){// 最近30天
          dateTime.value.startTime = moment().subtract(30,"days").startOf('day').format('YYYY/MM/DD HH:mm:ss')
          dateTime.value.endTime = moment().format('YYYY/MM/DD HH:mm:ss')//当前时间
      }
       getDataEchart(dateTime.value)
   }
     changeTime(0)
    return {
   
      timeCurrent,
      changeTime
    }
  }
})
</script>

<style lang="scss" scoped>
  .timebar{margin-left:20px; position: relative;top:15px; z-index: 999;}
  #lineChart{height:400px;}
@media screen and ( max-width: 750px ) {
   .timebar{ top:-15px; }
    #lineChart{height:300px;}
}
</style>