<script setup lang="ts">const modules = ["房间生命周期", "并发加入", "房主转移", "SSE 事件和超时解散"]</script>
<template><main><h1>多人游戏房间与匹配服务</h1><ul><li v-for="item in modules" :key="item">{{ item }}</li></ul></main></template>
<style>body{font-family:system-ui;background:#f5f7fb}main{max-width:900px;margin:3rem auto;background:white;padding:2rem;border-radius:16px}</style>
