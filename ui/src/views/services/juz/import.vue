<template>
  <div class="app-container">
      <div class="filter-container">
        <el-select  class="filter-item" v-model="selectedImport" style="width: 200px"  placeholder="select impor type">
          <el-option v-for="s in  imports" :key="s.value" :label="s.label" :value="s.value">
          </el-option>
        </el-select>
        
        <el-row style="margin-top:10px;">
            <el-col :span=12>
                <el-alert
                    :title="$t('juz.copyConfigTips')"
                    :closable=false
                    type="warning">
                </el-alert>
                <el-input type="textarea" :autosize="{ minRows: 10, maxRows: 50}" placeholder="Config JSON" v-model="content"></el-input>
                <el-button  style="margin-top:6px" @click="beginImport('edit')">Update</el-button>
                <el-button type="success" style="margin-top:6px" @click="beginImport('create')">Create</el-button>
            </el-col>
        </el-row>
      </div>
  </div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
export default {
  name: 'import',
  data() {
      return {
          imports :[{label:'API',value:2},{label:'Strategy',value:3}],
          selectedImport: 2,
          content: ''
      }
  },
  methods: {
      beginImport(type) {
          if (this.selectedImport=='') {
            this.$message({
                message: 'select import type',
                type: 'warning',
                duration: 5 * 1000,
                center: true
            }) 
            return 
          }
        
          switch (this.selectedImport) {
              case 2://导入API
                  var action = ''
                  if (type=='create') {
                    // 创建目标，id要重置为0
                    action = 'create'
                  }
                  var o = JSON.parse(this.content)
                  o.param_rules = JSON.stringify(o.param_rules)
                  var params = {
                        target_app: 'juzManage',
                        target_path: '/manage/api/define',
                        api : o,
                        action: action
                    }
                    proxy('POST',params).then(res => {
                        this.content = ''
                        this.$message({
                            message: 'imported ok',
                            type: 'success',
                            duration: 3 * 1000,
                            center: true
                        })
                    })
                  break;
            
              case 3://导入策略
                  if (type=='create') {
                        // 创建策略
                        var params = {
                            target_app: 'juzManage',
                            target_path: '/manage/strategy/create',
                            strategy: this.content
                        }
                        proxy('POST',params).then(res => {
                            this.content = ''
                            this.$message({
                                message: 'Strategy created ok',
                                type: 'success',
                                duration: 3 * 1000,
                                center: true
                            })
                        })
                  } else {
                    // 创建策略
                    var params = {
                        target_app: 'juzManage',
                        target_path: '/manage/strategy/update',
                        strategy: this.content
                    }
                    proxy('POST',params).then(res => {
                        this.content = ''
                        this.$message({
                            message: 'Strategy created ok',
                            type: 'success',
                            duration: 3 * 1000,
                            center: true
                        })
                    })
                  }
                  break;
              default:
                  break;
          }
      }
  },
  created() {
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
