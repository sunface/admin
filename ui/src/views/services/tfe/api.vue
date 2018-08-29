<template>
  <div class="app-container">
      <div class="filter-container">
        
        <el-select clearable class="filter-item" v-model="selectedService" @change='handleSelService' style="width: 200px"  placeholder="请选择Service">
          <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name">
          </el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"  v-waves type="primary" icon="el-icon-edit">创建API</el-button>

       

        <span style="float:right;margin-right:30px">
          <el-button class="filter-item"  @click="batchRelease"  v-waves icon="el-icon-upload2">批量发布</el-button>
          <el-button class="filter-item" style="margin-right: 20px;" @click="handleBatchStrategy"  v-waves icon="el-icon-bell">批量设置策略</el-button>
          <el-input  style="width: 250px;" class="filter-item" placeholder="模糊查询" v-model="q" clearable @clear="loadApis(1)" @keyup.enter.native="loadApis(1)" >
          </el-input>
          <el-button class="filter-item" type="success" v-waves icon="el-icon-search"   @click="loadApis(1)"></el-button>
        </span>

      </div>

      <div class="table">
        <el-table  :data="apis"  fit highlight-current-row style="width: 100%;"  :default-sort = "{prop: 'gmt_modified', order: 'descending'}" @selection-change="selTableRows">
           <el-table-column
            type="selection"
            width="55"> 
          </el-table-column>
          <el-table-column align="left" label="API ID" width="210" prop="api_id">
            <template slot-scope="scope">
              <span>{{scope.row.api_id}}</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="路径映射" width="100" prop="api_id">
            <template slot-scope="scope">
              <span v-if="scope.row.path_type==0">否</span>
              <span v-else>是</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="API名" width="160">
            <template slot-scope="scope">
              <span>{{showName(scope.row.api_id)}}</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="版本号" width="80">
            <template slot-scope="scope">
              <span>{{showVersion(scope.row.api_id)}}</span>
            </template>
          </el-table-column>
          <el-table-column width="300" align="left" label="目标地址" prop="route_addr">
            <template slot-scope="scope">
              <span>{{scope.row.route_addr}}</span>
            </template>
          </el-table-column>
          <el-table-column width="250" align="left" label="发布状态" >
             <template slot-scope="scope">
                <div><el-tag type="success" size="mini">编辑版本</el-tag> <span style="color:#337ab7;font-size:13px;">{{scope.row.revise_version}}</span></div>
                <div v-if="scope.row.release_version!=''" style="margin-top:3px"><el-tag  size="mini" :type="calcRelease(scope.row)">发布版本</el-tag>  <span style="color:#337ab7;font-size:13px;">{{scope.row.release_version}}</span></div>
                <div v-else style="margin-top:3px"><el-tag type="warning" size="mini">尚未发布</el-tag></div>
            </template>
          </el-table-column>
          <el-table-column align="center" label="操作" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button  v-if="scope.row.release_version!=scope.row.revise_version && scope.row.release_version!=''" type="text" @click="releaseAPI(scope.row)">重新发布</el-button>
              <el-button  v-if="scope.row.release_version=='' " type="text" @click="releaseAPI(scope.row)">发布上线</el-button>
              <el-button  type="text"  @click="handleManage(scope.row)">管理</el-button>
              <el-button  class="green-button"  type="text" @click="handleCopy(scope.row,$event)">复制配置</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
            layout="total,prev, pager, next"
            @current-change="loadApis"
            :total="totalApis">
        </el-pagination>
      </div>
      <el-dialog class="black-dialog" :title="'管理API ：'+tempApi.api_id" :visible.sync="apiManageVisible" top="40px"  :fullscreen=true style="backgroud:#545c64 !important"> 
          <el-row :gutter="20" style="margin-top:40px;">
            <el-col :span=3 :offset=1>
              <el-menu
                default-active="1"
                class="el-menu-vertical-demo"
                @select="onManageSelect"
                text-color="#fff"
                background-color="#545c64"
                active-text-color="#ffd04b"
                ref="apiMenu"
                >
                <el-menu-item index="1">
                   <i class="el-icon-setting"></i> 
                  <span slot="title">API定义</span>
                </el-menu-item>
                <el-menu-item index="2">
                  <i class="el-icon-caret-right"></i>
                  <span slot="title">API调试</span>
                </el-menu-item>
                <el-menu-item index="3">
                   <i class="el-icon-menu"></i>
                  <span slot="title">审计日志</span>
                </el-menu-item>
                <el-menu-item index="4">
                   <i class="el-icon-bell"></i>
                  <span slot="title">指标监控</span>
                </el-menu-item>
                <el-menu-item index="5">
                   <i class="el-icon-document"></i>
                  <span slot="title">API文档</span>
                </el-menu-item>
              </el-menu>
            </el-col>
            <el-col :span=20 :offset=0 style="color:#fff;padding-right:0;">
              <!-- api定义 -->
                <div v-show="selectManageIndex==1">
                     <div style="width:85%">
                        <span style="float:right;margin-top:-36px">
                          <el-button type="success" size="small" @click="handleEdit()">编辑</el-button>
                          <el-button v-if="tempApi.release_version!=tempApi.revise_version && tempApi.release_version!=''" type="warning" size="small" @click="releaseAPI(tempApi)">重新发布</el-button>
                          <el-button v-if="tempApi.release_version=='' " type="warning" size="small" @click="releaseAPI(tempApi)">发布上线</el-button>
                         
                          <el-button v-show="tempApi.release_version!=''" type="info" size="small" @click="offlineAPI(tempApi)">下线</el-button>
                           <el-button  type="" @click="delApi(tempApi)" size="small">删除</el-button>
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

      <el-dialog :title="apiTitle" :visible.sync="apiDefineVisible" top="40px" :before-close="cancelAPIDefine" :fullscreen=false> 
        <el-steps :active="apiDefineStep" finish-status="success">
          <el-step title="API信息"></el-step>
          <el-step title="请求设置"></el-step>
          <el-step title="高级设置"></el-step>
        </el-steps>

        <el-row :gutter="20" style="margin-top:40px;">
            <el-col :span=20 :offset=2>
              <div v-show="apiDefineStep==0">
                  <el-form label-position="left" label-width="70px" size="small">
                    <el-form-item label="Service">
                      {{tempApi.service}}
                    </el-form-item>
                    <el-form-item label="路径映射">
                       <el-switch
                          v-model="tempApi.path_type"
                          active-color="#13ce66"
                          :active-value='switchActive'
                          :inactive-value='switchInactive'
                          :disabled="defineStatus=='edit'">
                        </el-switch>
                        <el-alert
                        :closable=false
                        title="路径映射支持客户端直接通过http://yourip/a/b/c的形式进行访问,非路径映射是用http://yourip/service/api?api_name=xxx进行访问"
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
                    <el-form-item v-if="defineStatus=='create'" label="API名" prop="name" class="first-item">
                      <span v-if="tempApi.path_type==0">{{tempApi.service}}.</span>
                      <el-tooltip class="no-border-input" :content="apiNameTooltip()" placement="top">
                          <el-input  v-model="tempApi.api_id" placeholder=""  style="width:250px">
                          </el-input>
                       </el-tooltip>
                    </el-form-item>
                    <el-form-item v-else label="API名称" prop="name" class="first-item">
                      {{showName(tempApi.api_id)}}
                    </el-form-item>
                    <el-form-item label="版本号">
                      <el-input-number v-if="defineStatus=='create'" class="no-border-input" style="width:90px;border:none" v-model="tempApi.version"  :min="1" :max="9" size="mini"></el-input-number>
                      <span v-else>{{showVersion(tempApi.api_id)}}</span>
                    </el-form-item>
                    <el-form-item label="描述">
                      <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder="简单描述该API" v-model="tempApi.desc">
                      </el-input>
                    </el-form-item>
                  </el-form>
              </div>
              <div v-show="apiDefineStep==1">
                  <el-form label-position="left" label-width="150px" size="small">
                    <el-form-item label="代理类型" prop="route_type">
                      <el-select  class="filter-item" v-model="tempApi.route_type" style="width: 150px"  placeholder="请选择代理类型">
                        <el-option label="同步代理" :value=1></el-option>
                        <el-option  label="重定向" :value=2></el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item label="后端服务类型" prop="name" class="first-item">
                      <el-radio-group v-model="tempApi.route_proto">
                        <el-radio :label="1">HTTP(s)服务</el-radio>
                        <el-radio :label="2">Mock</el-radio>
                      </el-radio-group> 
                    </el-form-item>
                    <div v-if="tempApi.route_proto==1">
                      <el-form-item label="后端服务地址" required>
                        <el-tooltip content="待访问的后段服务地址" placement="top">
                            <el-input v-model="tempApi.route_addr"></el-input>
                        </el-tooltip>
                      </el-form-item>
                      <el-form-item label="Header和Cookie">
                          <el-tag type="success" size="large" style="border:none;">默认透传</el-tag>
                      </el-form-item>
                       <el-form-item label="请求参数">
                          <el-tag type="success" size="large" style="border:none;">默认透传</el-tag>
                      </el-form-item>
                    </div>
                    <div v-else>
                      <el-form-item label="Mock返回">
                        <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder="这里请填写数据的元信息，例如返回json就填写json格式,返回字符串就填写字符串，网关不做任何转换，直接返回给请求的客户端" v-model="tempApi.mock_data"></el-input>
                          <el-alert
                          :closable=false
                            title="注意！设置Mock后，网关会自动模拟数据，实际请求不会到达后端服务"
                            type="warning">
                          </el-alert>
                      </el-form-item>    
                    </div>
                  
                
                  </el-form>
              </div>
              <div v-show="apiDefineStep==2">
                  <el-form label-position="left" label-width="150px" size="small">
                      <div class="form-block">
                          <span>标签分组<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                          <el-form-item label="标签" prop="route_type" style="margin-top:10px"> 
                            <el-select  clearable class="filter-item" v-model="tempApi.label" style="width: 150px"  placeholder="选择标签">
                              <el-option v-for="l in labels"  :key="l" :label="l" :value="l"></el-option>
                            </el-select>
                          </el-form-item>
                          <el-alert
                            :closable=false
                            title="您可以为同一个系统或者功能的API设置一样的标签组，然后就可以对该系统进行监控和数据统计了"
                            type="success">
                          </el-alert>
                      </div>
                      <div class="form-block">
                        <span>参数验证<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                        <el-form-item label="开启验证" style="width:300px" class="first-item">
                            <el-switch
                              v-model="tempApi.verify_on"
                              active-color="#13ce66"
                              :active-value='switchActive'
                              :inactive-value='switchInactive'>
                            </el-switch>
                        </el-form-item>
                        <el-form-item label="验证表" style="width:300px" class="first-item">
                          <el-button type="success" icon="el-icon-edit"  @click="showVerify">编辑内容</el-button>
                        </el-form-item>
                      </div>

                      <div class="form-block">
                        <span>流量路由(金丝雀测试)<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                        <el-form-item label="开启路由" style="width:300px" class="first-item">
                            <el-switch
                              v-model="tempApi.traffic_on"
                              active-color="#13ce66" 
                              :active-value='switchActive'
                              :inactive-value='switchInactive'>
                            </el-switch>
                        </el-form-item>
                        <el-form-item label="目标api" prop="req_timeout" style="width:500px" class="first-item">
                          <el-tooltip content="指定路由到的api名" placement="top">
                            <el-input v-model="tempApi.traffic_api" placeholder="">
                            </el-input>
                          </el-tooltip>
                        </el-form-item>
                        <el-form-item label="路由比例" prop="req_timeout">
                          <el-tooltip content="将指定的比例路由到上面的目标API,在0和100之间取值" placement="top">
                              <el-input-number v-model="tempApi.traffic_ratio" :min="0" :max="100"></el-input-number>
                          </el-tooltip>
                        </el-form-item>
                        <el-form-item label="IP列表" style="width:300px" class="first-item">
                          <el-tooltip content="精确指定哪些IP的客户端被路由到上面的目标API" placement="top">
                            <el-input v-model="tempApi.traffic_ips" placeholder="ip之间空格分割">
                            </el-input>
                          </el-tooltip>
                        </el-form-item>
                      </div>

                       <el-alert
                        style="margin-top:8px"
                            title="小提示: 您可以去策略页面新建策略，再进行下面的设置"
                            :closable=false
                            type="success">
                          </el-alert>
                      <div class="form-block">
                        <span>黑白名单<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                        <el-form-item label="选择策略" prop="req_timeout">
                          <el-select clearable class="filter-item" v-model="selBW" @change="selectBW" style="width: 200px"  placeholder="请选择">
                            <el-option v-for="s in  filterStrategies(1)" :key="s.name" :label="s.name" :value="s.id">
                            </el-option>
                          </el-select>
                          <el-popover
                            v-show="this.selBW != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <bwlist :strategy="strategy"></bwlist>
                            <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(selBW)"></i> 
                          </el-popover>
                        </el-form-item>
                      </div>
                      <div class="form-block">
                        <span>超时重试<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                        <el-form-item label="选择策略" prop="req_timeout">
                          <el-select clearable class="filter-item" v-model="selRetry" @change="selectRetry" style="width: 200px"  placeholder="请选择">
                            <el-option v-for="s in  filterStrategies(2)" :key="s.name" :label="s.name" :value="s.id"> 
                            </el-option>
                          </el-select>
                          <el-popover
                            v-show="this.selRetry != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <retry :strategy="strategy"></retry>
                            <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(selRetry)"></i> 
                        </el-popover>
                        </el-form-item>
                      </div>
                      <div class="form-block">
                        <span>流量控制<el-tag type="success" size="mini" style="border:none;">可选</el-tag></span>
                        <el-form-item label="选择策略" prop="req_timeout">
                          <el-select clearable class="filter-item" v-model="selTraffic" @change="selectTraffic" style="width: 200px"  placeholder="请选择">
                            <el-option v-for="s in  filterStrategies(3)" :key="s.name" :label="s.name" :value="s.id"> 
                            </el-option>
                          </el-select>
                          <el-popover
                            v-show="this.selTraffic != ''"
                            placement="top"
                            width="600"
                            trigger="click">
                            <traffic :strategy="strategy"></traffic>
                            <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(selTraffic)"></i> 
                        </el-popover>
                        </el-form-item>
                      </div>
                  </el-form>
              </div>
            </el-col>
        </el-row>

         <div slot="footer" class="dialog-footer">
            <el-button @click="nextDefineStep(false)" v-show="apiDefineStep!=0">上一步</el-button>
            <el-button @click="nextDefineStep(true)" v-show="apiDefineStep!=2" type="success">下一步</el-button>
            <el-button  type="primary" @click="saveApi" v-show="apiDefineStep!=0">保存</el-button>
          </div>
      </el-dialog>


      <el-dialog
        title="参数验证"
        :fullscreen=true
        :visible.sync="veryfyVisible"
      >
        <el-row :gutter="20" style="margin-top:0px;">
            <el-col :xs="{span:24}" :sm="{span:18}" :md="{span: 8}" :lg="{ span: 6}" style="border-right:.5px solid #999">
                <div v-html="ruleContent" ></div>
            </el-col>
            <el-col :span=16 :offset=1>
                <el-form label-position="left" label-width="70px">
                  <div>
                    <span>为API: {{tempApi.api_id}} 新建规则</span>
                    <el-form-item label="参数名" style="margin-top:15px">
                      <el-input v-model="tempVerifyRule.param" placeholder="参数名" style="width:150px" >
                      </el-input>
                    </el-form-item>
                    <el-form-item label="规则">
                      <el-input v-model="tempVerifyRule.rule" placeholder="规则，正则表达式" style="width:350px">
                      </el-input>
                    </el-form-item>
                    <el-form-item label="测试数据">
                      <el-input v-model="tempVerifyRule.test_data" placeholder="给出测试数据" style="width:150px">
                      </el-input>
                    </el-form-item>
                     <el-button @click="addVerifyRule" type="success" >创建并测试</el-button>
                  </div>
              
                </el-form>

                <div class="table" style="margin-top:15px;overflow-y:auto">
                  <el-table  :data="tempApi.param_rules" border fit highlight-current-row style="width: 100%;">
                    <el-table-column align="center" label="参数名" width="150">
                      <template slot-scope="scope">
                        <span>{{scope.row.param}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column width="400" align="center" label="规则">
                      <template slot-scope="scope">
                        <span>{{scope.row.rule}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column width="250" align="center" label="测试数据" >
                      <template slot-scope="scope">
                          <span>{{scope.row.test_data}}</span>
                      </template>
                    </el-table-column>
                    <el-table-column align="center" label="操作" class-name="small-padding fixed-width">
                      <template slot-scope="scope"> 
                        <el-button  type="text" size="mini" icon="el-icon-minus" circle @click="delParamRule(scope.row.param)"></el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
            </el-col>
        </el-row>

        <span slot="footer" class="dialog-footer">
          <el-button @click="cancelVerify">关闭</el-button>
        </span>
      </el-dialog>

      <!-- 批量设置策略 -->
      <el-dialog title="批量策略管理" :visible.sync="batchStrategyVisible" close-on-press-escape=true>
        <el-form  label-position="left" label-width="120px" style='width: 650px; margin-left:50px;' size="mini">
            <div class="form-block">
              <span>批量设置策略</span>
              <el-alert
                  style="margin-top:8px"
                  title="警告 ：批量设置会强制覆盖之前的策略，请谨慎操作!"
                  :closable=false
                  type="warning">
              </el-alert>
              <el-alert
                  style="margin-top:8px"
                  title="提示 ：未选择的策略将被忽略，不会将API的对应策略设置为空"
                  :closable=false
                  type="success">
              </el-alert>
                <div class="form-block">
                  <span>黑白名单</span>
                  <el-form-item label="选择策略">
                    <el-select clearable class="filter-item" v-model="batchBW" style="width: 200px"  placeholder="请选择">
                      <el-option v-for="s in  filterStrategies(1)" :key="s.name" :label="s.name" :value="s.id">
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchBW != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <bwlist :strategy="strategy"></bwlist>
                      <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(batchBW)"></i> 
                    </el-popover>
                  </el-form-item>
                </div>
                <div class="form-block">
                  <span>超时重试</span>
                  <el-form-item label="选择策略" prop="req_timeout">
                    <el-select clearable class="filter-item" v-model="batchRetry"  style="width: 200px"  placeholder="请选择">
                      <el-option v-for="s in  filterStrategies(2)" :key="s.name" :label="s.name" :value="s.id"> 
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchRetry != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <retry :strategy="strategy"></retry>
                      <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(batchRetry)"></i> 
                  </el-popover>
                  </el-form-item>
                </div>
                <div class="form-block">
                  <span>流量控制</span>
                  <el-form-item label="选择策略" prop="req_timeout">
                    <el-select clearable class="filter-item" v-model="batchTraffic"  style="width: 200px"  placeholder="请选择">
                      <el-option v-for="s in  filterStrategies(3)" :key="s.name" :label="s.name" :value="s.id"> 
                      </el-option>
                    </el-select>
                    <el-popover
                      v-show="this.batchTraffic != ''"
                      placement="top"
                      width="600"
                      trigger="click">
                      <traffic :strategy="strategy"></traffic>
                      <i slot="reference" class="el-icon-view hover-cursor" style="margin-left: 8px;color: #67c23a" @click="getStrategy(batchTraffic)"></i> 
                  </el-popover>
                  </el-form-item>
                </div>
              </div>
              <div class="form-block" style="padding-bottom:20px">
                <span>批量删除策略</span>
                <el-alert
                  style="margin-top:8px"
                  title="警告 ：批量删除会删除所选API的对应策略，请谨慎操作!"
                  :closable=false
                  type="warning">
                </el-alert>
                <div style="margin-top:10px">
                  <el-button @click="batchDelStrategy(1)">删除黑白名单</el-button>
                  <el-button style="margin-left:10px;" @click="batchDelStrategy(2)">删除超时重试</el-button>
                  <el-button style="margin-left:10px;" @click="batchDelStrategy(2)">删除流量控制</el-button>
                </div>
             
              </div>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="batchStrategyVisible=false">取消</el-button>
          <el-button  type="primary" @click="submitBatchStrategy">
            <span>提交</span>
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
        createRules: {
          name: [{ required: false, message: '请输入API Name,例如party.account.create', trigger: 'blur' }],
          route_addr: [{ required: false, message: '请输入目标地址', trigger: 'blur' }]
        },
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

        labels: []
      }
  },
  methods: {
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
         return "API ID是由Service + API名 + 版本号自动生成的，无需手动设置"
      } else {
        return "API ID是由API名 + 版本号自动生成的，无需手动设置"
      }
    },
    apiNameTooltip() {
      if (this.tempApi.path_type==0) {
         return "必须是a.b这类形式，用.分割，只支持英文字母和数字"
      } else {
        return "必须是url路径形式，例如/path1/path2"
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
          message: '未选择任何API',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
      this.$confirm('您确定要批量发布选中的API吗？', '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_service: 'tfe.manage',
              target_path: '/manage/api/batchRelease',
              api_ids: s
          }
          proxy('POST',params).then(res => {
            this.loadApis(this.currentPage)
            this.$message({
              message: '批量发布成功',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        });
    },
    batchDelStrategy(type) {
      this.$confirm('警告 ：批量删除会删除所选API的对应策略，您确定要继续吗？', '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_service: 'tfe.manage',
              target_path: '/manage/api/batchDelStrategy',
              api_ids: s,
              type: type,
              service: this.selectedService
          }
          proxy('POST',params).then(res => {
            this.loadApis(this.currentPage)
            this.$message({
              message: '批量删除成功',
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
          message: '没有选择任何策略',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
      this.$confirm('警告 ：批量设置会强制覆盖已选择API的策略设置，您确定要继续吗？', '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var apiIDs = []
          for (var i=0;i<this.selectApis.length;i++) {
            apiIDs.push(this.selectApis[i].api_id)
          }
          var s = JSON.stringify(apiIDs)
          var params = {
              target_service: 'tfe.manage',
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
              message: '批量设置成功',
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
          message: '未选择任何API',
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
        message: '复制成功',
        type: 'success',
        duration: 2000
      })
    },
    getStrategy(id) {
      if (id=='') {
        this.$message({
          message: '没有选择策略',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }
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
      this.$confirm('您将要下线以下API：' + api.api_id + '\n下线后，所有的外部用户将无法访问该API对应的服务', '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var params = {
              target_service: 'tfe.manage',
              target_path: '/manage/api/offline',
              api_id: api.api_id,
          }
          proxy('POST',params).then(res => {
            api.release_version=''
            this.$message({
              message: 'API下线成功',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消下线'
          });          
        });
    },
    releaseAPI(api) {
       this.$confirm('您将要发布以下API：' + api.api_id, '提示', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: '取消',
          confirmButtonText: '确定',
          type: 'info'
        }).then(() => {
          var params = {
              target_service: 'tfe.manage',
              target_path: '/manage/api/release',
              api_id: api.api_id,
          }
          proxy('POST',params).then(res => {
            api.release_version=api.revise_version
            this.loadApis(this.currentPage)
            this.$message({
              message: 'API发布成功',
              type: 'success',
              duration: 5 * 1000,
              center: true
            }) 
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消发布'
          });          
        });
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
              message: '参数名已存在',
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
            target_service: 'tfe.manage',
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
         <p><strong>示例演示</strong></p>
        <p>mobile : ^1[0-9]{10}$</p>
        <p>email : ^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$</p>

        <p><strong>示例解析</strong></p>
        <p>1. 我们假定http的请求参数为k1=v1&k2=v2&k3=v3的形式</p>
        <p>2. 首先我们对mobile字段进行了限制，要求必须符合"^1[0-9]{10}$"这种形式的正则表达式</p>
        <p>3. 暂时只支持最外层参数</p>
        <p>综上，符合上述请求格式的参数形式为http://xxx.com/index?mobile=15880271182&email=aaa@bbb.com}</p>
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
          target_service: 'tfe.manage',
          target_path: '/manage/auditLog/count',
          target_type: 2,
          target_id: this.tempApi.api_id
      }
      proxy('POST',params).then(res => {
          this.totalLogs = res.data.data
      })

      var params = {
          target_service: 'tfe.manage',
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
      this.apiTitle = '编辑API'
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
          message: '请先下线API，再删除',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        })
        return 
      }
      this.$confirm('此操作将永久删除<strong>该API</strong>, 是否继续删除?', '提示', {
        dangerouslyUseHTMLString: true,
        cancelButtonText: '不，赶紧取消删除',
        confirmButtonText: 'ok',
        type: 'error'
      }).then(() => {
          var params = {
              target_service: 'tfe.manage',
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
                message: '删除API成功',
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
    saveApi() {
      if (this.tempApi.api_id == '') {
        this.$message({
              message: "API名不能为空,且必须以'SERVICE.'为前缀",
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
          target_service: 'tfe.manage',
          target_path: '/manage/api/define',
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
            message: '编辑API成功',
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
              message: '请先选择Service',
              type: 'warning',
              duration: 3 * 1000,
              center: true
          })
          return 
      }
      this.apiTitle = '创建API'
      this.apiDefineStep = 0
      this.tempApi = {
         service: this.selectedService,
         api_id: '',
         version: 1,
         desc:''
,        route_type: 1,
         route_addr: 'http://',
         route_proto: 1,
         mock_data: '',
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
    handleSelService() {
       this.loadApis(1)
       this.loadStrategy(this.selectedService)
       this.loadLabels()
    },
    loadStrategy(service) {
       var params = {
          target_service: 'tfe.manage',
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
            target_service: 'tfe.manage',
            target_path: '/manage/api/count',
            q: this.q,
            service: this.selectedService
        }
        proxy('POST',params).then(res => {
            this.totalApis = res.data.data
        })

       // 加载Service底下的所有API
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/api/query',
                service: this.selectedService,
                q: this.q,
                page: page
            }
        proxy('POST',params).then(res => {
           this.apis = res.data.data
           for (var i=0;i<this.apis.length;i++) {
             this.apis[i].param_rules = JSON.parse(this.apis[i].param_rules)
           }
        })
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
               this.loadApis(1)
               this.loadStrategy(this.selectedService)
               this.loadLabels()
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