<template>
  <div class="app-container">
      <div class="filter-container">
        <el-tag>select service</el-tag>
        <el-select clearable class="filter-item" :value="calcService()" @change='handleSelService' style="width: 200px"  placeholder="select a service">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"   type="primary" icon="el-icon-edit">Add Strategy</el-button>
      </div>
      
      <div class="table">
        <el-table  :data="strategies"  fit highlight-current-row style="width: 100%;" >
          <el-table-column align="left" label="Type" width="150" prop="name">
            <template slot-scope="scope">
              <span v-if="scope.row.type==1">White/Black List</span>
              <span v-if="scope.row.type==2">Timeout/Retry</span>
              <span v-if="scope.row.type==3">Traffic Controls</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="Name" width="200" prop="name">
            <template slot-scope="scope">
              <span>{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column width="180" align="left" label="Update Date" prop="route_addr">
            <template slot-scope="scope">
              <span>{{scope.row.modify_date}}</span>
            </template>
          </el-table-column>
          <el-table-column width="150" align="left" label="Status" prop="route_addr">
            <template slot-scope="scope">
              <span v-if="scope.row.status==0"><el-tag type="warning" style="border:none;">Off</el-tag></span>
              <span v-else><el-tag type="success" style="border:none;">On</el-tag></span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="Operate" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button size="small" type="text" @click="handleEdit(scope.row)">Edit</el-button>
              <el-button  class="green-button" size="small" type="text" @click="handleCopy(scope.row,$event)">Copy config</el-button>
              <el-button size="small" type="text" @click="changeStatus(scope.row)">
                <span v-if="scope.row.status==0">Start</span>
                <span v-else>Stop</span>
              </el-button>
              <el-button  size="small" type="text" @click="delStrategy(scope.row)">Delete</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <!-- 策略面板 -->
      <el-dialog :title="dialogTitle" :visible.sync="bwlistVisible">
        <el-form  label-position="left" label-width="120px" style='width: 650px; margin-left:50px;' size="mini">
            <div class="form-block">
                <span>Basic</span>
                <el-form-item label="Strategy Name" style="margin-top:10px">
                    <el-input  placeholder="" style="width:300px" v-model="tempStrategy.name">
                    </el-input>
                </el-form-item>
                <el-form-item label="Type" style="margin-top:10px">
                    <el-radio-group v-model="tempStrategy.type" @change="selType" :disabled="!isCretaeSubmit">
                        <el-radio :label="1">White/Black List</el-radio>
                        <el-radio :label="2">Timeout/Retry</el-radio>
                        <el-radio :label="3">Traffic Control</el-radio>
                    </el-radio-group>
                </el-form-item>
            </div>
            <div class="form-block" v-if="tempStrategy.type==1">
                <span>White/Black List</span>
                <el-form-item label="Type" style="margin-top:10px">
                    <el-radio-group v-model="tempStrategy.sub_type" >
                        <el-radio :label="1">Black</el-radio>
                        <el-radio :label="2">White</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-alert
                    title="If the list is long ,you can input and save at the the same time"
                    :closable=false
                    type="warning">
                </el-alert>
                <el-form-item label="Add" v-if="tempStrategy.sub_type==1 || tempStrategy.sub_type==2" style="margin-top:10px">
                    <el-radio-group v-model="selbwType" @change="selBwType">
                        <el-radio :label="1">IP</el-radio>
                        <el-radio :label="2">Param</el-radio>
                    </el-radio-group>
                    <div>
                        <el-input  placeholder="Key" style="width:80px" v-model="tempBWKey" :disabled="selbwType==1">
                        </el-input>
                        <el-input  placeholder="Val" style="width:160px" v-model="tempBWVal">
                        </el-input>
                        <el-button @click="addBW" type="success" size="mini" icon="el-icon-plus" circle style="margin-left:10px"></el-button>
                    </div>
                </el-form-item>
                <el-form-item label="Current List" v-if="tempStrategy.sub_type==1 || tempStrategy.sub_type==2">
                    <div v-for="kv in tempStrategy.content"  :key="kv.key" style="margin-top:-2px">
                        <el-tag type="success" size="large" style="width:70px;border:none;">{{kv.key}}</el-tag>
                        <el-tag type="info" size="large" style="width:200px;border:none">{{kv.val}}</el-tag>
                        <el-button @click.native="removeBW(kv)" type="text" size="mini" icon="el-icon-minus" style="margin-left:10px"></el-button>
                    </div>
                </el-form-item>
            </div>
            <div class="form-block" v-if="tempStrategy.type==2">
              <span>Request Timeout/Retry</span>
              <el-form-item label="Timeout(Second)" style="width:200px;margin-top:10px" >
                <el-tooltip content="0<X<=60" placement="top">
                    <el-input-number v-model="tempStrategy.content.req_timeout"  :min="1" :max="60" ></el-input-number>
                </el-tooltip>
              </el-form-item>
              <el-form-item label="Retry Times" style="width:200px">
                <el-tooltip content="When error happens, we will retry after some time" placement="top">
                    <el-input-number v-model="tempStrategy.content.retry_times"  :min="0" :max="5" ></el-input-number>
                </el-tooltip>
              </el-form-item>
              <el-form-item label="Retry Interval" style="width:200px" >
                <el-tooltip content="the interval when next retry happens, in seconds" placement="top">
                    <el-input-number v-model="tempStrategy.content.retry_interval"  :min="1" :max="30" ></el-input-number>
                </el-tooltip>
              </el-form-item>
            </div>
            <div v-if="tempStrategy.type==3">
              <div class="form-block" >
                  <span>API Traffic</span>
                  <el-form-item label="QPS" style="margin-top:10px" >
                      <el-input-number v-model="tempStrategy.content.qps"  :min="0" :max="10000" ></el-input-number>
                      <el-alert
                          title="Max QPS per second, 0 means no limit"
                          :closable=false
                          type="success">
                      </el-alert>
                    </el-form-item>
                    <el-form-item label="Concurrents">
                      <el-input-number v-model="tempStrategy.content.concurrent"  :min="0" :max="10000" ></el-input-number>
                      <el-alert
                          title="Max concurrents per second, 0 means no limit"
                          :closable=false
                          type="success">
                      </el-alert>
                    </el-form-item>
              </div>
              <div class="form-block" >
                  <span>User Quota</span>
                  <el-form-item label="Param" style="margin-top:10px;width:400px" >
                    <el-input v-model="tempStrategy.content.param" placeholder="empty means no limit"></el-input>
                  </el-form-item>
                  <el-form-item label="Time Span">
                    <el-input-number v-model="tempStrategy.content.span"  :min="2" :max="2592000" ></el-input-number>Second
                  </el-form-item>
                  <el-form-item label="Times">
                    <el-input-number v-model="tempStrategy.content.times"  :min="1" :max="10240000" ></el-input-number>
                    <el-alert
                      title="In the limit time span, the value of the param, only be allowed to access limit times"
                      :closable=false
                      type="success">
                    </el-alert>
                    <el-alert
                        title="e.g. param=mobile|time span=86400|times=1000, means each mobile number only be allowed to access 1000 times in 86400 seconds"
                        :closable=false
                        type="success">
                    </el-alert>
                    </el-form-item>
              </div>
            </div>
        </el-form>
 
        <div slot="footer" class="dialog-footer">
          <el-button @click="bwlistVisible=false">Cancel</el-button>
          <el-button  type="primary" v-if="isCretaeSubmit" @click="submitStrategy">
            <span>Submit</span>
          </el-button>
          <el-button  type="primary" v-else @click="submitEdit">
            <span>Submit</span>
          </el-button>
        </div>
      </el-dialog>

  </div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
import clipboard from '@/directive/clipboard/index.js'
import clip from '@/utils/clipboard'
import request from '@/utils/request'   
export default {
  name: 'strategy',
  cmponents:{clipboard,clip},
  data() {
      return {
        services:[],
        selectedService: '',

        strategies: [],

        dialogStatus: '',
        dialogTitle: '',

        tempStrategy: {
            content:{}
        },

        // 黑白名单
        bwlistVisible: false,
        tempBWKey: '',
        tempBWVal: '',
        selbwType: 1,

        apiShow: false,

        isCretaeSubmit: false
      }
  },
  methods: {
    handleSelService(service) {
       this.$store.dispatch('setService', service)
       this.selectedService = service
       this.loadStrategy(service)
    },
     handleCopy(api, event) {
      var text = JSON.stringify(api)
      clip(text, event)
      this.$message({
        message: 'Copied',
        type: 'success',
        duration: 2000
      })
    },
    changeStatus(s) {
      var str = ''
      if (s.status == 0) {
        str = 'You will start using strategy："' + s.name + '"！When started, all the apis using will  start using this strategy！'
      } else {
         str = 'You will stop using strategy："' + s.name + '"！When stopped, all the apis will stop using this strategy'
      }
      this.$confirm( str, 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/strategy/change',
              id: s.id,
              name : s.name,
              status: s.status
          }
          proxy('POST',params).then(res => {
            if (s.status==0) {
              s.status = 1
            } else {
              s.status = 0
            }
            this.$message({
              message: 'Status changed ok',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
          })
        })
    },
    delStrategy(s) {　
      this.$confirm('You will delete this strategy："' + s.name + '"! When deleted, all the apis will detele this strategy' , 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/strategy/delete',
              id: s.id,
              name: s.name,
              service: s.service,
              type: s.type
          }
          proxy('POST',params).then(res => {
            this.loadStrategy(this.selectedService)
            this.$message({
              message: 'Delete strategy ok',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
          })
        })
    },
    apiSet(s) {
      this.apiShow = true
      this.tempStrategy = s
    },
    submitEdit() {
      if (this.tempStrategy.name == '') {
        this.$message({
          message: 'Strategy name cant be empty',
          type: 'warning',
          duration: 3 * 1000,
          center: true
      })
      }
     if (this.tempStrategy.type == 3) {
        this.tempStrategy.content.param = this.tempStrategy.content.param.trim()
      }

    var s = Object.assign({},   this.tempStrategy) 
    s.content = JSON.stringify(s.content)
      // 创建策略
      var params = {
        target_app: 'juzManage',
        target_path: '/manage/strategy/update',
        strategy: s
      }
      proxy('POST',params).then(res => {
        //黑白名单类型为了方便输入保存，提交后不关闭面板,其它类型要关闭面板
       if (this.tempStrategy.type != 1) {
          this.bwlistVisible = false
        }
        this.loadStrategy(this.selectedService)
        this.$message({
          message: 'Submit ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        })
      })
    },
    handleEdit(s) {
      this.dialogStatus = 'edit'
      this.dialogTitle = 'Edit Strategy'
      this.tempStrategy = Object.assign({}, s) 
      this.tempStrategy.content = JSON.parse(this.tempStrategy.content)
      this.bwlistVisible = true
      this.tempBWKey = 'ip'
      this.isCretaeSubmit = false
    },
    //黑白名单
    removeBW(kv) {
      for (var i=0;i<this.tempStrategy.content.length;i++) {
        if (kv.key == this.tempStrategy.content[i].key && kv.val == this.tempStrategy.content[i].val) {
          this.tempStrategy.content.splice(i,1)
        }
      }
    },
    addBW() {
      if (this.tempBWKey == '' || this.tempBWVal == '') {
        this.$message({
            message: 'key、val cant be empty',
            type: 'warning',
            duration: 3 * 1000,
            center: true
        })
        return 
      }

      // 名单是否存在
      for ( var i=0;i<this.tempStrategy.content.length;i++) {
          if (this.tempBWKey == this.tempStrategy.content[i].key && this.tempBWVal == this.tempStrategy.content[i].val) {
            this.$message({
                message: 'List item already exist',
                type: 'warning',
                duration: 3 * 1000,
                center: true
            })
            return 
          }
      }
      this.tempStrategy.content.push({
        key: this.tempBWKey,
        val: this.tempBWVal
      })
      if (this.selbwType==1) {
        this.tempBWKey = 'ip'
      }  else {
        this.tempBWKey = ''
      }
      this.tempBWVal = ''
    },
    submitStrategy() {
       if (this.tempStrategy.name == '') {
         this.$message({
            message: 'Strategy name cant empty',
            type: 'warning',
            duration: 3 * 1000,
            center: true
        })
       }

      if (this.tempStrategy.type == 3) {
        this.tempStrategy.content.param = this.tempStrategy.content.param.trim()
      }
      var s = Object.assign({},   this.tempStrategy) 
      s.content = JSON.stringify(s.content)
      // 创建策略
      var params = {
        target_app: 'juzManage',
        target_path: '/manage/strategy/create',
        strategy: s
      }
      proxy('POST',params).then(res => {
          this.isCretaeSubmit = false
          if (this.tempStrategy.type != 1) {
              this.bwlistVisible = false
          }
        this.loadStrategy(this.selectedService)
        this.$message({
          message: 'Submit ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        })
      })
    },
    selBwType(val) {
      if (val == 1) {
        this.tempBWKey = 'ip'
      } else {
        this.tempBWKey = ''
      }
    },
    selType(type) {
        switch (type) {
            case 1:
                this.tempStrategy = {
                    name: this.tempStrategy.name,
                    type: 1,
                    sub_type: 1,
                    content: [],
                    service: this.selectedService,
                }
                break;
            case 2:
                this.tempStrategy = {
                    name: this.tempStrategy.name,
                    type: 2,
                    sub_type: 0,
                    service: this.selectedService,
                    content: {
                        req_timeout: 15,
                        retry_times: 0,
                        retry_interval: 3
                    },
                }
                break
            case 3:
              this.tempStrategy = {
                    name: this.tempStrategy.name,
                    type: 3,
                    sub_type: 0,
                    service: this.selectedService,
                    content: {
                        qps: -1,
                        concurrent: -1,

                        param: '', // 限制参数名
                        span: 2, // SPAN时间(秒)
                        times: 1, // 访问TIMES(次数)

                        // 熔断触发条件
                        fuse_error: 0, // 熔断错误率，大于该值时，触发熔断保护
                        fuse_error_count: 20,   // 熔断触发的最小请求次数
                        // 熔断保护
                        fuse_percent: 50, // 熔断触发后，允许访问的百分比，例如100次请求，只有50次允许通过
                        // 熔断恢复条件
                        fuse_recover: 25, // 熔断错误率小于该值时，取消熔断保护
                        fuse_recover_count: 10 // 熔断恢复的最小请求次数
                    },
                }
                break
            default:
                break;
        }
    },
    handleCreate() {
      if (this.selectedService == '') {
        this.$message({
            message: 'Selete a service first',
            type: 'warning',
            duration: 3 * 1000,
            center: true
        })
        return 
      }
      this.dialogStatus = 'create'
      this.dialogTitle = 'Create Strategy'
      this.tempBWKey = 'ip'

      this.bwlistVisible = true
      this.isCretaeSubmit = true
      this.tempStrategy = {
        type: 1,
        sub_type: 1,
        content: [],
        name: '',
        service: this.selectedService
      }
    },
    selStrategy(val) {
      this.selectedStrategy = val
    },
    loadStrategy(service) {
       var params = {
          target_app: 'juzManage',
          target_path: '/manage/strategy/load',
          service: this.selectedService,
          type: 0,
        }
        proxy('POST',params).then(res => {
          this.strategies = res.data.data
        })
    },
    calcService() {
        return this.selectedService || this.$store.getters.service
    },
    loadServices() {
      request({
          url: '/infra/service/query',
          method: 'GET', 
          params: {
          }
        }).then(res => {
            this.services = res.data.data
        })
    }
  },
  created() {
    this.loadServices()
    this.selectedService = this.$store.getters.service
    if (this.selectedService != '') {
      this.loadStrategy(this.selectedService)
    }
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
