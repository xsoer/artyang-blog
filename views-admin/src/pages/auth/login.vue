<template>
    <div>
        <h3>{{msg}}</h3>
        <input aria-label="姓名">
    </div>
</template>
<script>
export default {
  name: 'Login',
  data () {
    return {
      msg: '登陆'
    }
  }
}
</script>
<style scoped>
</style>
