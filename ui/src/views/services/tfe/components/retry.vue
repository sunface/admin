<template>
  <div class="retry">
       <el-form  label-position="left" label-width="70px"  size="mini">
          <el-form-item label="策略名称" style="margin-top:10px">
              {{strategy.name}}
          </el-form-item>
          <el-form-item label="超时时间" style="margin-top:10px">
              {{strategy.content.req_timeout}}
          </el-form-item>
          <el-form-item label="重试次数" style="margin-top:10px">
              {{strategy.content.retry_times}}
          </el-form-item>
          
          <el-form-item label="重试间隔" style="margin-top:10px">
              {{strategy.content.retry_interval}}
          </el-form-item>
        </el-form>
  </div>
</template>

<script>
/* eslint-disable */
export default {
  name: 'retry',
  props: {
      strategy : {
        type: Object,
        default: {}
      }
  },
  components: {},
  data() {
      return {
      }
  },
  methods: {
  },
  created() {

  },
  mounted() {
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>

