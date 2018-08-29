<template>
  <div class="app-container">
    <div class="filter-container">
        <el-select clearable class="filter-item" v-model="selectedService" @change='handleSelService' style="width: 200px"  placeholder="请选择Service">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="createVisible=true"  v-waves type="primary" icon="el-icon-edit">创建标签</el-button>
    </div>
  
    <div >
        <el-tag
            class="hover-cursor"
            v-for="(l,i) in labels"
            :key="l"
            :type="calcLabelType(i)"
            size='large'
            style="margin-right:15px"
            @click.native="setLabel(l)">
            {{l}}
        </el-tag>
    </div>
     <el-alert
            :closable=false
            title="点击上面的标签查看监控统计"
            type="success"
            style="width:250px;margin-top:10px">
        </el-alert> 

    <el-dialog :title="'当前标签 ：'+selectedLabel" :visible.sync="metricsVisible" top="40px"  :fullscreen=true style="backgroud:#545c64 !important">
        <iframe :src="labelGrafanaAddr(selectedLabel)" style="height:910px;width:100%;border:none;margin-top:-40px;"></iframe>
    </el-dialog>

     <el-dialog title="创建标签" :visible.sync="createVisible"  style="backgroud:#545c64 !important">
         <el-input v-model="tempLabel" placeholder="输入标签名" style="width:500px"></el-input>
        <div slot="footer" class="dialog-footer">
            <el-button @click="createVisible=false">取消</el-button>
            <el-button  type="primary"  @click="createLabel">提交</el-button>
          </div>
    </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
export default {
  name: 'labels',
  data() {
      return {
          services: [],
          labels: [],
          selectedService: '',
          selectedLabel: '',
          metricsVisible: false,
          createVisible:false,
          tempLabel: ''
      }
  },
  methods:{
    createLabel() {
        var params = {
            target_service: 'tfe.manage',
            target_path: '/manage/labels/create',
            service: this.selectedService,
            name: this.tempLabel
        }
        proxy('POST',params).then(res => {
            this.$message({
                message: '标签创建成功',
                type: 'success',
                duration: 3 * 1000,
                center: true
            }) 
        })
    } ,
    setLabel(l) {
        this.selectedLabel = l
        this.metricsVisible = true
    },
    labelGrafanaAddr(l) {
      var addr = process.env.LABEL_GRAFANA + l
       return  addr
    },
    calcLabelType(i) {
        var tp = ''
        switch (i % 5) {
            case 0:
                tp = ''
                break;
            case 1: 
                tp = 'success'
                break
            case 2: 
                tp = 'info'
                break
            case 3: 
                tp = 'warning'
                break
            case 4:
                tp = 'error'
                break
            default:
                break;
        }
        return tp
    },
    handleSelService() {
        if (this.selectedService == '') {
            this.labels = []
            return 
        }
        this.loadLabels()
    },
    loadLabels() {
       //加载审计日志
        var params = {
            target_service: 'tfe.manage',
            target_path: '/manage/labels/query',
            service: this.selectedService
        }
        proxy('POST',params).then(res => {
            this.labels = res.data.data
        })
    },
    loadServices() {
      // 加载该用户的所有service
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/service/query',
            }
        proxy('POST',params).then(res => {
           if (res.data.status == 200) {
             this.services = res.data.data
             if (this.services.length > 0) {
               this.selectedService = this.services[0].name
               this.loadLabels()
             }
           }
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
