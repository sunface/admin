<template>
  <div class="login-container">

    <el-form class="login-form" autoComplete="on" label-position="left">
      <div class="title-container">
        <h3 class="title">传化网关</h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container svg-container_login">
          <svg-icon icon-class="user" />
        </span>
        <el-input name="username" type="text" v-model="loginForm.username" autoComplete="on" placeholder="输入账号"
        />
      </el-form-item>

      <el-form-item prop="password">
        <span class="svg-container">
          <svg-icon icon-class="password" />
        </span>
        <el-input name="password" :type="passwordType" v-model="loginForm.password" autoComplete="on"
          placeholder="输入密码" @keyup.enter.native="handleLogin" />
        <span class="show-pwd" @click="showPwd">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>

      <el-button type="primary" style="width:100%;margin-bottom:30px;"  @click.native.prevent="handleLogin">登录</el-button>

      <div class="tips">
        <span>访客账号 : guest</span>
        <span>密码 : guest</span>
      </div>
    </el-form>

    <div class="index-bg-footer"><img src="../../assets/login-bg2.png" style="width:100%"></div>
  </div>
</template>

<script>
/* eslint-disable */
const config = require('@/../config')
import request from '@/utils/request' 
export default {
  name: 'login',
  data() {
    return {
      passwordType: 'password',
      loginForm:{}
    }
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
    },
    handleLogin() {
      var _this = this
      // 登陆
      request({
          url: '/login',
          method: 'POST', 
          params: {
            user: this.loginForm.username,
            pw: this.loginForm.password
          }
      }).then(response => {
          _this.$store.dispatch('SetUserInfo', response.data.data).then(() => {  
              _this.$router.push({ path: '/adminui' })
          })
      }).catch(error => {
      })
    }
  },
  created() {
 
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
$bg:#2d3a4b;
$light_gray:#eee;

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;
    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      height: 40px;
      &:-webkit-autofill {
        // -webkit-box-shadow: 0 0 0px 1000px $bg inset !important;
        // -webkit-text-fill-color: #fff !important;
      }
    }
  }
  .el-form-item {
    border-bottom: 1px solid rgba(255, 255, 255, 0.9);
    // background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    // color: #454545;
  }
}
</style>

<style rel="stylesheet/scss" lang="scss" scoped>
$dark_gray:#889aa4;
$light_gray:#eee;

.login-container {
  width: 100%;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  background-color: #fff;
  background-image: url(../../assets/login-bg1.jpg);
  background-repeat: no-repeat;
  background-size: cover;
  background-position: 50%;
  overflow: hidden;
  .login-form {
    position: absolute;
    left: 50px;
    width: 460px;
    padding: 35px 35px 15px 35px;
    margin-top: 150px;
    z-index: 99;
  }
  .tips {
    font-size: 13px;
    color: rgb(180,180,180);
    margin-bottom: 10px;
    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }
  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
    &_login {
      font-size: 20px;
    }
  }
  .title-container {
    position: relative;
    .title {
      font-size: 36px;
      font-weight: bold;
    }
    .title::after {
        content: " ";
        border-bottom-style: solid;
        border-bottom-width: 3px;
        display: block;
        width: 35px;
        color: #2d8cf0;
        padding-top: 26px;
        margin-left: 6px
    }
    .set-language {
      color: #fff;
      position: absolute;
      top: 5px;
      right: 0px;
    }
  }
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
  .thirdparty-button {
    position: absolute;
    right: 35px;
    bottom: 28px;
  }
}

.index-bg-footer {
  width: 260px;
  height: 90px;
  position: absolute;
  bottom: 0;
  left: 0;
  img {
    z-index: 1
  }
}

 @media only screen and (max-width: 800px) {
     .login-form { 
        margin-top: 30px !important;
     }
     .index-bg-footer {
       display: none !important
     }
}
</style>
