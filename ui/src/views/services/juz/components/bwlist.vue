<template>
  <div class="bwlist">
       <el-form  label-position="left" label-width="70px"  size="mini">
          <el-form-item label="Name" style="margin-top:10px">
              {{strategy.name}}
          </el-form-item>
          <el-form-item label="Type" style="margin-top:10px">
              <span v-if="strategy.sub_type==1">Black</span>
              <span v-else>White</span>
          </el-form-item>
          
          <el-form-item label="List" style="margin-top:10px">
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

