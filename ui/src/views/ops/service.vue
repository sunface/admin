<template>
  <div class="app-container">
      <div class="filter-container">
        <!-- <el-input  style="width: 200px;" class="filter-item" placeholder="Service">
        </el-input>
        <el-button class="filter-item" type="primary" style="margin-left: 10px;" v-waves icon="el-icon-search" @click="handleCreate">搜索</el-button> -->
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"  v-waves type="success" icon="el-icon-plus">Add Service</el-button>
        <!-- <el-button class="filter-item" type="primary" :loading="downloadLoading" v-waves icon="el-icon-download" @click="handleDownload">导出</el-button> -->
      </div>

      <div class="table">
        <el-table  :data="services" border fit highlight-current-row style="width: 100%"  :default-sort = "{prop: 'modify_date', order: 'descending'}">
          <el-table-column align="center" label="Name" width="250" prop="name" sortable>
            <template slot-scope="scope">
              <span>{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column width="200" align="center" label="Owner" sortable>
             <template slot-scope="scope">
              <span>{{scope.row.owner}}</span>
            </template>
          </el-table-column>
          <el-table-column width="350" align="center" label="Description">
            <template slot-scope="scope">
              <span>{{scope.row.desc}}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="Operate" class-name="small-padding">
            <template slot-scope="scope">
              <el-button size="mini" @click="handleEdit(scope.row)">Edit</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>


      <el-dialog title="Create service" :visible.sync="createServiceVisible">
        <el-form :model="tempCreate" label-position="left" label-width="70px" style='width: 500px; margin-left:50px;'>
          <el-form-item label="name" prop="name">
            <el-input v-model="tempCreate.name" placeholder="only support alhpabet,once created,it cant be changed"></el-input>
          </el-form-item>
          <el-form-item label="owner" style="margin-top:25px">
            <el-select class="filter-item" v-model="tempCreate.owner" style="width: 200px"  placeholder="select user" filterable>
                <el-option v-for="u in users" :key="u.id" :label="u.username" :value="u.username"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="desc">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="describe the service" v-model="tempCreate.desc">
            </el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="createServiceVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="createSubmit">Submit</el-button>
        </div>
      </el-dialog>

       <el-dialog title="Edit Service" :visible.sync="editServiceVisible">
        <el-form :model="tempEdit" label-position="left" label-width="120px" style='width: 550px; margin-left:50px;'>
          <el-form-item label="Name" prop="name">
              {{tempEdit.name}}
          </el-form-item>
          <el-form-item label="desc">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="describe the service" v-model="tempEdit.desc">
            </el-input>
          </el-form-item>
          <!-- <el-form-item label="描述">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="简单描述该Service" v-model="tempEdit.desc">
            </el-input>
          </el-form-item> -->

          <el-form-item label="set user priv">
            <el-select class="filter-item" v-model="editServiceUser" style="width: 150px"  placeholder="select user" filterable>
                <el-option v-for="u in users" :key="u.id" :label="u.username" :value="u.username"></el-option>
            </el-select>
             <el-select class="filter-item" v-model="editServicePriv" style="width: 120px"  placeholder="select priv">
              <el-option label="normal" value="normal"></el-option>
              <el-option label="admin" value="admin"></el-option>
            </el-select>
            <el-button  class="button-new-tag" icon="el-icon-plus" @click="addServiceUser" style="border:none"></el-button>

            <div>
                <el-tag :key="u" v-for="u in editServiceUsers" closable :disable-transitions="false" @close="delServiceUser(u.username)">{{u.username}}</el-tag>
            </div>
          </el-form-item>
          
          <div class="form-block">
             <el-tag type="warning" style="border:none">危险区域</el-tag>
             <el-form-item label="change owner" class="margin-top-20">
                <el-select class="filter-item" v-model="tempEdit.owner" style="width: 150px"  placeholder="select user" filterable>
                    <el-option v-for="u in users" :key="u.id" :label="u.username" :value="u.username"></el-option>
                </el-select>
                 <el-alert
                    :closable=false
                    title="Once the owner is changed to others, you will loose the control of this service"
                    type="warning">
                </el-alert> 
            </el-form-item>
          </div>

        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="editServiceVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="editSubmit">Submit</el-button>
        </div>
      </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import waves from '@/directive/waves' // 水波纹指令
import {proxy} from '@/api/apiProxy'
import {loadUsers} from '@/api/user'
import request from '@/utils/request' 
export default {
  name: 'service',
  directives: {
    waves
  },
  data() {
      return {
        services:[],
        users: [],
        
        createServiceVisible: false,
        tempCreate: {
          name: ''
        },

        tempEdit: {

        },
        editServiceVisible: false,

        tempService: {},
        metricsVisible: false,

        editServiceUser:'',
        editServicePriv: '',

        editServiceUsers: []
      } 
  },
  methods : {
    delServiceUser(u) {
        this.$confirm('Are you sure to remove the user : ' + u +'?', 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'cancel',
          confirmButtonText: 'submit',
          type: 'warning'
        }).then(() => {
            request({
                url: '/ops/service/delUser',
                method: 'POST', 
                params: {
                    service: this.tempEdit.name,
                    username: u
                }
            }).then(res => {
                for (var i=0; i<this.editServiceUsers.length;i++) {
                    if (this.editServiceUsers[i].username == u) {
                        this.editServiceUsers.splice(i,1)
                    }
                }
            })
        });
    },
    addServiceUser() {
        request({
          url: '/ops/service/addUser',
          method: 'POST', 
          params: {
            service: this.tempEdit.name,
            username: this.editServiceUser,
            priv: this.editServicePriv
          }
        }).then(res => {
           this.editServiceUsers.push({
               username: this.editServiceUser,
               priv: this.editServicePriv
           })
           this.editServiceUser = ''
           this.editServicePriv = ''
        })
    },
    handleEdit(row) {
      this.tempEdit = Object.assign({}, row) 
      // this.tempEdit.admins = []
      this.loadServiceUsers()
      this.editServiceVisible = true
    },
    handleCreate() {
        this.tempCreate = {
        }
        this.createServiceVisible = true
    },
    editSubmit() {
        request({
          url: '/ops/service/update',
          method: 'POST', 
          params: {
            service: this.tempEdit.name,
            desc: this.tempEdit.desc,
            owner: this.tempEdit.owner,
          }
        }).then(res => {
            this.loadData()
            this.tempEdit = {
            }
            this.editServiceVisible = false
            this.$message({
                message: '编辑成功',
                type: 'success',
                duration: 3 * 1000,
                center: true
            })
        })
    },
    createSubmit() {
        request({
          url: '/ops/service/create',
          method: 'POST', 
          params: {
            service: this.tempCreate.name,
            desc: this.tempCreate.desc,
            owner: this.tempCreate.owner,
          }
        }).then(res => {
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
        })
    },
    loadData() {
       request({
          url: '/ops/service/query',
          method: 'GET'
        }).then(res => {
            this.services = res.data.data
        })
    },
    loadServiceUsers() {
        request({
          url: '/ops/service/queryUser',
          method: 'GET',
          params: {
              service: this.tempEdit.name
          }
        }).then(res => {
            this.editServiceUsers = res.data.data
        })
    }
  },
  mounted() {
    this.loadData()
    loadUsers().then(res => {
        this.users = res.data.data
    })
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


