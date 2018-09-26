<template>
  <div class="api-define">
    <div>
        <el-alert
        :closable="false"
        :title="$t('juz.basic')"
        type="warning">
        </el-alert>
    </div>
    <el-form  :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item :label="$t('common.service')" class="first-item">
            {{api.service}}
        </el-form-item>
        <el-form-item :label="$t('common.application')"  class="first-item">
            {{api.app}}
        </el-form-item>
        <el-form-item label="API ID"  class="first-item">
            {{api.api_id}}
        </el-form-item>
        <el-form-item :label="$t('juz.apiName')"  class="first-item">
            {{showName(api.api_id)}}
        </el-form-item>
        <el-form-item :label="$t('common.version')" class="first-item">
            {{showVersion(api.api_id)}}
        </el-form-item>
        <el-form-item :label="$t('common.desc')">
            {{api.desc}}
        </el-form-item>
        <el-form-item :label="$t('common.status')">
            <div><el-tag type="success">{{$t('common.edit')}}</el-tag> <span style="font-size:13px;">{{api.revise_version}}</span></div>
            <div v-if="api.release_version!=''" style="margin-top:3px"><el-tag :type="calcRelease(api)">{{$t('common.release')}}</el-tag>  <span style="font-size:13px;">{{api.release_version}}</span></div>
            <div v-else style="margin-top:3px"><el-tag type="warning">{{$t('common.notReleased')}}</el-tag></div>
        </el-form-item>
    </el-form>
    <div>
        <el-alert
        :closable="false"
        :title="$t('juz.request')"
        type="success">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <el-form-item :label="$t('juz.backendType')" class="first-item">
        <span v-if="api.backend_type==1">HTTP(s)</span>
        <span v-if="api.backend_type==2">Mock</span>
        </el-form-item>
        <div v-if="api.backend_type==1">
            <el-form-item :label="$t('juz.addrType')" >
                <span v-if="api.addr_type==1">URL</span>
                <span v-else>ETCD</span>
            </el-form-item>
            <el-form-item :label="$t('juz.backendAddr')">
                {{api.backend_addr}}
            </el-form-item>
            <el-form-item :label="$t('juz.backendURI')" v-if="api.addr_type==2">
                {{api.backend_uri}}
            </el-form-item>
            <el-form-item label="Header,Cookie">
                <el-tag type="success" size="large" style="border:none;">Pass through</el-tag>
            </el-form-item>
            <el-form-item :label="$t('juz.params')">
                <el-tag type="success" size="large" style="border:none;">Pass through</el-tag>
            </el-form-item>
            </div>
            <div v-else>
            <el-form-item label="Mock Data">
                {{api.mock_data}}
            </el-form-item>    
        </div>
    </el-form>
    <div>
    <el-alert
        :closable="false"
        :title="$t('juz.advance')"
        type="info">
        </el-alert>
    </div>
    <el-form   :class="{'black-form':isBlackForm()}" label-position="left" label-width="150px" size="small">
        <div class="form-block" style="margin-top:10px">
            <span>{{$t('juz.paramVerify')}}</span>
            <el-form-item :label="$t('common.status')" style="width:300px" class="first-item">
                <span v-if="api.verify_on==1">On</span>
                <span v-else>Off</span>
            </el-form-item>
            <el-form-item :label="$t('juz.verifyTable')" style="width:300px" class="first-item">
                <el-popover
                    placement="top"
                    width="800"
                    trigger="click">
                    <paramRules :rules="rules"></paramRules>
                    <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getParamRules()"></i> 
                </el-popover>
            </el-form-item>
        </div>
        <div class="form-block">
            <span>{{$t('juz.canaryTest')}}</span>
            <el-form-item :label="$t('common.status')" style="width:300px" class="first-item">
                <span v-if="api.traffic_on==1">On</span>
                <span v-else>Off</span>
            </el-form-item>
            <el-form-item :label="$t('juz.targetAPIID')" prop="req_timeout" style="width:500px" class="first-item">
                {{api.traffic_api}}
            </el-form-item>
            <el-form-item :label="$t('juz.trafficRatio')" prop="req_timeout">
                {{api.traffic_ratio}}
            </el-form-item>
            <el-form-item :label="$t('juz.ipList')" style="width:300px" class="first-item">
                {{api.traffic_ips}}
            </el-form-item>
        </div>

        <div class="form-block">
            <span>{{$t('juz.bwList')}}</span>
            <el-form-item v-if="api.bw_strategy!=0" label="Strategy"  class="first-item">
                {{api.bw_strategy}} 
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <bwlist :strategy="strategy"></bwlist>
                    <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(api.bw_strategy)"></i> 
                </el-popover>
                
            </el-form-item>
             <el-form-item v-else :label="$t('juz.noStrategy')"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>{{$t('juz.timoutRetry')}}</span>
            <el-form-item v-if="api.retry_strategy!=0" label="Strategy"  class="first-item">
                {{api.retry_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <retry :strategy="strategy"></retry>
                    <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px " @click="getStrategy(api.retry_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else :label="$t('juz.noStrategy')"  class="first-item">
            </el-form-item>
        </div>
        <div class="form-block">
            <span>{{$t('juz.trafficControl')}}</span>
            <el-form-item v-if="api.traffic_strategy!=0" label="Strategy"  class="first-item">
                {{api.traffic_strategy}}
                <el-popover
                    placement="top"
                    width="600"
                    trigger="click">
                    <traffic :strategy="strategy"></traffic>
                    <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(api.traffic_strategy)"></i> 
                </el-popover>
            </el-form-item>
            <el-form-item v-else :label="$t('juz.noStrategy')"  class="first-item">
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