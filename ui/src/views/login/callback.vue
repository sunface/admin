<template>
  <div class="callback">
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
export default {
  name: 'callback',
  data() {
      return {}
  },
  created() {
    var subToken = this.$route.query.subToken 
    var _this = this
    // 获取用户信息s
        request({
            url: '/login',
            method: 'POST', 
            params: {
                subToken: subToken,
            }
        }).then(response => {
            // 存在，设置用户信息
            _this.$store.dispatch('SetUserInfo', response.data.data).then(() => {  
                _this.$router.push({ path: '/' })
            })
        }).catch(error => {
          _this.$router.push({ path: '/' })
        })
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>
