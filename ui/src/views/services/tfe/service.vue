<template>
  <div class="app-container">
      <div class="filter-container">
        <!-- <el-input  style="width: 200px;" class="filter-item" placeholder="Service">
        </el-input>
        <el-button class="filter-item" type="primary" style="margin-left: 10px;" v-waves icon="el-icon-search" @click="handleCreate">搜索</el-button> -->
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"  v-waves type="primary" icon="el-icon-plus">添加</el-button>
        <!-- <el-button class="filter-item" type="primary" :loading="downloadLoading" v-waves icon="el-icon-download" @click="handleDownload">导出</el-button> -->
      </div>

      <div class="table">
        <el-table  :data="list" border fit highlight-current-row style="width: 100%;min-height:1000px;"  :default-sort = "{prop: 'gmt_modified', order: 'descending'}">
          <el-table-column align="center" label="Name" width="200" prop="name" sortable>
            <template slot-scope="scope">
              <span>{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column width="100" align="center" label="创建者" prop="creator" sortable>
             <template slot-scope="scope">
              <span>{{scope.row.creator}}</span>
            </template>
          </el-table-column>
          <el-table-column width="250" align="center" label="管理员" prop="admins">
            <template slot-scope="scope">
              <span>{{scope.row.admins}}</span>
            </template>
          </el-table-column>
           <el-table-column width="200" align="center" label="你的权限" prop="priv" sortable>
            <template slot-scope="scope">
              <span v-show="scope.row.priv==0">创建者</span>
              <span v-show="scope.row.priv==1">管理员</span>
              <span v-show="scope.row.priv==2">普通用户</span>
              <span v-show="scope.row.priv==3">访客</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="操作" class-name="small-padding">
            <template slot-scope="scope">
              <el-button v-show="scope.row.priv==0 || scope.row.priv==1 || scope.row.priv==3"  size="mini" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button  size="mini" @click="handleMetrics(scope.row)">指标监控</el-button>
              <el-button v-show="scope.row.priv==0 || scope.row.priv== 3" size="mini" type="danger" @click="delService(scope.row.name)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>


      <el-dialog title="创建Service" :visible.sync="createServiceVisible">
        <el-form :rules="createRules" ref="dataForm" :model="tempCreate" label-position="left" label-width="70px" style='width: 400px; margin-left:50px;'>
          <el-form-item label="Name" prop="name">
            <el-input v-model="tempCreate.name" placeholder="只支持英文字母+数字,一旦输入不可更改!"></el-input>
          </el-form-item>

          <el-form-item label="描述">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="简单描述该Service" v-model="tempCreate.desc">
            </el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="createServiceVisible = false">取消</el-button>
          <el-button  type="primary" @click="createSubmit">创建</el-button>
        </div>
      </el-dialog>

       <el-dialog title="编辑Service" :visible.sync="editServiceVisible">
        <el-form :rules="editRules" :model="tempEdit" label-position="left" label-width="70px" style='width: 400px; margin-left:50px;'>
          <el-form-item label="Name" prop="name">
            <el-input v-model="tempEdit.name" :disabled="true"></el-input>
          </el-form-item>

          <!-- <el-form-item label="描述">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="简单描述该Service" v-model="tempEdit.desc">
            </el-input>
          </el-form-item> -->

          <el-form-item label="管理员">
             <el-tag
              :key="u"
              v-for="u in tempEdit.admins"
              closable
              :disable-transitions="false"
              @close="editAdminClose(u)">
              {{u}}
            </el-tag>
            <el-input
              class="input-new-tag"
              v-if="editAdminInputVisible"
              v-model="editAdminInputValue"
              ref="editAdminInput"
              size="small"
              @keyup.enter.native="adminInputConfirm"
              @blur="adminInputConfirm"
            >
               <el-button slot="append" size="small" @click="adminInputConfirm">ok</el-button>
            </el-input>
            <el-button v-else class="button-new-tag" size="small" @click="showEditAdminInput">+ New Tag</el-button>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="editServiceVisible = false">关闭</el-button>
        </div>
      </el-dialog>

      <el-dialog :title="'当前Service ：'+tempService.name" :visible.sync="metricsVisible" top="40px"  :fullscreen=true style="backgroud:#545c64 !important">
        <iframe :src="grafanaAddr(tempService.name)" style="height:910px;width:100%;border:none;margin-top:-40px;"></iframe>
      </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import waves from '@/directive/waves' // 水波纹指令
import {proxy} from '@/api/apiProxy' 
export default {
  name: 'service',
  directives: {
    waves
  },
  data() {
      return {
        list:[
        ],

        createServiceVisible: false,
        tempCreate: {
          name: ''
        },
        createRules: {
          name: [{ required: true, message: '请输入Service Name', trigger: 'blur' }]
        },

        tempEdit: {

        },
        editServiceVisible: false,
        editRules: {
        },
        editAdminInputVisible: false,
        editAdminInputValue: '',

        tempService: {},
        metricsVisible: false
      } 
  },
  methods : {
    grafanaAddr(serviceName) {
      var addr = process.env.SERVICE_GRAFANA + serviceName
       return  addr
    },
    handleMetrics(service) {
      this.tempService = service 
      this.metricsVisible = true
    },
     editAdminClose(u) {
          // 删除admin
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/privilege/delete',
                service: this.tempEdit.name,
                target_user: u,
                privilege: 1
            }
        proxy('POST',params).then(res => {
           if (res.data.status == 200) {
              this.tempEdit.admins.splice(this.tempEdit.admins.indexOf(u), 1);
              this.$message({
                    message: '删除管理成功',
                    type: 'success',
                    duration: 3 * 1000,
                    center: true
                })
           }
        })
       
      },
     adminInputConfirm() {
        // 添加admin
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/privilege/set',
                service: this.tempEdit.name,
                target_user: this.editAdminInputValue,
                privilege: 1
            }
        proxy('POST',params).then(res => {
           if (res.data.status == 200) {
              let inputValue = this.editAdminInputValue;
              if (inputValue) {
                this.tempEdit.admins.push(inputValue);
              }
              this.editAdminInputVisible = false;
              this.editAdminInputValue = '';
              this.$message({
                    message: '添加管理成功',
                    type: 'success',
                    duration: 3 * 1000,
                    center: true
                })
           }
        })

      },
    showEditAdminInput() {
      this.editAdminInputVisible = true
      this.$nextTick(_ => {
          this.$refs.editAdminInput.$refs.input.focus();
      });
    },
    handleEdit(row) {
      this.tempEdit = Object.assign({}, row) 
      // this.tempEdit.admins = []
      this.editServiceVisible = true
    },
    delService(name) {
      // 查询该service关联多少api，如果不等于0，不能删除
      var params = {
          target_service: 'tfe.manage',
          target_path: '/manage/service/countApi',
          service: name,
      }
      proxy('POST',params).then(res => {
          if (res.data.data != 0) {
             this.$message({
                  message: '当前service还有拥有API: '+res.data.data +'个，请先删除所有API后再删除Service',
                  type: 'warning',
                  duration: 5 * 1000,
                  center: true
              })
              return 
          }

          this.$confirm('此操作将永久删除<strong>该Service及关联API</strong>, 是否继续删除?', '提示', {
            dangerouslyUseHTMLString: true,
            cancelButtonText: '不，赶紧取消删除',
            confirmButtonText: 'ok',
            type: 'error'
          }).then(() => {
              // 获取用户信息
              var params = {
                      target_service: 'tfe.manage',
                      target_path: '/manage/service/delete',
                      name: name,
                  }
              proxy('POST',params).then(res => {
                if (res.data.status == 200) {
                    this.loadData()
                    this.$message({
                          message: '删除成功',
                          type: 'success',
                          duration: 3 * 1000,
                          center: true
                      })
                }
              })
          }).catch(() => {
            this.$message({
              type: 'info',
              message: '已取消删除'
            });          
          });
      })
    },
    handleCreate() {
        this.tempCreate = {
          name: ''
        }
        this.createServiceVisible = true
        this.$nextTick(() => {
          this.$refs['dataForm'].clearValidate()
        })
    },
    createSubmit() {
        // 获取用户信息
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/service/create',
                name: this.tempCreate.name,
                desc: this.tempCreate.desc
            }
        proxy('POST',params).then(res => {
           if (res.data.status == 200) {
              this.loadData()
              this.tempCreate = {
                name: ''
              }
              this.createServiceVisible = false
              this.$message({
                    message: '创建成功',
                    type: 'success',
                    duration: 3 * 1000,
                    center: true
                })
           }
        })
    },
    loadData() {
      // 加载该用户的所有service
        var params = {
                target_service: 'tfe.manage',
                target_path: '/manage/service/query',
            }
        proxy('POST',params).then(res => {
           if (res.data.status == 200) {
             this.list = res.data.data
           }
        })
    },
  },
  created() {
    this.loadData()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>

<style scoped>
 .el-tag + .el-tag {
    margin-left: 10px;
  }
  .button-new-tag {
    /* margin-left: 10px; */
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .input-new-tag {
    width: 120px;
    /* margin-left: 10px; */
    vertical-align: bottom;
  }
</style>


