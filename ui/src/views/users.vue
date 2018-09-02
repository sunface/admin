<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate" type="primary" icon="el-icon-edit">Add User</el-button>
    </div>

    <div class="table">
      <el-table  :data="users" border fit highlight-current-row style="width: 100%;min-height:1000px;"  :default-sort = "{prop: 'create_date', order: 'ascending'}">
        <el-table-column align="center" label="user name" width="250" prop="name" sortable>
          <template slot-scope="scope">
            <span>{{scope.row.username}}</span>
          </template>
        </el-table-column>
        <el-table-column width="250" align="center" label="privilege" prop="creator" sortable>
            <template slot-scope="scope">
            <span>{{scope.row.priv}}</span>
          </template>
        </el-table-column>
        <el-table-column width="250" align="center" label="create date" prop="admins">
          <template slot-scope="scope">
            <span>{{scope.row.create_date}}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" class-name="small-padding">
          <template slot-scope="scope">
            <el-button  size="mini" @click="editUser(scope.row)" v-if="scope.row.username!='admin'">Edit</el-button>
            <el-button  size="mini" type="danger" @click="deleteUser(scope.row.username)" v-if="scope.row.username!='admin'">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>


    <el-dialog title="Add User" :visible.sync="createVisible">
        <el-form  label-position="left" label-width="100px" style='width: 400px; margin-left:50px;'>
          <el-form-item label="user name">
            <el-input v-model="tempCreate.username" placeholder="only support alphabet and numeric"></el-input>
          </el-form-item>

          <el-form-item label="privilege">
            <el-select class="filter-item" v-model="tempCreate.priv" style="width: 200px"  placeholder="select priv">
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
          <el-button @click="createVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="createSubmit">Submit</el-button>
        </div>
      </el-dialog>
    
   <el-dialog title="Edit User" :visible.sync="editVisible">
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
