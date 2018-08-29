<template>
  <div class="app-container">
      <div class="filter-container">
        <el-select clearable class="filter-item" v-model="selectedService" @change='handleSelService' style="width: 200px"  placeholder="请选择Service">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"   type="primary" icon="el-icon-edit">新建策略</el-button>
      </div>
      
      <div class="table">
        <el-table  :data="strategies"  fit highlight-current-row style="width: 100%;" >
          <el-table-column align="left" label="类型" width="150" prop="name">
            <template slot-scope="scope">
              <span v-if="scope.row.type==1">黑白名单</span>
              <span v-if="scope.row.type==2">超时重试</span>
              <span v-if="scope.row.type==3">流量控制</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="名称" width="200" prop="name">
            <template slot-scope="scope">
              <span>{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column width="180" align="left" label="更新时间" prop="route_addr">
            <template slot-scope="scope">
              <span>{{scope.row.gmt_modified}}</span>
            </template>
          </el-table-column>
          <el-table-column width="150" align="left" label="当前状态" prop="route_addr">
            <template slot-scope="scope">
              <span v-if="scope.row.status==0"><el-tag type="warning" style="border:none;">已停用</el-tag></span>
              <span v-else><el-tag type="success" style="border:none;">已启用</el-tag></span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="操作" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button size="small" type="text" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button  class="green-button" size="small" type="text" @click="handleCopy(scope.row,$event)">复制配置</el-button>
              <el-button size="small" type="text" @click="changeStatus(scope.row)">
                <span v-if="scope.row.status==0">启用</span>
                <span v-else>停用</span>
              </el-button>
              <el-button  size="small" type="text" @click="delStrategy(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <!-- 策略面板 -->
      <el-dialog :title="dialogTitle" :visible.sync="bwlistVisible">
            <el-form  label-position="left" label-width="120px" style='width: 650px; margin-left:50px;' size="mini">
                <div class="form-block">
                    <span>基本设置</span>
                    <el-form-item label="策略名称" style="margin-top:10px">
                        <el-input  placeholder="请输入名称,不能重复,支持汉字等字符" style="width:300px" v-model="tempStrategy.name">
                        </el-input>
                    </el-form-item>
                    <el-form-item label="策略类型" style="margin-top:10px">
                        <el-radio-group v-model="tempStrategy.type" @change="selType" :disabled="!isCretaeSubmit">
                            <el-radio :label="1">黑白名单</el-radio>
                            <el-radio :label="2">超时重试</el-radio>
                            <el-radio :label="3">流量控制</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </div>
                <div class="form-block" v-if="tempStrategy.type==1">
                    <span>黑白名单</span>
                    <el-form-item label="名单类型" style="margin-top:10px">
                        <el-radio-group v-model="tempStrategy.sub_type" >
                            <el-radio :label="1">黑名单</el-radio>
                            <el-radio :label="2">白名单</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-alert
                        title="提示：若名单较长，可以边输入边提交，避免输入全部丢失"
                        :closable=false
                        type="warning">
                    </el-alert>
                    <el-form-item label="添加名单" v-if="tempStrategy.sub_type==1 || tempStrategy.sub_type==2" style="margin-top:10px">
                        <el-radio-group v-model="selbwType" @change="selBwType">
                            <el-radio :label="1">IP</el-radio>
                            <el-radio :label="2">参数</el-radio>
                        </el-radio-group>
                        <div>
                            <el-input  placeholder="Key" style="width:80px" v-model="tempBWKey" :disabled="selbwType==1">
                            </el-input>
                            <el-input  placeholder="Val" style="width:160px" v-model="tempBWVal">
                            </el-input>
                            <el-button @click="addBW" type="success" size="mini" icon="el-icon-plus" circle style="margin-left:10px"></el-button>
                        </div>
                    </el-form-item>
                    <el-form-item label="当前名单" v-if="tempStrategy.sub_type==1 || tempStrategy.sub_type==2">
                        <div v-for="kv in tempStrategy.content"  :key="kv.key" style="margin-top:-2px">
                            <el-tag type="success" size="large" style="width:70px;border:none;">{{kv.key}}</el-tag>
                            <el-tag type="info" size="large" style="width:200px;border:none">{{kv.val}}</el-tag>
                            <el-button @click.native="removeBW(kv)" type="text" size="mini" icon="el-icon-minus" style="margin-left:10px"></el-button>
                        </div>
                    </el-form-item>
                </div>
                <div class="form-block" v-if="tempStrategy.type==2">
                    <span>超时重试</span>
                    <el-form-item label="超时设置(秒)" style="width:200px;margin-top:10px" >
                        <el-tooltip content="取值必须在0和60之间，默认为15" placement="top">
                            <el-input-number v-model="tempStrategy.content.req_timeout"  :min="1" :max="60" ></el-input-number>
                        </el-tooltip>
                      </el-form-item>
                      <el-form-item label="重试次数" style="width:200px">
                        <el-tooltip content="发生网络错误时，自动进行请求重试,取值在0和5之间" placement="top">
                            <el-input-number v-model="tempStrategy.content.retry_times"  :min="0" :max="5" ></el-input-number>
                        </el-tooltip>
                      </el-form-item>
                      <el-form-item label="重试间隔(秒)" style="width:200px" >
                        <el-tooltip content="每次重试的间隔时间，取值在0和30之间" placement="top">
                            <el-input-number v-model="tempStrategy.content.retry_interval"  :min="1" :max="30" ></el-input-number>
                        </el-tooltip>
                    </el-form-item>
                </div>
                <div v-if="tempStrategy.type==3">
                  <div class="form-block" >
                      <span>接口流量</span>
                      <el-form-item label="QPS" style="margin-top:10px" >
                          <el-input-number v-model="tempStrategy.content.qps"  :min="0" :max="10000" ></el-input-number>
                          <el-alert
                              title="该API接口允许的每秒最大请求数,0代表不限制"
                              :closable=false
                              type="success">
                          </el-alert>
                        </el-form-item>
                        <el-form-item label="并发数">
                          <el-input-number v-model="tempStrategy.content.concurrent"  :min="0" :max="10000" ></el-input-number>
                          <el-alert
                              title="该API接口允许的每秒最大并发数,0代表不限制"
                              :closable=false
                              type="success">
                          </el-alert>
                        </el-form-item>
                  </div>
                  <div class="form-block" >
                      <span>用户流量</span>
                      <el-form-item label="参数名" style="margin-top:10px;width:400px" >
                        <el-input v-model="tempStrategy.content.param" placeholder="参数名为空，则不限制"></el-input>
                      </el-form-item>
                      <el-form-item label="限制时间">
                        <el-input-number v-model="tempStrategy.content.span"  :min="2" :max="2592000" ></el-input-number>秒
                      </el-form-item>
                      <el-form-item label="限制次数">
                        <el-input-number v-model="tempStrategy.content.times"  :min="1" :max="10240000" ></el-input-number>
                        <el-alert
                          title="在限制时间内，参数名对应的每个值，只允许访问限制的次数"
                          :closable=false
                          type="success">
                        </el-alert>
                        <el-alert
                            title="例如mobile参数名对应的所有手机号，在86400秒之内，只允许访问该接口1000次"
                            :closable=false
                            type="success">
                        </el-alert>
                        </el-form-item>
                  </div>
                  <!-- <div class="form-block" >
                      <span>熔断设置</span>
                      <el-form-item label="熔断错误率" style="margin-top:10px;" >
                        <el-input-number v-model="tempStrategy.content.fuse_error"  :min="0" :max="100" ></el-input-number>%
                        <el-alert
                          title="60秒内，超过一定错误率后，触发熔断，0表示不启用熔断"
                          :closable=false
                          type="success">
                        </el-alert>
                      </el-form-item>
                      <el-form-item label="熔断请求次数" style="margin-top:10px" >
                        <el-input-number v-model="tempStrategy.content.fuse_error_count"  :min="10" :max="100000" ></el-input-number>
                        <el-alert
                          title="60秒内，请求超过一定的次数，且满足上面的错误率，才能触发熔断"
                          :closable=false
                          type="success">
                        </el-alert>
                      </el-form-item>
                      <el-form-item label="熔断后百分比" style="margin-top:10px;" >
                        <el-input-number v-model="tempStrategy.content.fuse_percent"  :min="1" :max="100" ></el-input-number>%
                        <el-alert
                          title="触发熔断后，只有一定比例的请求能继续执行"
                          :closable=false
                          type="success">
                        </el-alert>
                      </el-form-item>
                      <el-form-item label="恢复错误率" style="margin-top:10px;" >
                        <el-input-number v-model="tempStrategy.content.fuse_recover"  :min="1" :max="99" ></el-input-number>%
                        <el-alert
                          title="触发熔断后，当过去60秒错误率小于该值后，取消熔断"
                          :closable=false
                          type="success">
                        </el-alert>
                      </el-form-item>
                      <el-form-item label="恢复请求次数" style="margin-top:10px;" >
                        <el-input-number v-model="tempStrategy.content.fuse_recover_count"  :min="1" :max="99" ></el-input-number>
                        <el-alert
                          title="60秒内，请求超过一定的次数，且小于上面的错误率，才能取消熔断"
                          :closable=false
                          type="success">
                        </el-alert>
                      </el-form-item>
                  </div> -->
                </div>
            </el-form>
 
        <div slot="footer" class="dialog-footer">
          <el-button @click="bwlistVisible=false">取消</el-button>
          <el-button  type="primary" v-if="isCretaeSubmit" @click="submitStrategy">
            <span>提交</span>
          </el-button>
          <el-button  type="primary" v-else @click="submitEdit">
            <span>提交</span>
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
     handleCopy(api, event) {
      var text = JSON.stringify(api)
      clip(text, event)
      this.$message({
        message: '复制成功',
        type: 'success',
        duration: 2000
      })
    },
    changeStatus(s) {
      var str = ''
      if (s.status == 0) {
        str = '您将要启用以下策略："' + s.name + '"！ 启用后，所有关联该策略得api将启用使用该策略！'
      } else {
         str = '您将要停用以下策略："' + s.name + '"！ 停用后，所有关联该策略得api将停止使用该策略！'
      }
      this.$confirm( str, '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var params = {
              target_service: 'tfe.manage',
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
              message: '策略状态变更成功',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消删除'
          });          
        });
    },
    delStrategy(s) {　
      this.$confirm('您将要删除以下策略："' + s.name + '", 注意删除后，所有关联该策略得api将自动解除关联' , '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var params = {
              target_service: 'tfe.manage',
              target_path: '/manage/strategy/delete',
              id: s.id,
              name: s.name,
              service: s.service,
              type: s.type
          }
          proxy('POST',params).then(res => {
            this.loadStrategy(this.selectedService)
            this.$message({
              message: '策略删除成功',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消删除'
          });          
        });
    },
    apiSet(s) {
      this.apiShow = true
      this.tempStrategy = s
    },
    submitEdit() {
      if (this.tempStrategy.name == '') {
        this.$message({
          message: '策略名不能为空',
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
        target_service: 'tfe.manage',
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
          message: '提交成功',
          type: 'success',
          duration: 3 * 1000,
          center: true
        })
      })
    },
    handleEdit(s) {
      this.dialogStatus = 'edit'
      this.dialogTitle = '编辑策略'
      this.tempStrategy = Object.assign({}, s) 
      this.tempStrategy.content = JSON.parse(this.tempStrategy.content)
      this.bwlistVisible = true
      this.tempBWKey = 'ip'
      this.isCretaeSubmit = false
    },
    handleSelService(service) {
      this.loadStrategy(service)
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
            message: 'key、val不能为空',
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
                message: '该名单项已存在',
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
            message: '策略名不能为空',
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
        target_service: 'tfe.manage',
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
          message: '提交成功',
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
            default:
                break;
        }
    },
    handleCreate() {
      if (this.selectedService == '') {
        this.$message({
            message: '请先选择一个Service',
            type: 'warning',
            duration: 3 * 1000,
            center: true
        })
        return 
      }
      this.dialogStatus = 'create'
      this.dialogTitle = '创建策略'
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
          target_service: 'tfe.manage',
          target_path: '/manage/strategy/load',
          service: this.selectedService,
          type: 0,
        }
        proxy('POST',params).then(res => {
          this.strategies = res.data.data
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
            this.loadStrategy(this.selectedService)
          }
        }
      })
    }
  },
  created() {
    this.loadServices()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
