<template>
  <div class="api-define">
    <div>
        <el-alert
        :closable="false"
        title="基本信息"
        type="warning">
        </el-alert>
    </div>
    <el-form  :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item label="Service" class="first-item">
            {{api.service}}
        </el-form-item>
        <el-form-item label="API ID"  class="first-item">
            {{api.api_id}}
        </el-form-item>
        <el-form-item label="API名"  class="first-item">
            {{showName(api.api_id)}}
        </el-form-item>
        <el-form-item label="版本号" class="first-item">
            {{showVersion(api.api_id)}}
        </el-form-item>
        <el-form-item label="描述">
            {{api.desc}}
        </el-form-item>
        <el-form-item label="发布状态">
            <div><el-tag type="success">编辑版本</el-tag> <span style="font-size:13px;">{{api.revise_version}}</span></div>
            <div v-if="api.release_version!=''" style="margin-top:3px"><el-tag :type="calcRelease(api)">发布版本</el-tag>  <span style="font-size:13px;">{{api.release_version}}</span></div>
            <div v-else style="margin-top:3px"><el-tag type="warning">尚未发布</el-tag></div>
        </el-form-item>
    </el-form>
    <div>
        <el-alert
        :closable="false"
        title="基本设置"
        type="success">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item label="后端类型" class="first-item">
        <span v-if="api.route_proto==1">HTTP(s)服务</span>
        <span v-if="api.route_proto==2">Mock</span>
        </el-form-item>
        <div v-if="api.route_proto==1">
        <el-form-item label="后端地址">
            {{api.route_addr}}
        </el-form-item>
        <el-form-item label="Header">
            <el-tag type="success" size="large" style="border:none;">默认透传</el-tag>
        </el-form-item>
        <el-form-item label="请求参数">
            <el-tag type="success" size="large" style="border:none;">默认透传</el-tag>
        </el-form-item>
        </div>
        <div v-else>
        <el-form-item label="Mock返回">
            <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder="这里请填写数据的元信息，例如返回json就填写json格式,返回字符串就填写字符串，网关不做任何转换，直接返回给请求的客户端" v-model="api.mock_data"></el-input>
            <el-alert
                title="注意！设置Mock后，网关会自动模拟数据，实际请求不会到达后端服务"
                type="warning">
            </el-alert>
        </el-form-item>    
        </div>
    </el-form>
    <div>
    <el-alert
        :closable="false"
        title="高级设置"
        type="info">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <div class="form-block" style="margin-top:10px">
            <span>参数验证</span>
            <el-form-item label="开启验证" style="width:300px" class="first-item">
                <span v-if="api.verify_on==1">开启</span>
                <span v-else>关闭</span>
            </el-form-item>
            <el-form-item label="验证表" style="width:300px" class="first-item">
                <el-popover
                    placement="top"
                    width="800"
                    trigger="click">
                    <paramRules :rules="rules"></paramRules>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getParamRules()"></i> 
                </el-popover>
            </el-form-item>
        </div>
        <div class="form-block">
            <span>流量路由(金丝雀测试)</span>
            <el-form-item label="开启路由" style="width:300px" class="first-item">
                <span v-if="api.traffic_on==1">开启</span>
                <span v-else>关闭</span>
            </el-form-item>
            <el-form-item label="目标api" prop="req_timeout" style="width:500px" class="first-item">
                {{api.traffic_api}}
            </el-form-item>
            <el-form-item label="路由比例" prop="req_timeout">
                {{api.traffic_ratio}}
            </el-form-item>
            <el-form-item label="IP列表" style="width:300px" class="first-item">
                {{api.traffic_ips}}
            </el-form-item>
        </div>

        <div class="form-block">
            <span>黑白名单</span>
            <el-form-item v-if="api.bw_strategy!=0" label="已选策略"  class="first-item">
                {{api.bw_strategy}} 
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <bwlist :strategy="strategy"></bwlist>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.bw_strategy)"></i> 
                </el-popover>
                
            </el-form-item>
             <el-form-item v-else label="尚未选择策略"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>超时重试</span>
            <el-form-item v-if="api.retry_strategy!=0" label="已选策略"  class="first-item">
                {{api.retry_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <retry :strategy="strategy"></retry>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.retry_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else label="尚未选择策略"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>流量控制</span>
            <el-form-item v-if="api.traffic_strategy!=0" label="已选策略"  class="first-item">
                {{api.traffic_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <traffic :strategy="strategy"></traffic>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.traffic_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else label="尚未选择策略"  class="first-item">
            </el-form-item>
        </div>
    </el-form>
</div>
</template>

<script>
/* eslint-disable */
import {proxy} from '@/api/apiProxy'
import bwlist from '../components/bwlist'
import retry from '../components/retry'
import traffic from '../components/traffic'
import paramRules from '../components/paramRules'
export default {
  name: 'apiDefine',
   props: {
    api: {
      type: Object,
      default: {}
    },
    formType: {
        type: String,
        default: 'white-form'
    }
  },
  components:{bwlist,retry,paramRules,traffic},
  data() {
      return {
          strategies: [],
          strategy: {
              content:{}
          },
          rules: []
      }
  },
  methods: {
    showName(apiID) {
        return apiID.substring(0,apiID.length-3)
    },
    showVersion(apiID) {
        return apiID[apiID.length-1]
    },
    getParamRules() {
        this.rules = JSON.parse(this.api.param_rules)
    },
    getStrategy(id) {
        var params = {
          target_service: 'tfe.manage',
          target_path: '/manage/strategy/query',
          id: id
        }
      proxy('POST',params).then(res => {
        this.strategy = res.data.data
        this.strategy.content = JSON.parse(this.strategy.content)
      })
    },
    calcRelease(api) {
      if (api.release_version == api.revise_version) {
        return 'success'
      } 
      return 'warning'
    },
    isBlackForm() {
        return this.formType== 'black-form'
    }
  },
  mounted() {

  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>

<style lang="scss" scoped>


</style>

<style lang="scss">
.el-dialog__body {
    padding-top :0 !important
}
</style>