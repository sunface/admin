<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate" type="success" icon="el-icon-edit">{{$t('ops.addUser')}}</el-button>
    </div>

    <div class="table">
      <el-table  :data="users" border fit highlight-current-row style="width: 100%"  :default-sort = "{prop: 'create_date', order: 'ascending'}">
        <el-table-column align="center" :label="$t('ops.userName')" width="250" prop="name" sortable>
          <template slot-scope="scope">
            <span>{{scope.row.username}}</span>
          </template>
        </el-table-column>
        <el-table-column width="250" align="center" :label="$t('common.privilege')" prop="creator" sortable>
            <template slot-scope="scope">
            <span>{{scope.row.priv}}</span>
          </template>
        </el-table-column>
        <el-table-column width="250" align="center" :label="$t('common.createDate')" prop="admins">
          <template slot-scope="scope">
            <span>{{scope.row.create_date}}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" :label="$t('common.operate')" class-name="small-padding">
          <template slot-scope="scope">
            <el-button  size="mini" @click="editUser(scope.row)" v-if="scope.row.username!='admin'">{{$t('common.edit')}}</el-button>
            <el-button  size="mini" type="danger" @click="deleteUser(scope.row.username)" v-if="scope.row.username!='admin'">{{$t('common.delete')}}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>


    <el-dialog  class="mf-dialog"  title="Add User" :visible.sync="createVisible">
        <el-form  label-position="left" label-width="100px" style='width: 400px; margin-left:50px;'>
          <el-form-item :label="$t('ops.userName')">
            <el-input v-model="tempCreate.username" :placeholder="$t('common.onlyNumAndAlpha')"></el-input>
          </el-form-item>

          <el-form-item :label="$t('common.privilege')">
            <el-select class="filter-item" v-model="tempCreate.priv" style="width: 200px"  :placeholder="$t('common.selectPriv')">
              <el-option label="normal" value="normal"></el-option>
              <el-option label="admin" value="admin"></el-option>
            </el-select>
          </el-form-item>
           <el-alert
            :closable=false
            title="Only username 'admin' can operate"
            type="warning">
          </el-alert> 
        </el-form>
       
        <div slot="footer" class="dialog-footer">
          <el-button @click="createVisible = false">{{$t('common.cancel')}}</el-button>
          <el-button  type="primary" @click="createSubmit">{{$t('common.submit')}}</el-button>
        </div>
      </el-dialog>
    
   <el-dialog class="mf-dialog"  title="Edit User" :visible.sync="editVisible">
        <el-form  label-position="left" label-width="100px" style='width: 400px; margin-left:50px;'>
          <el-form-item label="user name">
             {{tempEdit.username}}
          </el-form-item>

          <el-form-item label="privilege">
            <el-select class="filter-item" v-model="tempEdit.priv" style="width: 200px"  placeholder="select priv" :disabled="tempEdit.username=='admin'">
              <el-option label="admin" value="admin"></el-option>
              <el-option label="normal" value="normal"></el-option>
            </el-select>
          </el-form-item>
           <el-alert
            :closable=false
            title="Only username 'admin' can operate"
            type="warning">
          </el-alert> 
        </el-form>
       
        <div slot="footer" class="dialog-footer">
          <el-button @click="editVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="editSubmit">Submit</el-button>
        </div>
    </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
import {loadUsers} from '@/api/user'
export default {
  name: 'users',
  data() {
      return {
        createVisible: false,
        tempCreate: {},

        users: [],

        editVisible: false,
        tempEdit :{}
      }
  },
  methods: {
    editSubmit() {
      request({
          url: '/user/edit',
          method: 'POST', 
          params: {
            username: this.tempEdit.username,
            priv: this.tempEdit.priv
          }
      }).then(res => {
        this.$message({
          message: 'edit user ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        }) 
        this.editVisible = false
        loadUsers().then(res => {
          this.users = res.data.data
        })
      })
    },
    editUser(user) {
      this.tempEdit = user
      this.editVisible = true
    },
    handleCreate() {
      this.createVisible = true
    },
    createSubmit() {
      request({
          url: '/user/add',
          method: 'POST', 
          params: {
            username: this.tempCreate.username,
            priv: this.tempCreate.priv
          }
      }).then(res => {
        this.$message({
          message: 'add user ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        }) 
        this.tempCreate = {}
        this.createVisible = false
        loadUsers().then(res => {
          this.users = res.data.data
        })
      })
    },
    deleteUser(username) {
        this.$confirm('Are you sure to delete the user : ' + username +'?', 'Warning', {
          dangerouslyUseHTMLString: true,
          cancelButtonText: 'cancel',
          confirmButtonText: 'submit',
          type: 'warning'
        }).then(() => {
          request({
              url: '/user/delete',
              method: 'POST', 
              params: {
                username: username
              }
          }).then(res => {
           loadUsers().then(res => {
              this.users = res.data.data
            })
            this.$message({
              message: 'delete user ok',
              type: 'success',
              duration: 3 * 1000,
              center: true
            }) 
          })
        });
    },
  },
  mounted() {
    // load users
    loadUsers().then(res => {
      this.users = res.data.data
    })
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
