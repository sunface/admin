<template>
  <div  :class="{'black-theme':isBlackTheme(),'audit-log':true}">
      <div class="table">
        <el-table  :data="logs"  fit highlight-current-row style="width: 100%;"  :default-sort = "{prop: 'gmt_modified', order: 'descending'}">
          <el-table-column align="left" label="发生时间" width="160" prop="gmt_modified">
            <template slot-scope="scope">
              <span>{{scope.row.gmt_modified}}</span>
            </template>
          </el-table-column>
         <el-table-column width="100" align="left" label="目标类型">
            <template slot-scope="scope">
              <span v-if="scope.row.target_type==1">Service</span>
              <span v-else-if="scope.row.target_type==2">API</span>
              <span v-else-if="scope.row.target_type==3">策略</span>
              <span v-else-if="scope.row.target_type==4">设置权限</span>
              <span v-else-if="scope.row.target_type==5">批量设置</span>
            </template>
          </el-table-column>
          <el-table-column width="250" align="left" label="目标名称(ID)">
            <template slot-scope="scope">
              <span>{{scope.row.target_id}}</span>
            </template>
          </el-table-column>
          <el-table-column width="150" align="left" label="目标内容">
            <template slot-scope="scope">
              <span v-if="scope.row.target_type==1">
                <span v-if="scope.row.content != ''">{{scope.row.content}}</span>
                <span v-else>无描述</span>
              </span>
              <span v-if="scope.row.target_type==4">{{scope.row.content}}</span>
              <span v-if="scope.row.target_type==3">
                  <i class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="showStrategy(scope.row.content)"></i> 
              </span>
              <span v-if="scope.row.target_type==2">
                  <i v-if="scope.row.content!=''" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="viewApi(scope.row.content)"></i> 
              </span>
              <span v-if="scope.row.target_type==5">
                  <el-popover
                    placement="top"
                    width="800"
                    trigger="click">
                    {{JSON.parse(scope.row.content)}}
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a"></i> 
                </el-popover>
              </span>
            </template>
          </el-table-column>
          <el-table-column width="150" align="left" label="Service">
            <template slot-scope="scope">
              <span>{{scope.row.service}}</span>
            </template>
          </el-table-column>
          <el-table-column width="100" align="left" label="操作类型">
            <template slot-scope="scope"> 
              <span v-if="scope.row.op_type==1">创建</span>
              <span v-else-if="scope.row.op_type==2">编辑</span>
              <span v-else-if="scope.row.op_type==3">发布</span>
              <span v-else-if="scope.row.op_type==4">下线</span>
              <span v-else-if="scope.row.op_type==5">删除</span>
            </template>
          </el-table-column>
          <el-table-column width="200" align="left" label="操作人" prop="user_id">
            <template slot-scope="scope">
              <span>{{scope.row.user_id}}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="管理" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button  v-if="scope.row.content!=''" class="green-button"  type="text" @click="handleCopy(scope.row.content,$event)">复制配置</el-button>
            </template>
          </el-table-column>
        </el-table>

         <el-pagination
            v-show="totalLogs>0"
            layout="total,prev, pager, next"
            @current-change="loadLogs"
            :total="totalLogs">
        </el-pagination>
      </div>

      <el-dialog  title="API定义" :visible.sync="apiViewVisible" top="20px">
        <apiDefine :api="tempApi"></apiDefine>
        <div slot="footer" class="dialog-footer">
          <el-button @click="apiViewVisible = false">关闭</el-button>
        </div>
      </el-dialog>
      <el-dialog  title="黑白名单" :visible.sync="bwVisible" top="20px">
        <bwlist :strategy="strategy"></bwlist>
        <div slot="footer" class="dialog-footer">
          <el-button @click="bwVisible = false">关闭</el-button>
        </div>
      </el-dialog>
       <el-dialog  title="超时重试" :visible.sync="retryVisble" top="20px">
        <retry :strategy="strategy"></retry>
        <div slot="footer" class="dialog-footer">
          <el-button @click="retryVisble = false">关闭</el-button>
        </div>
      </el-dialog>
      <el-dialog  title="超时重试" :visible.sync="trafficVisible" top="20px">
        <traffic :strategy="strategy"></traffic>
        <div slot="footer" class="dialog-footer">
          <el-button @click="trafficVisible = false">关闭</el-button>
        </div>
      </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
import apiDefine from '../components/apiDefine'
import bwlist from './bwlist'
import retry from './retry'
import traffic from './traffic'

import clipboard from '@/directive/clipboard/index.js'
import clip from '@/utils/clipboard'
export default {
  name: 'auditLog',
  props: {
    totalLogs: {
        type: Number,
        default: 0
    },
    logs : {
        type: Array,
        default: []
    },
    themeType: {
        type: String,
        default: 'white-theme' 
    },
    targetType: {
        type: Number,
        default: 0
    },
    targetID: {
        type: String,
        default: ''
    },
    q : {
      type: String,
      default: ''
    }
  },
  directives: {
    clipboard
  },
  components: {apiDefine,retry,bwlist,traffic},
  data() {
      return {
          tempApi: {},
          apiViewVisible : false,
          strategy: {
              content:{}
          },
          bwVisible: false,
          retryVisble: false,
          trafficVisible: false
      }
  },
  methods: {
    showStrategy(content) {
        this.strategy = JSON.parse(content)
        this.strategy.content = JSON.parse(this.strategy.content)
        switch (this.strategy.type) {
            case 1:
                this.bwVisible=true
                break;
            case 2:
                this.retryVisble=true
                break;
            case 3:
                this.trafficVisible = true
            default:
                break;
        }
    },
    isBlackTheme() {
        return this.themeType== 'black-theme'
    },
    viewApi(api) {
        this.tempApi = JSON.parse(api)
        this.apiViewVisible = true
    },

    handleCopy(text, event) {
      clip(text, event)
      this.$message({
        message: '复制成功',
        type: 'success',
        duration: 2000
      })
    },
    loadLogs(page) {
        var params = {
            target_service: 'tfe.manage',
            target_path: '/manage/auditLog/load',
            target_type: this.targetType,
            target_id: this.targetID,
            q: this.q,
            page: page
        }
        proxy('POST',params).then(res => {
            this.logs = res.data.data
        })
    },
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

