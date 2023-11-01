<template>
  <div class="logo-container">
    <!-- <img src="@/assets/logo.png" alt=""> -->
    <!-- <h1>{{ systemTitle }}</h1> -->
    <h3>{{titleText}}</h3>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { useStore } from 'vuex'
import { systemTitle } from '@/config'
import { getTitle} from '@/api/index'
export default defineComponent({
  setup() {
    const store = useStore()
     const titleText=ref('')
       const getTitleText = () => {
      let params={
        rule_id:store.state.user.ruleId,
      }
        getTitle(params).then(res=>{
          titleText.value=res.result.app_name
       
        }).catch(error => {
       
        })
   }
   getTitleText()
    // const isCollapse = computed(() => store.state.app.isCollapse)
    return {
      // isCollapse,
      systemTitle,
      titleText
    }
  }
})
</script>

<style lang="scss" scoped>
  .logo-container {
    min-height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    background-color: var(--system-logo-background);
    h3 {
      font-size: 16px;
      white-space: nowrap;
      color: var(--system-logo-color);
    }
  }
</style>