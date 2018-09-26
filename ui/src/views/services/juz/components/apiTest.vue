<template>
  <div class="api-test">
    <el-form class="black-form"  ref="dataForm" label-position="left" label-width="130px" >
        <el-form-item label="API ID"  >
            {{api.api_id}}
        </el-form-item>    
        <el-form-item :label="$t('juz.backendInfo')"  >
            <span v-if="api.backend_type==1">HTTP(S)-</span>
            <span v-else>Mock-</span>
            <span v-if="api.addr_type==1">URL-</span>
            <span v-else>ETCD-</span>
            {{api.backend_addr}}
            <span v-if="api.addr_type==2">-{{api.backend_uri}}</span>
        </el-form-item>
        <el-form-item :label="$t('juz.mode')">
            <el-radio-group v-model="debugOn" @change="handleDebug">
                <el-radio :label=true style="color:white">ON</el-radio>
                <el-radio :label=false style="color:white">OFF</el-radio>
            </el-radio-group> 
        </el-form-item>

        <!-- <el-form-item label="Juz-Api Addr">
            {{tfeAddr}}
        </el-form-item> -->


        <el-form-item :label="$t('common.param')" style="width:600px">
            <el-input type="textarea" :autosize="{ minRows: 10, maxRows: 12}" :value="makeParam(testParams)" @change="paramEnter" style="margin-top:10px">
            </el-input>
            <!-- <el-alert
                style="margin-top:8px"
                title="Form params only, e.g. http://test.com/service/api?api_name=test&api_version=1"
                :closable=false
                type="warning">
            </el-alert> -->
        </el-form-item>
        <el-form-item :label="$t('common.output')" style="width:800px">
            <el-input type="textarea" :autosize="{ minRows: 15, maxRows: 50}" :placeholder="$t('juz.testResult')" :value="calcResult()" :disabled=true>
            </el-input>
        </el-form-item>

        <el-form-item>
            <el-button type="success"  @click="startTest()">Start</el-button>
        </el-form-item>    
    </el-form>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request'
import {proxy} from '@/api/apiProxy'
export default {
  name: 'apiTest',
   props: {
    api: {
      type: Object,
      default: {}
    }
  },
  data() {
      return {
        tfeAddr: '',
        testResult: '',
        testPass : '',
        testParams: {}, 

        debugOn: false
      }
  },
  watch: {
    api() {
      this.initTest()
    }
  },
  methods: {
    calcResult() {
        if (this.testResult == '') {
            return ''
        }
        
        var e = JSON.stringify(JSON.parse(this.testResult),null,'\t')
        return e
    },
    handleDebug() {
        if (this.debugOn) {
            this.$set(this.testParams,'debug_on',true)
        } else {
            this.$set(this.testParams,'debug_on',undefined)
        }
    },
    paramEnter(e) {
        try {
            var t = JSON.parse(e)
            this.testParams = t
            if (this.testParams.debug_on) {
                this.debugOn = true
            } else {
                this.debugOn = false
            }
            // save params
            var params = {
                target_app: 'juzManage',
                target_path: '/manage/api/saveParam',
                api_id: this.api.api_id,
                params: e
             }
            proxy('POST',params).then(res => {
            })
        } catch (e) {
            this.$message({
                message: 'Params json invalid',
                type: 'error',
                duration: 5 * 1000,
                center: true
              })
        }
    },
    showName(apiID) {
        return apiID.substring(0,apiID.length-3)
    },
    showVersion(apiID) {
        return apiID[apiID.length-1]
    },
    startTest() {
      this.testResult = ''
      this.testPass = ''
      request({
            url: '/tools/testApi',
            method: 'POST',
            params: {
                tfe_addr: process.env.TFE_ADDR,
                params: this.testParams,
            }
        }).then(response => {
          if (response.data.status == 200) {
              this.$message({
                message: '测试成功',
                type: 'success',
                duration: 5 * 1000,
                center: true
              })
              this.testResult = response.data.data
              this.testPass = 'true'
          } else {
            this.$message({
                message: '测试失败,http code:'+ response.data.status,
                type: 'warning',
                duration: 5 * 1000,
                center: true
              })
             this.testResult = response.data.message
             this.testPass = 'false'
          }
            
        })
    },
    makeParam(p) {
      return JSON.stringify(p,null,'\t')
    },
    initTest() {
      this.testParams = {}
      this.tfeAddr = process.env.TFE_ADDR
      this.testResult = ''
      this.testPass = ''
      this.testParams.api_name=this.showName(this.api.api_id)
      this.testParams.api_version = this.showVersion(this.api.api_id)  
      // load old params
      var params = {
          target_app: 'juzManage',
          target_path: '/manage/api/queryParam',
          api_id: this.api.api_id
        }
      proxy('GET',params).then(res => {
          console.log(res.data.data)
        if (res.data.data != '') {
            this.testParams = JSON.parse(res.data.data)
            if (this.testParams.debug_on) {
                this.debugOn = true
            }
        }
      })
    }
  },
  mounted() {
      this.initTest()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
