<template>
  <div class="app-container">
    <div class="filter-container">
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"   type="success" icon="el-icon-edit">Create Server</el-button>
    </div>

    <div class="table">
        <el-table  :data="servers" border fit highlight-current-row style="width: 100%" size="small" :default-sort = "{prop: 'modify_date', order: 'descending'}">
            <el-table-column align="center" label="IP" width="250" prop="name" sortable>
                <template slot-scope="scope">
                    <div><el-tag size="mini" type="success" style="border:none;margin-right:5px">private</el-tag>{{scope.row.private_ip}}</div>
                    <div><el-tag size="mini" style="border:none;margin-right:5px">public</el-tag>{{scope.row.public_ip}}</div>
                </template>
            </el-table-column>
            <el-table-column width="150" align="center" label="Cloud Provider">
                <template slot-scope="scope">
                    <span>{{scope.row.cloud_provider}}</span>
                </template>
            </el-table-column>
            <el-table-column width="150" align="center" label="Region">
                <template slot-scope="scope">
                    <span>{{scope.row.region}}</span>
                </template>
            </el-table-column>
             <el-table-column width="100" align="center" label="Ssh User">
                <template slot-scope="scope">
                    <span>{{scope.row.ssh_user}}</span>
                </template>
            </el-table-column>
            <el-table-column width="150" align="center" label="Ssh PW">
                <template slot-scope="scope">
                    <span v-if="scope.row.ssh_pw != ''">{{scope.row.ssh_pw}}</span>
                    <span v-else> <i class="el-icon-view hover-cursor icon-primary"  @click="queryPW(scope.row)"></i> </span>
                </template>
            </el-table-column>
            <el-table-column width="250" align="center" label="Configure">
                <template slot-scope="scope">
                    <span>{{scope.row.configure}}</span>
                </template>
            </el-table-column>
             <el-table-column width="250" align="center" label="Expire">
                <template slot-scope="scope">
                    <span>{{scope.row.expre}}</span>
                </template>
            </el-table-column>
            <el-table-column align="center" label="Operate" class-name="small-padding">
                <template slot-scope="scope">
                    <el-button size="mini" @click="handleEdit(scope.row)">Edit</el-button>
                </template>
            </el-table-column>
        </el-table>
    </div>

    <el-dialog :title="defineTitle()" :visible.sync="defineVisible">
        <el-form  label-position="left" label-width="130px" style='width: 550px; margin-left:50px;'>
          <el-form-item label="public ip">
            <el-input v-model="tempServer.public_ip" placeholder="a.b.c.d"></el-input>
          </el-form-item>
          <el-form-item label="private ip">
            <el-input v-model="tempServer.private_ip" placeholder="a.b.c.d"></el-input>
          </el-form-item>
           <el-form-item label="cloud provider">
             <el-select  class="filter-item" v-model="tempServer.cloud_provider" style="width: 200px;margin-top:8px"  placeholder="select cloud">
                <el-option label="private" value="private"></el-option>
                <el-option  label="aliyun" value="aliyun"></el-option>
                <el-option  label="aws" value="aws"></el-option>
                <el-option  label="gce" value="gce"></el-option>
                <el-option  label="azure" value="azure"></el-option>
                <el-option  label="tencent" value="tencent"></el-option>
                <el-option  label="other" value="other"></el-option>
             </el-select>
          </el-form-item>
          <el-form-item label="ssh user">
            <el-input v-model="tempServer.ssh_user" placeholder="e.g. root"></el-input>
          </el-form-item>
          <el-form-item label="ssh password">
              <i class="el-icon-view hover-cursor icon-primary" style="color: #67c23a" @click="queryPW(tempServer)" v-if="tempServer.ssh_pw==''"></i>
              <el-input v-else v-model="tempServer.ssh_pw"></el-input>
          </el-form-item>
          <el-form-item label="region">
            <el-input v-model="tempServer.region" placeholder="HongKong.B"></el-input>
          </el-form-item>
          <el-form-item label="configure">
            <el-input v-model="tempServer.configure" placeholder="2C.4G.50G"></el-input>
          </el-form-item>
          <el-form-item label="expire date">
            <el-date-picker
                v-model="tempServer.expire"
                type="date"
                placeholder="pick date">
            </el-date-picker>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="defineVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="defineSubmit">Submit</el-button>
        </div>
      </el-dialog>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
export default {
  name: 'server',
  data() {
      return {
          servers: [],

          defineVisible: false,
          tempServer: {},
          defineStatus: 'create',

          tempPW : ''
      }
  },
  methods: {
      handleEdit(server) {
          this.tempServer = server
          this.tempServer.ssh_pw = ''
          this.defineStatus = 'edit'
          this.defineVisible = true
      },
      queryPW(server) {
           request({
          url: '/ops/server/queryPW',
          method: 'GET', 
          params: {
              id: server.id
          }
        }).then(res => {
            server.ssh_pw = res.data.data
        })
      },
      defineSubmit() {
        request({
          url: '/ops/server/define',
          method: 'POST', 
          params: {
            server: this.tempServer,
            op : this.defineStatus
          }
        }).then(res => {
            this.defineVisible = false
            this.loadServers()
            this.$message({
                message: 'create ok',
                type: 'success',
                duration: 3 * 1000,
                center: true
            })
        })
      },
      handleCreate() {
          this.defineStatus = 'create'
          this.defineVisible = true
      },
      defineTitle() {
          if (this.defineStatus == 'create') {
              return 'Create Server'
          } else {
              return 'Edit Server'
          }
      },
      loadServers() {
        request({
          url: '/ops/server/query',
          method: 'GET', 
          params: {
          }
        }).then(res => {
            this.servers = res.data.data
        })
      }
  },
  mounted() {
    this.loadServers()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
