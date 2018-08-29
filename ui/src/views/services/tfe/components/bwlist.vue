<template>
  <div class="bwlist">
       <el-form  label-position="left" label-width="70px"  size="mini">
          <el-form-item label="名称" style="margin-top:10px">
              {{strategy.name}}
          </el-form-item>
          <el-form-item label="类型" style="margin-top:10px">
              <span v-if="strategy.sub_type==1">黑名单</span>
              <span v-else>白名单</span>
          </el-form-item>
          
          <el-form-item label="名单" style="margin-top:10px">
              <div v-for="c in strategy.content" :key="c.val">
                  <el-tag type="success" size="small">{{c.key}}</el-tag>
                  <el-tag type="warning" size="small">{{c.val}}</el-tag>
              </div>
          </el-form-item>
        </el-form>
  </div>
</template>

<script>
/* eslint-disable */
export default {
  name: 'bwlist',
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

