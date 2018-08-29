<template>
  <div class="app-container">
      <div class="filter-container">
        <el-select  class="filter-item" v-model="selectedImport" style="width: 200px"  placeholder="请选择导入类型">
          <el-option v-for="s in  imports" :key="s.value" :label="s.label" :value="s.value">
          </el-option>
        </el-select>
        
        <el-row style="margin-top:10px;">
            <el-col :span=12>
                <el-alert
                    title="提示：请先在相应的地方复制配置，然后选择上面的导入类型，再再下方粘贴导入"
                    :closable=false
                    type="warning">
                </el-alert>
                <el-input type="textarea" :autosize="{ minRows: 10, maxRows: 50}" placeholder="配置JSON" v-model="content"></el-input>
                <el-button type="success" style="margin-top:6px" @click="beginImport('edit')">更新目标</el-button>
                <el-button  style="margin-top:6px" @click="beginImport('create')">创建目标</el-button>
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
          imports :[{label:'API',value:2},{label:'策略',value:3}],
          selectedImport: 2,
          content: ''
      }
  },
  methods: {
      beginImport(type) {
          if (this.selectedImport=='') {
            this.$message({
                message: '请先选择导入类型',
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
                        target_service: 'tfe.manage',
                        target_path: '/manage/api/define',
                        api : o,
                        action: action
                    }
                    proxy('POST',params).then(res => {
                        this.content = ''
                        this.$message({
                            message: '导入API成功',
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
                            target_service: 'tfe.manage',
                            target_path: '/manage/strategy/create',
                            strategy: this.content
                        }
                        proxy('POST',params).then(res => {
                            this.content = ''
                            this.$message({
                                message: '策略创建成功',
                                type: 'success',
                                duration: 3 * 1000,
                                center: true
                            })
                        })
                  } else {
                    // 创建策略
                    var params = {
                        target_service: 'tfe.manage',
                        target_path: '/manage/strategy/update',
                        strategy: this.content
                    }
                    proxy('POST',params).then(res => {
                        this.content = ''
                        this.$message({
                            message: '策略更新成功',
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
