<template>
  <div class="app-container">
    <div class="filter-container">
        <el-tag type="info"> service</el-tag>
        <el-select clearable class="filter-item" :value="calcService()" @change='handleSelService' style="width: 160px"  placeholder="select service">
            <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
            </el-option>
        </el-select>
    </div>
    <div class="filter-container">
        <el-tag type="info">app</el-tag>
        <el-select clearable class="filter-item margin-left-20" v-model="selectedApp" style="width: 160px"  placeholder="select app">
            <el-option v-for="s in  apps" :key="s.name" :label="s.name" :value="s.name">
            </el-option>
        </el-select>
    </div>
    <el-alert
        title="You can view either service or app's metrics"
        :closable=false
        style="width:400px"
        type="warning">
    </el-alert>
    <div class="filter-container margin-top-10">
        <el-button class="filter-item" type="success" icon="el-icon-edit" @click="viewMetrics">View metrics</el-button>
    </div>

    <el-dialog :title="'now' + targetType + '：'+ targetName" :visible.sync="metricsVisible" top="40px"  :fullscreen=true style="backgroud:rgb(57, 79, 90) !important">
        <iframe :src="viewUrl" style="height:910px;width:100%;border:none;margin-top:-40px;"></iframe>
    </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request'   
export default {
  name: 'metrics',
  data() {
      return {
          services: [],
          selectedService: '',

          apps: [],
          selectedApp: '',

          metricsVisible: false,

          targetType: '',
          targetName : '',
          viewUrl : ''
      }
  },
  methods: {
    viewMetrics() {
        if (this.selectedApp != '') {
            this.targetName = this.selectedApp
            this.viewUrl = process.env.APP_GRAFANA + this.targetName
            this.targetType = 'Service'
            this.metricsVisible = true
            return 
        }

        if (this.selectedService != '') {
            this.targetName = this.selectedService
            this.viewUrl = process.env.SERVICE_GRAFANA + this.targetName
            this.targetType = 'APP'
            this.metricsVisible = true
            return 
        }

        this.$message({
            message: 'You have not select any service or app',
            type: 'warning',
            duration: 3000
        })


    },
    handleSelService(service) {
        this.$store.dispatch('setService', service)
        this.selectedService = service
    },
    calcService() {
        return this.selectedService || this.$store.getters.service
    },
    loadApps() {
        request({
          url: '/ops/app/query',
          method: 'GET', 
          params: {
              service: this.selectedService
          }
        }).then(res => {
            this.apps = res.data.data
        })
    },
    loadServices() {
      request({
          url: '/ops/service/query',
          method: 'GET', 
          params: {
          }
        }).then(res => {
            this.services = res.data.data
        })
    }
  },
  mounted() {
      this.loadServices()
      this.selectedService = this.$store.getters.service
      if (this.selectedService != '') {
          this.loadApps()
      }
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
