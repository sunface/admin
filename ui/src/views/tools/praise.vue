<template>
  <div class="app-container">
        <el-row :gutter="20" style="margin-top:50px;">
            <el-col :xs="{span:24}" :sm="{span:18}" :md="{span: 8}" :lg="{ span: 6}">
                <el-card class="box-card">
                <div slot="header" class="clearfix">
                    <span>输入被惩罚者姓名</span>
                </div>
                <div style="height:100px;">
                    <el-form>
                    <el-form-item prop="title">
                        <md-input icon="search" name="title" placeholder="姓名之间以空格分割" v-model="punishee"></md-input>
                    </el-form-item>
                    </el-form>
                </div>

                <el-button v-waves type="primary" @click="submit">提交</el-button>
                </el-card>
            </el-col>
            <el-col :xs="{span:24}" :sm="{span:18}" :md="{span: 14}" :lg="{ span: 12}" v-show="rewarder.length != 0">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>奖励结果</span>
                    </div>
                    <div>{{rewarder}}</div>
                </el-card>
            </el-col>
        </el-row>
  </div> 
</template>

<script>
/* eslint-disable */
import MdInput from '@/components/MDinput'
import waves from '@/directive/waves/index.js' // 水波纹指令
import request from '@/utils/request' 

export default {
  name: 'praise',
   directives: {
    waves
  },
   components: {
    MdInput
  },
  data() {
      return {
          punishee: '',
          rewarder: []
      }
  },
  methods: {
      submit() {
          if (this.punishee == '') {
               this.$message({
                    message: '用户不能为空',
                    type: 'warning',
                    duration: 3 * 1000,
                    center: true
                })
                return 
          }
        var _this = this
        // 获取用户信息
        request({
            url: '/tools/praise',
            method: 'POST',
            params: {
                punishee: this.punishee,
            }
        }).then(response => {
            if (response.data == '输入不合法') {
                _this.$message({
                    message: '用户不存在',
                    type: 'error',
                    duration: 3 * 1000,
                    center: true
                })
                _this.rewarder = []
                return 
            }
            _this.$message({
                    message: '抽奖成功',
                    type: 'success',
                    duration: 3 * 1000,
                    center: true
                })
            _this.rewarder = response.data
        }).catch(error => {
           
        })
      }
  },
  created() {
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>


<style scoped>
.app-container {
  background-color: #f0f2f5;
  min-height: calc(100vh - 84px);
}
</style>
