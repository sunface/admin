<template>
  <div class="api-define">
    <div>
        <el-alert
        :closable="false"
        title="Basic"
        type="warning">
        </el-alert>
    </div>
    <el-form  :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item label="Service" class="first-item">
            {{api.service}}
        </el-form-item>
        <el-form-item label="App" class="first-item">
            {{api.app}}
        </el-form-item>
        <el-form-item label="API ID"  class="first-item">
            {{api.api_id}}
        </el-form-item>
        <el-form-item label="API Name"  class="first-item">
            {{showName(api.api_id)}}
        </el-form-item>
        <el-form-item label="Version" class="first-item">
            {{showVersion(api.api_id)}}
        </el-form-item>
        <el-form-item label="Desc">
            {{api.desc}}
        </el-form-item>
        <el-form-item label="Status">
            <div><el-tag type="success">editVer</el-tag> <span style="font-size:13px;">{{api.revise_version}}</span></div>
            <div v-if="api.release_version!=''" style="margin-top:3px"><el-tag :type="calcRelease(api)">releaseVer</el-tag>  <span style="font-size:13px;">{{api.release_version}}</span></div>
            <div v-else style="margin-top:3px"><el-tag type="warning">not released</el-tag></div>
        </el-form-item>
    </el-form>
    <div>
        <el-alert
        :closable="false"
        title="Request"
        type="success">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item label="Backend Type" class="first-item">
        <span v-if="api.backend_type==1">HTTP(s)</span>
        <span v-if="api.backend_type==2">Mock</span>
        </el-form-item>
        <div v-if="api.backend_type==1">
            <el-form-item label="Addr Type">
                <span v-if="api.addr_type==1">URL</span>
                <span v-else>ETCD</span>
            </el-form-item>
            <el-form-item label="Backend Addr">
                {{api.backend_addr}}
            </el-form-item>
            <el-form-item label="Backend Uri" v-if="api.addr_type==2">
                {{api.backend_uri}}
            </el-form-item>
            <el-form-item label="Header,Cookie">
                <el-tag type="success" size="large" style="border:none;">Pass through</el-tag>
            </el-form-item>
            <el-form-item label="Params">
                <el-tag type="success" size="large" style="border:none;">Pass through</el-tag>
            </el-form-item>
            </div>
            <div v-else>
            <el-form-item label="Mock Data">
                <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}"  v-model="api.mock_data"></el-input>
            </el-form-item>    
        </div>
    </el-form>
    <div>
    <el-alert
        :closable="false"
        title="Advance"
        type="info">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <div class="form-block" style="margin-top:10px">
            <span>Param Verify</span>
            <el-form-item label="Status" style="width:300px" class="first-item">
                <span v-if="api.verify_on==1">On</span>
                <span v-else>Off</span>
            </el-form-item>
            <el-form-item label="Verify Table" style="width:300px" class="first-item">
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
            <span>Canary Test</span>
            <el-form-item label="Status" style="width:300px" class="first-item">
                <span v-if="api.traffic_on==1">On</span>
                <span v-else>Off</span>
            </el-form-item>
            <el-form-item label="Target APIID" prop="req_timeout" style="width:500px" class="first-item">
                {{api.traffic_api}}
            </el-form-item>
            <el-form-item label="Traffic Ratio" prop="req_timeout">
                {{api.traffic_ratio}}
            </el-form-item>
            <el-form-item label="IP List" style="width:300px" class="first-item">
                {{api.traffic_ips}}
            </el-form-item>
        </div>

        <div class="form-block">
            <span>White/Black List</span>
            <el-form-item v-if="api.bw_strategy!=0" label="Strategy"  class="first-item">
                {{api.bw_strategy}} 
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <bwlist :strategy="strategy"></bwlist>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.bw_strategy)"></i> 
                </el-popover>
                
            </el-form-item>
             <el-form-item v-else label="no strategy"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>Timeout/Retry</span>
            <el-form-item v-if="api.retry_strategy!=0" label="Strategy"  class="first-item">
                {{api.retry_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <retry :strategy="strategy"></retry>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.retry_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else label="no strategy"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>Traffic Control</span>
            <el-form-item v-if="api.traffic_strategy!=0" label="Strategy"  class="first-item">
                {{api.traffic_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <traffic :strategy="strategy"></traffic>
                    <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(api.traffic_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else label="no strategy"  class="first-item">
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
          target_app: 'juzManage',
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