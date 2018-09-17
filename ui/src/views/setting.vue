<template>
  <div class="app-container">
        <el-form  label-position="left" label-width="120px" style="width: 650px;margin-left:15px" size="mini">
            <div class="form-block" style="padding:15px 15px">
              <el-tag  size="small" style="border:none;font-size:14px">My pw</el-tag>
              <el-form-item label="old" style="margin-top:25px">
                <el-input  v-model="oldPw" placeholder="input the old pw" style="width:250px">
                </el-input>
                </el-form-item>
                    <el-form-item label="new">
                    <el-input  v-model="newPW1" placeholder="input the new pw" style="width:250px">
                    </el-input>
                </el-form-item>
                <el-form-item label="re-new">
                <el-input v-model="newPW2" placeholder="reinput the new pw" style="width:250px">
                </el-input>
              </el-form-item>
              <el-button  type="success" size="small" @click="setPW">Submit</el-button>
            </div>


            <div class="form-block" style="padding:15px 15px" v-if="priv=='admin' || priv=='super_admin'">
              <el-tag  size="small" style="border:none;font-size:14px">User's pw</el-tag>
              <el-form-item label="username" style="margin-top:25px">
                <el-select class="filter-item" v-model="resetUsername" style="width: 200px"  placeholder="select user" filterable>
                  <el-option v-for="u in users" :key="u.id" :label="u.username" :value="u.username"></el-option>
                </el-select>
                </el-form-item>
                    <el-form-item label="new">
                    <el-input  v-model="resetPW1" placeholder="input the new pw" style="width:250px">
                    </el-input>
                </el-form-item>
                <el-form-item label="re-new">
                <el-input v-model="resetPW2" placeholder="reinput the new pw" style="width:250px">
                </el-input>
              </el-form-item>
              <el-button  type="success" size="small" @click="setUserPW">Submit</el-button>
            </div>
        </el-form>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
import {getLoginInfo} from '@/api/login'
import {loadUsers} from '@/api/user'
export default {
  name: 'setting',
  data() {
      return {
          users: [],

          oldPw : '',
          newPW1 : '',
          newPW2 : '',
          priv: 'normal',

          resetUsername: '',
          resetPW1: '',
          resetPW2: ''
      }
  },
  methods: {
    setUserPW() {
       if (this.resetPW1 == '' || this.resetPW2 == '' || this.resetUsername == '') {
        this.$message({
          message: 'required fields must not be empty',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }

      if (this.resetPW1 != this.resetPW2) {
        this.$message({
          message: 'the two new pw is not same',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }

      request({
          url: '/user/setUserPW',
          method: 'POST', 
          params: {
            username : this.resetUsername,
            new_pw: this.resetPW1
          }
      }).then(res => {
        this.$message({
          message: 'setting pw ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        }) 
      })
    },
    setPW() {
      if (this.oldPw == '' || this.newPW1 == '' || this.newPW2 == '') {
        this.$message({
          message: 'required fields must not be empty',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }

      if (this.newPW1 != this.newPW2) {
        this.$message({
          message: 'the two new pw is not same',
          type: 'warning',
          duration: 3 * 1000,
          center: true
        }) 
        return 
      }

      request({
          url: '/user/setPW',
          method: 'POST', 
          params: {
            old_pw: this.oldPw,
            new_pw: this.newPW1
          }
      }).then(res => {
        this.$message({
          message: 'setting pw ok',
          type: 'success',
          duration: 3 * 1000,
          center: true
        }) 
        // re-login
        this.$store.dispatch('LogOut').then(res => {
          this.$router.push('/adminui/login')
        }) 
      })
    }
  },
  mounted() {
     getLoginInfo().then((res) => {
        this.priv = res.data.data.priv
     })
     loadUsers().then(res => {
        this.users = res.data.data
      })
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
