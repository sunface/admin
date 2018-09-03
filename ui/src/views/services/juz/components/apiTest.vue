<template>
  <div class="api-test">
    <el-form class="black-form"  ref="dataForm" label-position="left" label-width="70px" >
        <el-form-item label="API ID"  style="width:300px">
            {{api.api_id}}
        </el-form-item>

        <el-form-item label="Backend URL"  style="width:300px">
            {{api.route_addr}}
        </el-form-item>
        
        <el-form-item label="Method">
            <el-select  class="filter-item" v-model="testMethod" style="width: 200px"  placeholder="选择HTTP Method">
            <el-option label="GET" value="GET"></el-option>
            <el-option  label="POST" value="POST"></el-option>
            <el-option  label="PUT" value="PUT"></el-option>
            <el-option  label="DELETE" value="DELETE"></el-option>
            </el-select>
        </el-form-item>

        <el-form-item label="JUZ API Address" style="width:400px">
            {{tfeAddr}}
        </el-form-item>


        <el-form-item label="Params" style="width:600px">
            <el-input type="textarea" :autosize="{ minRows: 6, maxRows: 8}" :value="makeParam(testParams)" :disabled=true>
            </el-input>
        </el-form-item>
        <el-form-item label="Add Param">
            <el-input   style="width:150px" v-model="tempKey">
                <template slot="prepend">Key</template>
            </el-input>
            <el-input  style="width:150px" v-model="tempVal">
                <template slot="prepend">Val</template>
            </el-input>
            <el-button @click="addParam" type="success" size="medium">Add</el-button>
            <el-button @click="clearParam" type="warning" size="medium">Clear All</el-button>
        </el-form-item>

        <el-form-item label="Output" style="width:600px">
            <el-input type="textarea" :autosize="{ minRows: 10, maxRows: 12}" placeholder="Test Result" v-model="testResult" :disabled=true>
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
        testMethod: 'GET',
        testPass : '',
        testParams: {}, 
        tempKey:'',
        tempVal:'',

        tempBWKey:'',
        tempBWVal:'',

        inRefresh : false
      }
  },
  watch: {
    api() {
      this.initTest()
    }
  },
  methods: {
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
                method: this.testMethod,
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
            
        }).catch(error => {
        })
    },
    clearParam() {
      this.testParams = {}
      this.testParams.api_name=this.showName(this.api.api_id)
      this.testParams.api_version = this.showVersion(this.api.api_id)
    },
    addParam() {
      this.testParams[this.tempKey] = this.tempVal
      this.tempKey = ''
      this.tempVal = ''
    },
    makeParam(p) {
      return JSON.stringify(p,null,'\t')
    },
    initTest() {
      this.tfeAddr = process.env.TFE_ADDR
      this.testResult = ''
      this.testPass = ''      
      this.testParams.api_name=this.showName(this.api.api_id)
      this.testParams.api_version = this.showVersion(this.api.api_id)
    }
  },
  created() {
      this.initTest()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
