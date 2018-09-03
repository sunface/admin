<template>
  <div class="traffic">
       <el-form  label-position="left" label-width="70px"  size="mini">
          <el-form-item label="Name" style="margin-top:10px">
              {{strategy.name}}
          </el-form-item>
          <div class="form-block" >
            <span>Api Traffic</span>
            <el-form-item label="QPS" style="margin-top:10px">
                {{strategy.content.qps}}
            </el-form-item>
            <el-form-item label="Concurrent" style="margin-top:10px">
                {{strategy.content.concurrent}}
            </el-form-item>
          </div>
          <div class="form-block" >
              <span>User Quota</span>
                <el-form-item label="Param" style="margin-top:10px">
                    {{strategy.content.param}}
                </el-form-item>
                <el-form-item label="Time Span" style="margin-top:10px">
                    {{strategy.content.span}}
                </el-form-item>
                <el-form-item label="Access Times" style="margin-top:10px">
                    {{strategy.content.times}}
                </el-form-item>
          </div>
        </el-form>
  </div>
</template>

<script>
/* eslint-disable */
export default {
  name: 'traffic',
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

