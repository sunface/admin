<template>
  <div class="app-container">
      <div class="filter-container">
        <el-tag>select service</el-tag>
        <el-select  class="filter-item" :value="calcService()" @change='handleSelService' style="width: 200px"  placeholder="Selete a service">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>

         <!-- <span style="float:right;margin-right:30px">
          <el-input  style="width: 250px;" class="filter-item" placeholder="模糊搜索" v-model="q" clearable @clear="loadLogs(selectedService)">
          </el-input>
          <el-button class="filter-item" type="success" v-waves icon="el-icon-search" @click="loadLogs(selectedService)"></el-button>
        </span> -->
      </div>
      <auditLog :totalLogs="totalLogs" :logs="logs" themeType='white-theme' :targetType=0 :targetID="sName" :q="q"></auditLog>
  </div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
import auditLog from './components/auditLog'
import request from '@/utils/request' 
export default {
  name: 'auditLogs',
  data() {
      return {
        services: [{name:'All'}],
        totalLogs: 0,
        logs: [],
        sName : '',

        selectedService: '',

        q: ''
      }
  },
  components:{auditLog},
  methods :{
    calcService() {
        return this.selectedService
    },
    handleSelService(val) {
       this.selectedService = val
       if (val=='All') {
         this.sName = this.genAllService(this.services)
       } else {
         this.sName = "'"+val+"'"
       }  
       this.loadLogs(this.sName)
    },
    loadLogs(service) {
      //加载审计日志
      var params = {
          target_app: 'juzManage',
          target_path: '/manage/auditLog/count',
          target_type: 0,
          target_id: service
      }
      proxy('POST',params).then(res => {
          this.totalLogs = res.data.data
      })

      var params = {
          target_app: 'juzManage',
          target_path: '/manage/auditLog/load',
          target_type: 0,
          target_id: service,
          q: this.q,
          page: 1
      }
      proxy('POST',params).then(res => {
          this.logs = res.data.data
      })
    },
    genAllService(services) {
      var s = ""
      for (var i=0;i<services.length;i++) {
        if (i==0) {
          s = "'"+services[i].name+"'"
        } else {
          s = s + ',' + "'"+services[i].name+"'"
        }
      }
      return s
    },
    loadServices() {
       request({
          url: '/infra/service/query',
          method: 'GET', 
          params: {
          }
        }).then(res => {
            this.services = res.data.data
            this.selectedService = 'All'
             this.sName = this.genAllService(res.data.data)
             this.loadLogs(this.sName)
        })
    }
  },
  mounted() {
    this.loadServices()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
