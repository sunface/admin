<template>
  <div class="app-container">
      <div class="filter-container">
        <el-tag> {{$t('common.service')}}</el-tag>
        <el-select clearable class="filter-item" :value="calcService()" @change='handleSelService' style="width: 160px"  :placeholder="$t('common.selectService')">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"  v-waves type="success" icon="el-icon-edit">{{$t('juz.createApi')}}</el-button>

       

        <span style="float:right;margin-right:30px">
          <el-button class="filter-item"  @click="batchRelease"  v-waves icon="el-icon-upload2">{{$t('juz.batchReleas')}}</el-button>
          <el-button class="filter-item" style="margin-right: 20px;" @click="handleBatchStrategy"  v-waves icon="el-icon-bell">{{$t('juz.batchStrategy')}}</el-button>
          <el-input  style="width: 250px;" class="filter-item" :placeholder="$t('common.fuzzingSearch')" v-model="q" clearable @clear="loadApis(1)" @keyup.enter.native="loadApis(1)" >
          </el-input>
          <el-button class="filter-item" type="text" v-waves icon="el-icon-search"   @click="loadApis(1)"></el-button>
        </span>

      </div>

      <div class="table">
        <el-table  :data="apis"  fit highlight-current-row style="width: 100%;"  :default-sort = "{prop: 'modify_date', order: 'descending'}" @selection-change="selTableRows">
           <el-table-column
            type="selection"
            width="55"> 
          </el-table-column>
          <el-table-column align="left" label="api id" width="210" prop="api_id">
            <template slot-scope="scope">
              <span>{{scope.row.api_id}}</span>
            </template>
          </el-table-column>
          <el-table-column align="left" :label="$t('common.application')" width="160">
            <template slot-scope="scope">
              <span>{{scope.row.app}}</span>
            </template>
          </el-table-column>
          <el-table-column width="400" align="left" :label="$t('juz.backendInfo')" >
            <template slot-scope="scope">
              <span v-if="scope.row.backend_type==1">HTTP(S)-</span>
              <span v-else>Mock-</span>
              <span v-if="scope.row.addr_type==1">URL-</span>
              <span v-else>ETCD-</span>
              {{scope.row.backend_addr}}
              <span v-if="scope.row.addr_type==2">-{{scope.row.backend_uri}}</span>
            </template>
          </el-table-column>
          <el-table-column width="250" align="left" :label="$t('common.status')" >
             <template slot-scope="scope">
                <div><el-tag type="success" size="mini" >{{$t('common.edit')}}</el-tag> <span style="color:#337ab7;font-size:13px;">{{scope.row.revise_version}}</span></div>
                <div v-if="scope.row.release_version!=''" style="margin-top:3px"><el-tag  size="mini"   :type="calcRelease(scope.row)">{{$t('common.release')}}</el-tag>  <span style="color:#337ab7;font-size:13px;">{{scope.row.release_version}}</span></div>
                <div v-else style="margin-top:3px"><el-tag type="warning" size="mini">{{$t('common.notReleased')}}</el-tag></div>
            </template>
          </el-table-column>
          <el-table-column align="center" :label="$t('common.operate')" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <span class="table-op-btn" v-if="scope.row.release_version!=scope.row.revise_version && scope.row.release_version!=''" @click="releaseAPI(scope.row)">{{$t('common.release')}}</span>
              <span  class="table-op-btn" v-if="scope.row.release_version=='' " @click="releaseAPI(scope.row)">{{$t('common.release')}}</span>
              <span  class="table-op-btn" @click="handleManage(scope.row)">{{$t('common.manage')}}</span>
              <span   class="table-op-btn" @click="handleCopy(scope.row,$event)">{{$t('common.copyConfig')}}</span>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
            layout="total,prev, pager, next"
            @current-change="loadApis"
            :total="totalApis">
        </el-pagination>
      </div>
      <el-dialog class="black-dialog" :title="'Manage API ：'+tempApi.api_id" :visible.sync="apiManageVisible" top="40px"  :fullscreen=true style="backgroud:rgb(57, 79, 90) !important"> 
          <el-row :gutter="20" style="margin-top:40px;">
            <el-col :span=3 :offset=1>
              <el-menu
                default-active="1"
                class="el-menu-vertical-demo"
                @select="onManageSelect"
                text-color="#fff"
                background-color="rgb(57, 79, 90)"
                active-text-color="#ffd04b"
                ref="apiMenu"
                >
                <el-menu-item index="1">
                   <i class="el-icon-setting"></i> 
                  <span slot="title">{{$t('juz.apiDefine')}}</span>
                </el-menu-item>
                <el-menu-item index="2">
                  <i class="el-icon-caret-right"></i>
                  <span slot="title">{{$t('juz.apiTest')}}</span>
                </el-menu-item>
                <el-menu-item index="3">
                   <i class="el-icon-menu"></i>
                  <span slot="title">{{$t('juz.auditLog')}}</span>
                </el-menu-item>
                <el-menu-item index="4">
                   <i class="el-icon-bell"></i>
                  <span slot="title">{{$t('juz.metrics')}}</span>
                </el-menu-item>
                <!-- <el-menu-item index="5">
                   <i class="el-icon-document"></i>
                  <span slot="title">API文档</span>
                </el-menu-item> -->
              </el-menu>
            </el-col>
            <el-col :span=20 :offset=0 style="color:#fff;padding-right:0;">
              <!-- api定义 -->
                <div v-show="selectManageIndex==1">
                     <div style="width:85%">
                        <span style="float:right;margin-top:-36px">
                          <el-button type="success" size="small" @click="handleEdit()">{{$t('common.edit')}}</el-button>
                          <el-button v-if="tempApi.release_version!=tempApi.revise_version && tempApi.release_version!=''" type="warning" size="small" @click="releaseAPI(tempApi)">{{$t('common.release')}}</el-button>
                          <el-button v-if="tempApi.release_version=='' " type="warning" size="small" @click="releaseAPI(tempApi)">{{$t('common.release')}}</el-button>
                         
                          <el-button v-show="tempApi.release_version!=''" type="info" size="small" @click="offlineAPI(tempApi)">{{$t('common.offline')}}</el-button>
                           <el-button  type="" @click="delApi(tempApi)" size="small">{{$t('common.delete')}}</el-button>
                        </span>
                        <apiDefine :api="tempApi" formType='black-form'></apiDefine>
                    </div>
                </div>
                <!-- api测试 -->
                <div v-show="selectManageIndex==2">  
                  <apiTest :api="tempApi"></apiTest>
                </div>

                <!-- 审计日志 -->
                <div v-show="selectManageIndex==3">
                   <auditLog :totalLogs="totalLogs" :logs="logs" themeType='black-theme' :targetType=2 :targetID="tempApi.api_id"></auditLog>
                </div>
                
                <!-- 监控统计 -->
                <div v-show="selectManageIndex==4">
                   <iframe :src="apiGrafanaAddr(tempApi.api_id)" style="height:910px;width:100%;border:none;margin-top:-40px;"></iframe>
                </div>
                <!--  -->
            </el-col>
          </el-row>
      </el-dialog>

      <el-dialog  class="mf-dialog" :title="apiTitle" :visible.sync="apiDefineVisible" top="40px" :before-close="cancelAPIDefine" :fullscreen=false> 
        <el-steps :active="apiDefineStep" finish-status="success">
          <el-step :title="$t('juz.basic')"></el-step>
          <el-step :title="$t('juz.request')"></el-step>
          <el-step :title="$t('juz.advance')"></el-step>
        </el-steps>

        <el-row :gutter="20" style="margin-top:40px;">
            <el-col :span=20 :offset=2>
              <div v-show="apiDefineStep==0">
                  <el-form label-position="left" label-width="110px" size="small">
                    <el-form-item :label="$t('common.service')">
                      {{tempApi.service}}
                    </el-form-item>
                    <el-form-item :label="$t('common.application')" required>
                      <el-select clearable class="filter-item" v-model="tempApi.app" style="width: 200px"  :placeholder="$t('common.select')">
                            <el-option v-for="s in  apps" :key="s.name" :label="s.name" :value="s.name">
                            </el-option>
                          </el-select>
                    </el-form-item>
                    <el-form-item label="URL Type">
                       <el-switch
                          v-model="tempApi.path_type"
                          active-color="#13ce66"
                          :active-value='switchActive'
                          :inactive-value='switchInactive'
                          :disabled="defineStatus=='edit'">
                        </el-switch>
                        <el-alert
                        :closable=false
                        v-if="tempApi.path_type==1"
                        :title="$t('juz.urlType1Tips')"
                        type="success">
                      </el-alert> 
                       <el-alert
                        :closable=false
                        v-else
                         :title="$t('juz.urlType2Tips')"
                        type="success">
                      </el-alert> 
                    </el-form-item>
                   
                    <el-form-item label="API ID">
                      <span v-if="defineStatus=='create'" >{{genAPIIDByType()}}</span>
                      <span v-else >{{tempApi.api_id}}</span>
                      <el-alert
                        :closable=false
                        :title="apiNameAlert()"
                        type="warning">
                      </el-alert>
                    </el-form-item>
                    <el-form-item v-if="defineStatus=='create'" :label="$t('juz.apiName')"  class="first-item" required>
                      <span v-if="tempApi.path_type==0">{{tempApi.service}}.</span>
                      <el-tooltip class="no-border-input" :content="apiNameTooltip()" placement="top">
                          <el-input  v-model="tempApi.api_id" placeholder=""  style="width:250px">
                          </el-input>
                       </el-tooltip>
                    </el-form-item>
                    <el-form-item v-else :label="$t('juz.apiName')" class="first-item" required>
                      {{showName(tempApi.api_id)}}
                    </el-form-item>
                    <el-form-item :label="$t('common.version')">
                      <el-input-number v-if="defineStatus=='create'" class="no-border-input" style="width:90px;border:none" v-model="tempApi.version"  :min="1" :max="9" size="mini"></el-input-number>
                      <span v-else>{{showVersion(tempApi.api_id)}}</span>
                    </el-form-item>
                    <el-form-item :label="$t('common.desc')">
                      <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder="Describe your api" v-model="tempApi.desc">
                      </el-input>
                    </el-form-item>
                  </el-form>
              </div>
              <div v-show="apiDefineStep==1">
                  <el-form label-position="left" label-width="150px" size="small">
                    <el-form-item :label="$t('juz.routeType')" prop="route_type">
                      <el-select  class="filter-item" v-model="tempApi.route_type" style="width: 150px"  placeholder="please choose">
                        <el-option :label="$t('juz.routeDirect')" :value=1></el-option>
                        <el-option   :label="$t('juz.routeRedirect')" :value=2></el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item :label="$t('juz.backendType')" prop="name" class="first-item">
                      <el-radio-group v-model="tempApi.backend_type">
                        <el-radio :label="1">HTTP(s)</el-radio>
                        <el-radio :label="2">Mock</el-radio>
                      </el-radio-group> 
                    </el-form-item>
                    <div v-if="tempApi.backend_type==1">
                      <el-form-item label="HTTP Method">
                          <el-select  class="filter-item" v-model="tempApi.method" style="width: 200px;margin-top:8px"  placeholder="select method">
                            <el-option label="GET" value="GET"></el-option>
                            <el-option  label="POST" value="POST"></el-option>
                            <el-option  label="PUT" value="PUT"></el-option>
                            <el-option  label="DELETE" value="DELETE"></el-option>
                          </el-select>
                      </el-form-item>
                       <el-form-item :label="$t('juz.addrType')" prop="name" class="first-item">
                        <el-radio-group v-model="tempApi.addr_type">
                          <el-radio :label="1">URL</el-radio>
                          <el-radio :label="2">ETCD</el-radio>
                        </el-radio-group> 
                      </el-form-item>
                      <div v-if="tempApi.addr_type==1">
                        <el-form-item :label="$t('juz.backendAddr')" required>
                          <el-tooltip content="The application url you want to access" placement="top">
                              <el-input v-model="tempApi.backend_addr" placeholder="http://"></el-input>
                          </el-tooltip>
                        </el-form-item>
                      </div>
                      <div v-else>
                        <el-form-item :label="$t('juz.backendAddr')" required>
                          <el-select clearable class="filter-item" v-model="tempService" @change='handleTempService' style="width: 160px"  placeholder="select service">
                            <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
                            </el-option>
                          </el-select>
                          <el-select clearable class="filter-item" v-model="tempApi.backend_addr" style="width: 200px"  placeholder="select app">
                            <el-option v-for="s in  tempApps" :key="s.name" :label="s.name" :value="s.name">
                            </el-option>
                          </el-select>
                        </el-form-item>
                        <el-form-item :label="$t('juz.backendURI')" required>
                          <el-tooltip content="Http uri,e.g. if target is 'http://test.com/a/b',here you shold fill in '/a/b'" placement="top">
                              <el-input v-model="tempApi.backend_uri" placeholder="/a/b"></el-input>
                          </el-tooltip>
                        </el-form-item>
                      </div>
                    
                      <el-form-item label="Header,Cookie">
                          <el-tag type="success" size="large" style="border:none;">{{$t('juz.passThrough')}}</el-tag>
                      </el-form-item>
                       <el-form-item label="Params">
                          <el-tag type="success" size="large" style="border:none;">{{$t('juz.passThrough')}}</el-tag>
                      </el-form-item>
                    </div>
                    <div v-else>
                      <el-form-item label="Mock Data">
                        <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder='if you fill with {"a":"b"},juz will directly return the filling content' v-model="tempApi.mock_data"></el-input>
                          <el-alert
                          :closable=false
                            title="Attention! When mock is used, the client request will get the mock data, the request will never reach the backend service"
                            type="warning">
                          </el-alert>
                      </el-form-item>    
                    </div>
                  
                
                  </el-form>
              </div>
              <div v-show="apiDefineStep==2">
                  <el-form label-position="left" label-width="150px" size="small">
                      <div class="form-block">
                        <span>{{$t('juz.paramVerify')}}<el-tag type="success" size="mini" style="border:none;">optional</el-tag></span>
                        <el-form-item label="On/Off" style="width:300px" class="first-item">
                            <el-switch
                              v-model="tempApi.verify_on"
                              active-color="#13ce66"
                              :active-value='switchActive'
                              :inactive-value='switchInactive'>
                            </el-switch>
                        </el-form-item>
                        <el-form-item :label="$t('juz.verifyTable')" style="width:300px" class="first-item">
                          <el-button type="success" icon="el-icon-edit"  @click="showVerify">{{$t('common.edit')}}</el-button>
                        </el-form-item>
                      </div>

                      <div class="form-block">
                        <span>{{$t('juz.canaryTest')}}<el-tag type="success" size="mini" style="border:none;">optional</el-tag></span>
                        <el-form-item label="On/Off" style="width:300px" class="first-item">
                            <el-switch
                              v-model="tempApi.traffic_on"
                              active-color="#13ce66" 
                              :active-value='switchActive'
                              :inactive-value='switchInactive'>
                            </el-switch>
                        </el-form-item>
                        <el-form-item :label="$t('juz.targetAPIID')" style="width:500px" class="first-item">
                          <el-tooltip :content="$t('juz.targetIDTips')" placement="top">
                            <el-input v-model="tempApi.traffic_api" placeholder="">
                            </el-input>
                          </el-tooltip>
                        </el-form-item>
                        <el-form-item :label="$t('juz.trafficRatio')" prop="req_timeout">
                          <el-tooltip :content="$t('juz.trafficRatioTips')" placement="top">
                              <el-input-number v-model="tempApi.traffic_ratio" :min="1" :max="100"></el-input-number>
                          </el-tooltip>
                        </el-form-item>
                        <el-form-item :label="$t('juz.ipList')" style="width:400px" class="first-item">
                          <el-tooltip :content="$t('juz.ipListTips')" placement="top">
                            <el-input v-model="tempApi.traffic_ips" :placeholder="$t('juz.sepIpWithSpace')">
                            </el-input>
                          </el-tooltip>
                        </el-form-item>
                      </div>

                       <el-alert
                        style="margin-top:8px"
                            :title="$t('juz.createStrategyTips')"
                            :closable=false
                            type="success">
                          </el-alert>
                      <div class="form-block">
                        <span>{{$t('juz.bwList')}}<el-tag type="success" size="mini" style="border:none;">optional</el-tag></span>
                        <el-form-item :label="$t('juz.strategy')" prop="req_timeout" class="first-item">
                          <el-select clearable class="filter-item" v-model="selBW" @change="selectBW" style="width: 200px"  placeholder="please choose">
                            <el-option v-for="s in  filterStrategies(1)" :key="s.name" :label="s.name" :value="s.id">
                            </el-option>
                          </el-select> 
                          <el-popover
                            v-show="this.selBW != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <bwlist :strategy="strategy"></bwlist>
                            <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(selBW)"></i> 
                          </el-popover>
                        </el-form-item>
                      </div>
                      <div class="form-block">
                        <span>{{$t('juz.timoutRetry')}}<el-tag type="success" size="mini" style="border:none;">optional</el-tag></span>
                        <el-form-item :label="$t('juz.strategy')" prop="req_timeout" class="first-item">
                          <el-select clearable class="filter-item" v-model="selRetry" @change="selectRetry" style="width: 200px"  placeholder="please choose">
                            <el-option v-for="s in  filterStrategies(2)" :key="s.name" :label="s.name" :value="s.id"> 
                            </el-option>
                          </el-select>
                          <el-popover
                            v-show="this.selRetry != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <retry :strategy="strategy"></retry>
                            <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(selRetry)"></i> 
                        </el-popover>
                        </el-form-item>
                      </div>
                      <div class="form-block">
                        <span>{{$t('juz.trafficControl')}}<el-tag type="success" size="mini" style="border:none;">optional</el-tag></span>
                        <el-form-item :label="$t('juz.strategy')" prop="req_timeout" class="first-item">
                          <el-select clearable class="filter-item" v-model="selTraffic" @change="selectTraffic" style="width: 200px"  placeholder="please choose">
                            <el-option v-for="s in  filterStrategies(3)" :key="s.name" :label="s.name" :value="s.id"> 
                            </el-option>
                          </el-select>
                          <el-popover
                            v-show="this.selTraffic != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <traffic :strategy="strategy"></traffic>
                            <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px;icon-primary" @click="getStrategy(selTraffic)"></i> 
                        </el-popover>
                        </el-form-item>
                      </div>
                  </el-form>
              </div>
            </el-col>
        </el-row>

         <div slot="footer" class="dialog-footer">
            <el-button @click="nextDefineStep(false)" v-show="apiDefineStep!=0">Prev</el-button>
            <el-button @click="nextDefineStep(true)" v-show="apiDefineStep!=2" type="success">Next</el-button>
            <el-button  type="primary" @click="saveApi" v-show="apiDefineStep!=0">Save</el-button>
          </div>
      </el-dialog>


      <el-dialog
        class="mf-dialog" 
        :title="$t('juz.paramVerify')"
        :fullscreen=true
        :visible.sync="veryfyVisible"
      >
        <el-row :gutter="20" style="margin-top:0px;">
            <el-col :xs="{span:24}" :sm="{span:18}" :md="{span: 8}" :lg="{ span: 6}" style="border-right:.5px solid #999">
                <div v-html="ruleContent" ></div>
            </el-col>
            <el-col :span=16 :offset=1>
                <el-form label-position="left" label-width="120px">
                  <div>
                    <span>{{$t('juz.createRule')}}</span>
                    <el-form-item :label="$t('common.param')" style="margin-top:15px">
                      <el-input v-model="tempVerifyRule.param" placeholder="Param name" style="width:150px" >
                      </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('common.rule')">
                      <el-input v-model="tempVerifyRule.rule" placeholder="regexp format" style="width:350px">
                      </el-input>
                    </el-form-item>
                    <el-form-item :label="$t('juz.testData')">
                      <el-input v-model="tempVerifyRule.test_data" placeholder="input test data for test the rule" style="width:150px">
                      </el-input>
                    </el-form-item>
                     <el-button @click="addVerifyRule" type="success" >{{$t('juz.testCreate')}}</el-button>
                  </div>
              
                </el-form>

                <div class="table" style="margin-top:15px;overflow-y:auto">
                  <el-table  :data="tempApi.param_rules" border fit highlight-current-row style="width: 100%;">
                    <el-table-column align="center" :label="$t('common.param')" width="150">
                      <template slot-scope="scope">
                        <span>{{scope.row.param}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column width="400" align="center" :label="$t('common.rule')">
                      <template slot-scope="scope">
                        <span>{{scope.row.rule}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column width="250" align="center" :label="$t('juz.testData')" >
                      <template slot-scope="scope">
                          <span>{{scope.row.test_data}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column align="center" :label="$t('common.operate')" class-name="small-padding fixed-width">
                      <template slot-scope="scope"> 
                        <el-button  type="text" size="mini" icon="el-icon-minus" circle @click="delParamRule(scope.row.param)"></el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
            </el-col>
        </el-row>

        <span slot="footer" class="dialog-footer">
          <el-button @click="cancelVerify">Close</el-button>
        </span>
      </el-dialog>

      <!-- 批量设置策略 -->
      <el-dialog title="Batch Strategy Setting" class="mf-dialog"  :visible.sync="batchStrategyVisible">
        <el-form  label-position="left" label-width="120px" style='width: 650px; margin-left:50px;' size="mini">
            <div class="form-block">
              <span>Batch Strategy</span>
              <el-alert
                  style="margin-top:8px"
                  title="Warning: batch setting will overide the api's specific setting, please operate with caution!"
                  :closable=false
                  type="warning">
              </el-alert>
              <el-alert
                  style="margin-top:8px"
                  title="Tips: the unselected strategy will be ignored, it will not setting your api's corresponding strategy to null"
                  :closable=false
                  type="success">
              </el-alert>
                <div class="form-block">
                  <span>White/Black List</span>
                  <el-form-item label="Stategy">
                    <el-select clearable class="filter-item" v-model="batchBW" style="width: 200px"  placeholder="please choose">
                      <el-option v-for="s in  filterStrategies(1)" :key="s.name" :label="s.name" :value="s.id">
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchBW != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <bwlist :strategy="strategy"></bwlist>
                      <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(batchBW)"></i> 
                    </el-popover>
                  </el-form-item>
                </div>
                <div class="form-block">
                  <span>Timeout/Retry</span>
                  <el-form-item label="Strategy" prop="req_timeout">
                    <el-select clearable class="filter-item" v-model="batchRetry"  style="width: 200px"  placeholder="please choose">
                      <el-option v-for="s in  filterStrategies(2)" :key="s.name" :label="s.name" :value="s.id"> 
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchRetry != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <retry :strategy="strategy"></retry>
                      <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px;" @click="getStrategy(batchRetry)"></i> 
                  </el-popover>
                  </el-form-item>
                </div>
                <div class="form-block">
                  <span>Traffic Control</span>
                  <el-form-item label="选择策略" prop="req_timeout">
                    <el-select clearable class="filter-item" v-model="batchTraffic"  style="width: 200px"  placeholder="please choose">
                      <el-option v-for="s in  filterStrategies(3)" :key="s.name" :label="s.name" :value="s.id"> 
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchTraffic != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <traffic :strategy="strategy"></traffic>
                      <i slot="reference" class="el-icon-view hover-cursor icon-primary" style="margin-left: 8px" @click="getStrategy(batchTraffic)"></i> 
                  </el-popover>
                  </el-form-item>
                </div>
              </div>
              <div class="form-block" style="padding-bottom:20px">
                <span>Batch Delete Strategy</span>
                <el-alert
                  style="margin-top:8px"
                  title="Warning: batch setting will delete the api's specific setting, please operate with caution!"
                  :closable=false
                  type="warning">
                </el-alert>
                <div style="margin-top:10px">
                  <el-button @click="batchDelStrategy(1)">Del White/Black</el-button>
                  <el-button style="margin-left:10px;" @click="batchDelStrategy(2)">Del Timeout/Retry</el-button>
                  <el-button style="margin-left:10px;" @click="batchDelStrategy(2)">Del Traffic Control</el-button>
                </div>
             
              </div>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="batchStrategyVisible=false">Cancel</el-button>
          <el-button  type="primary" @click="submitBatchStrategy">
            <span>Submit</span>
          </el-button>
        </div>
      </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import waves from '@/directive/waves' // 水波纹指令
import {proxy} from '@/api/apiProxy'
import request from '@/utils/request' 
import apiDefine from './components/apiDefine'
import apiTest from './components/apiTest'
import bwlist from './components/bwlist'
import retry from './components/retry'
import traffic from './components/traffic'
import auditLog from './components/auditLog'
import clipboard from '@/directive/clipboard/index.js'
import clip from '@/utils/clipboard'

export default {
  name: 'api',
  directives: {
    waves
  },
  components: { apiDefine,apiTest,auditLog,bwlist,retry,traffic,clipboard,clip},
  data() {
      return {
        services:[],
        selectedService: '',
        
        apis: [],
        totalApis: 0,
        currentPage: 1,

        tempApi: {
          route_type: '1'
        },

        apiDefineVisible: false,
        defineStatus: 'create',


   
        testApi: {
        },
        
        selbwType: 1,

        switchActive: 1,
        switchInactive: 0,
        
        veryfyVisible: false,
        ruleContent: '',

        tempVerifyRule: {},

        // api定义
        apiDefineStep: 0,

        // api管理
        apiManageVisible: false,
        selectManageIndex: 1,

        apiTitle: '',

        strategies: [],
        selBW: '',
        selRetry: '',
        selTraffic: '',

        strategy: {
            content:{}
        },

        // 审计日志 
        totalLogs: 0,
        logs: [],

        q: '',

        batchStrategyVisible:false,
        batchBW: '',
        batchRetry: '',
        batchTraffic: '',
        selectApis: [],

        apps: [],

        tempService: '',
        tempApps : []
      }
  },
  methods: {
    handleTempService(s) {
      this.tempService = s
      if (s == '') {
        this.tempApps = []
        return 
      }
      // load apps
      request({
        url: '/ops/app/query',
        method: 'GET', 
        params: {
            service: s
        }
      }).then(res => {
          this.tempApps = res.data.data
      })
    },
    apiGrafanaAddr(api_id) {
      var addr = process.env.API_GRAFANA + api_id
       return  addr
    },
    genAPIIDByType() {
       if (this.tempApi.path_type==0) {
         return this.tempApi.service + '.' + this.tempApi.api_id + '.v' + this.tempApi.version
      } else {
        return this.tempApi.api_id + '.v' +this.tempApi.version
      }
    },
    apiNameAlert() {
      if (this.tempApi.path_type==0) {
         return this.$t('juz.apiIDTips1')
      } else {
        return this.$t('juz.apiIDTips2')
      }
    },
    apiNameTooltip() {
      if (this.tempApi.path_type==0) {
         return this.$t('juz.apiNameTips1')
      } else {
        return  this.$t('juz.apiNameTips2')
      }
    },
    showName(apiID) {
        return apiID.substring(0,apiID.length-3)
    },
    showVersion(apiID) {
        return apiID[apiID.length-1]
    },
    batchRelease() {
      if (this.selectApis.length == 0) {
        this.$message({
          message: 'No apis being selected',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
      this.$confirm('Are you sure to batch release the selected apis?', 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/api/batchRelease',
              api_ids: s
          }
          proxy('POST',params).then(res => {
            this.loadApis(this.currentPage)
            this.$message({
              message: 'Batch release successfully',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        });
    },
    batchDelStrategy(type) {
      this.$confirm('Are you sure to batch delete the strategy?', 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/api/batchDelStrategy',
              api_ids: s,
              type: type,
              service: this.selectedService
          }
          proxy('POST',params).then(res => {
            this.loadApis(this.currentPage)
            this.$message({
              message: 'Batch delete successfully',
              type: 'success',
              duration: 3 * 1000,
              center: true
            })
            this.batchStrategyVisible = false
          })
        });
    },
    submitBatchStrategy() {
      if (this.batchBW =='' && this.batchRetry == '' && this.batchTraffic == '') {
        this.$message({
          message: 'No strategy being selected',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
      
      this.$confirm("Batch setting will overide the api's specific setting,want to continue?", 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/api/batchStrategy',
              api_ids: s,
              batch_bw: this.batchBW,
              batch_retry: this.batchRetry,
              batch_traffic: this.batchTraffic,
              service: this.selectedService
          }
          proxy('POST',params).then(res => {
            this.loadApis(this.currentPage)
            this.$message({
              message: 'Batch strategy successfully',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
              this.batchStrategyVisible = false
          })
        });
    },
    handleBatchStrategy() {
      if (this.selectApis.length == 0) {
        this.$message({
          message: 'No apis being selected',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
      this.batchStrategyVisible = true
    },
     selTableRows(val) {
       this.selectApis = val
     },
     showName(apiID) {
       return apiID.substring(0,apiID.length-3)
     },
     showVersion(apiID) {
       return apiID[apiID.length-1]
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
    getStrategy(id) {
      if (id=='') {
        this.$message({
          message: 'No strategy being selected',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
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
    selectBW(val) {
      if (val=='') {
        this.tempApi.bw_strategy = 0
      } else {
        this.tempApi.bw_strategy = val 
      }
      
      this.selBW = val
    },
    selectRetry(val) {
       if (val=='') {
         this.tempApi.retry_strategy = 0
       } else {
         this.tempApi.retry_strategy = val 
       }
       
       this.selRetry = val
    },
    selectTraffic(val) {
       if (val=='') {
         this.tempApi.traffic_strategy = 0
       } else {
         this.tempApi.traffic_strategy = val 
       }
       
       this.selTraffic = val
    },
    filterStrategies(tp) {
      var temp = []
      for (var i=0;i<this.strategies.length;i++) {
        if (this.strategies[i].type == tp) {
          temp.push(this.strategies[i])
        }
      }
      return temp
    },
    strategyName(id) {
      if (id == 0) {
        return 
      }
      for (var i=0;i<this.strategies.length;i++) {
        if (this.strategies[i].id == id) {
          return this.strategies[i].name
        }
      }
    },
    offlineAPI(api) {
      this.$confirm('You will offline this api：' + api.api_id + '\n,when offlined, clients will not access this api anymore', 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/api/offline',
              api_id: api.api_id,
          }
          proxy('POST',params).then(res => {
            api.release_version=''
            this.$message({
              message: 'Offline successfully',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        })
    },
    releaseAPI(api) {
       this.$confirm('You will release this api：' + api.api_id, 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'Cancel',
          confirmButtonText: 'Submit',
          type: 'info'
        }).then(() => {
          var params = {
              target_app: 'juzManage',
              target_service: api.service,
              target_path: '/manage/api/release',
              api_id: api.api_id,
          }
          proxy('POST',params).then(res => {
            api.release_version=api.revise_version
            this.loadApis(this.currentPage)
            this.$message({
              message: 'Api released successfully',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        })
    },
    delParamRule(param) {
      for (var i=0;i<this.tempApi.param_rules.length;i++) {
        if (this.tempApi.param_rules[i].param == param) {
          this.tempApi.param_rules.splice(i,1)
        }
      }
    },
    addVerifyRule() {
      if (this.tempApi.param_rules == null) {
        this.tempApi.param_rules = []
      } else {
        // 检查是否存在重复的规则
        for (var i=0;i<this.tempApi.param_rules.length;i++) {
          if (this.tempVerifyRule.param == this.tempApi.param_rules[i].param) {
            this.$message({
              message: 'Param exists',
              type: 'warning',
              duration: 3 * 1000,
              center: true
            })
            return
          }
        }
      }
      // 验证正则是否合法
       var params = {
            target_app: 'juzManage',
            target_path: '/manage/api/verifyParamRule',
            param: this.tempVerifyRule.param,
            rule: this.tempVerifyRule.rule,
            test_data: this.tempVerifyRule.test_data
        }
        proxy('POST',params).then(res => {
          this.tempApi.param_rules.unshift({
            param: this.tempVerifyRule.param,
            rule: this.tempVerifyRule.rule,
            test_data: this.tempVerifyRule.test_data
          }) 
          this.tempVerifyRule = {}
          this.$message({
              message: 'Test successfully',
              type: 'success',
              duration: 3 * 1000,
              center: true
            })
        })
    },
    cancelAPIDefine() {
      this.apiDefineVisible = false
    },
    cancelVerify() {
      this.veryfyVisible = false
    },
    showVerify() { 
      this.tempVerifyRule.param = ''
      this.tempVerifyRule.rule = ''
      this.tempVerifyRule.test_data = ''

      this.veryfyVisible = true
      this.ruleContent = `
         <p><strong>Demo</strong></p>
        <p>mobile : ^1[0-9]{10}$</p>
        <p>email : ^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$</p>

        <p><strong>Demo analysis</strong></p>
        <p>1. Let's assume our http params is like this: k1=v1&k2=v2&k3=v3</p>
        <p>2. Now we want to verify the mobile param,we use the regexp rule: ^1[0-9]{10}$</p>
        <p>so，a client can access with http://xxx.com/index?mobile=15880271182&email=aaa@bbb.com}</p>
      `
      
    },

    onManageSelect(index) {
        this.selectManageIndex = index
    },
    handleManage(api) {
      this.apiDefineStep = 0
      this.defineStatus = 'edit'
      this.tempApi = api
      // 打开API面板时，因为数据都是一次性加载，因此我们需要重新加载数据
      //加载审计日志
      var params = {
          target_app: 'juzManage',
          target_path: '/manage/auditLog/count',
          target_type: 2,
          target_id: this.tempApi.api_id
      }
      proxy('POST',params).then(res => {
          this.totalLogs = res.data.data
      })

      var params = {
          target_app: 'juzManage',
          target_path: '/manage/auditLog/load',
          target_type: 2,
          target_id: this.tempApi.api_id,
          page: 1
      }
      proxy('POST',params).then(res => {
          this.logs = res.data.data
      })
      this.apiManageVisible = true
    },
    cancelManage() {
      this.apiManageVisible = false
    },
    handleEdit() {
      this.apiTitle = 'Edit API'
      this.apiDefineVisible = true
      if (this.tempApi.bw_strategy != 0) {
        this.selBW = this.tempApi.bw_strategy
      }

      if (this.tempApi.retry_strategy != 0) {
        this.selRetry = this.tempApi.retry_strategy
      }

       if (this.tempApi.traffic_strategy != 0) {
        this.selTraffic = this.tempApi.traffic_strategy
      }
    },
    delApi(api) {
      if (api.release_version != '') {
        this.$message({
          message: 'Please offline the api first',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        })
        return 
      }
      this.$confirm('Are you sure to delete this api?', 'Warning', {
        dangerouslyUseHTMLString: true,
        cancelButtonText: 'Cancel',
        confirmButtonText: 'Submit',
        type: 'error'
      }).then(() => {
          var params = {
              target_app: 'juzManage',
              target_path: '/manage/api/delete',
              service: api.service,
              api_id: api.api_id
          }
          proxy('POST',params).then(res => {
            if (res.data.status != 200) {
                this.$message({
                  message: res.data.message,
                  type: 'warning',
                  duration: 3 * 1000,
                  center: true
                })
                return 
            }
            this.apiManageVisible = false
            this.loadApis(this.currentPage)
            this.$message({
                message: 'Delete api successfully',
                type: 'success',
                duration: 3 * 1000,
                center: true
            })
          })
        })
    },
    saveApi() {
      if (this.tempApi.api_id == '') {
        this.$message({
              message: "Api name cant not be empty",
              type: 'error',
              duration: 3 * 1000,
              center: true
        })
        return 
      }
      

      if (this.tempApi.app == '') {
        this.$message({
              message: "You must choose a application the api belongs",
              type: 'error',
              duration: 3 * 1000,
              center: true
        })
        return 
      }

      var s = Object.assign({},   this.tempApi) 
      s.param_rules = JSON.stringify(s.param_rules)
      if (this.defineStatus == 'create') {
        s.api_id = this.genAPIIDByType()
      } 
      var params = {
          target_app: 'juzManage',
          target_path: '/manage/api/define',
          target_service: this.tempApi.service,
          api : s,
          action: this.defineStatus
      }
      proxy('POST',params).then(res => {
        res.data.data.param_rules = JSON.parse(res.data.data.param_rules)
        this.apiDefineVisible = false 
        if (this.defineStatus =='create') {
          this.apis.unshift(res.data.data)
        } else {
          this.tempApi = res.data.data
          this.loadApis(this.currentPage)
        }
        if (res.data.message != '') {
          this.$message({
            message:res.data.message,
            type: 'info',
            duration: 3 * 1000,
            center: true
          })
        } else {
          this.$message({
            message: 'Edit api ok',
            type: 'success',
            duration: 3 * 1000,
            center: true
        })
        }
        
      })
    },
      // 定义API
    nextDefineStep(forward) {
      if (forward) {
         this.apiDefineStep++
      } else {
         this.apiDefineStep--
      }
    },
    handleCreate() {
      if (this.selectedService=='') {
          this.$message({
              message: 'Select a service first',
              type: 'warning',
              duration: 3 * 1000,
              center: true
          })
          return 
      }
      this.apiTitle = 'Create API'
      this.apiDefineStep = 0
      this.tempApi = {
         service: this.selectedService,
         app: '',
         api_id: '',
         version: 1,
         path_type:0,
         desc:'',        
         route_type: 1,
         addr_type: 1,
         backend_addr: '',
         backend_uri: '',
         backend_type: 1,
         mock_data: '',
         method: 'POST',
         bw_strategy:0,
         retry_strategy:0,
         traffic_api: '',
         traffic_ratio: 0,
         traffic_on: 0,
         traffic_ips: '',
         verify_on: 0,
         cached_time: 0,
         param_rules: []
      }

      this.apiDefineVisible = true
      this.defineStatus = 'create'
    },
    calcService() {
        return this.selectedService || this.$store.getters.service
    },
    handleSelService(service) {
       this.$store.dispatch('setService', service)
       this.selectedService = service
       this.loadApis(1)
       this.loadStrategy(this.selectedService)
       this.loadApps()
    },
    loadStrategy(service) {
       var params = {
          target_app: 'juzManage',
          target_path: '/manage/strategy/load',
          service: this.selectedService,
          type: 0, //加载所有的策略
        }
        proxy('POST',params).then(res => {
          this.strategies = res.data.data
        })
    },
    loadApis(page) {
       this.currentPage = page
        //加载审计日志
        var params = {
            target_service: this.selectedService,
            target_app: 'juzManage',
            target_path: '/manage/api/count',
            q: this.q,
            service: this.selectedService
        }
        proxy('POST',params).then(res => {
            this.totalApis = res.data.data
        })

       // 加载Service底下的所有API
        var params = {
          target_service: this.selectedService,
          target_app: 'juzManage',
          target_path: '/manage/api/query',
          service: this.selectedService,
          q: this.q,
          page: page
        }
        proxy('POST',params).then(res => {
           this.apis = res.data.data
           console.log(this.apis)
           for (var i=0;i<this.apis.length;i++) {
             this.apis[i].param_rules = JSON.parse(this.apis[i].param_rules)
           }
        })
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
        this.loadApis(1)
        this.loadStrategy(this.selectedService)
    }
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>

<style lang="scss"  scoped>

</style>

<style lang="scss">
.no-border-input {
  .el-input__inner {
    border-top:none;
    border-left:none;
    border-right:none;
    border-radius:0;
    height: 26px
  }

  .el-input--mini {
    input {
      border-bottom:none;
    }
  }
    span {
      border: none !important
    }
}
.el-input.is-disabled {
  .el-input__inner {
    color:black
  }
}

.el-textarea.is-disabled {
  .el-textarea__inner {
    color:black
  }
}
.el-input.el-input-group--prepend {
  .el-input-group__prepend {
    color:black
  }
}

.note {
  font-size: 13px;
  margin-top:0px;
}


</style>