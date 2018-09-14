<template>
  <div class="app-container">
    <div class="filter-container">
        <el-tag>service</el-tag>
        <el-select clearable class="filter-item" :value="calcService()" @change='handleSelService' style="width: 200px"  placeholder="请选择Service">
            <el-option v-for="s in  services" :key="s.name" :label="s.name" :value="s.name"></el-option>
        </el-select>
        <el-button class="filter-item" style="margin-left: 10px;" @click="handleCreate"  v-waves type="success" icon="el-icon-plus">Add App</el-button>
    </div>

    <div class="table">
        <el-table  :data="apps" border fit highlight-current-row style="width: 100%;"  :default-sort = "{prop: 'modify_date', order: 'descending'}">
            <el-table-column align="center" label="Name" width="250" prop="name" sortable>
            <template slot-scope="scope">
                <span>{{scope.row.name}}</span> 
            </template>
            </el-table-column>
            <el-table-column width="200" align="center" label="Service" sortable>
                <template slot-scope="scope">
                <span>{{scope.row.service}}</span>
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

     <el-dialog title="Create application" :visible.sync="createVisible">
        <el-form  label-position="left" label-width="70px" style='width: 500px; margin-left:50px;'>
          <el-form-item label="name" prop="name">
            <el-input v-model="tempCreate.name" placeholder="only support alhpabet,once created,it cant be changed"></el-input>
          </el-form-item>
          <el-form-item label="service">
            {{selectedService}}
          </el-form-item>
          <el-form-item label="desc">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="describe the app" v-model="tempCreate.desc">
            </el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="createVisible = false">Cancel</el-button>
          <el-button  type="primary" @click="createSubmit">Submit</el-button>
        </div>
      </el-dialog>

      <el-dialog title="Edit application" :visible.sync="editVisible">
        <el-form  label-position="left" label-width="70px" style='width: 500px; margin-left:50px;'>
          <el-form-item label="name" prop="name">
            {{tempEdit.name}}
          </el-form-item>
          <el-form-item label="service">
            {{tempEdit.service}}
          </el-form-item>
          <el-form-item label="desc">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 4}" placeholder="describe the app" v-model="tempEdit.desc">
            </el-input>
          </el-form-item>
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
import waves from '@/directive/waves' // 水波纹指令
export default {
  name: 'application',
  data() {
      return {
          services :[],
          apps: [],
          selectedService: '',

          createVisible: false,
          tempCreate: {},

          editVisible: false,
          tempEdit: {}
      }
  },
  directives: {
    waves
  },
  methods: {
    editSubmit() {
        request({
          url: '/ops/app/update',
          method: 'POST', 
          params: {
            app: this.tempEdit.name,
            service: this.tempEdit.service,
            desc: this.tempEdit.desc
          }
        }).then(res => {
            this.editVisible = false
            this.loadApps()
            this.$message({
                message: 'eidt ok',
                type: 'success',
                duration: 3 * 1000,
                center: true
            })
        })
    },
    handleEdit(app) {
        this.tempEdit = app 
        this.editVisible = true
    },
    createSubmit() {
        request({
          url: '/ops/app/create',
          method: 'POST', 
          params: {
            app: this.tempCreate.name,
            service: this.selectedService,
            desc: this.tempCreate.desc
          }
        }).then(res => {
            this.tempCreate = {}
            this.createVisible = false
            this.loadApps()
            this.$message({
                message: '创建成功',
                type: 'success',
                duration: 3 * 1000,
                center: true
            })
        })
    },
    handleCreate() {
        this.createVisible = true
    },
    handleSelService(service) {
        this.$store.dispatch('setService', service)
        this.selectedService = service
        this.loadApps()
    },
    calcService() {
        return this.selectedService || this.$store.getters.service
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
    }
  },
  mounted() {
      this.loadServices()
      this.selectedService = this.$store.getters.service
      if (this.selectedService != '') {
          this.loadApps()
      }
    //   
    //   alert()
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
