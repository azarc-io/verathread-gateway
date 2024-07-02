<template>
  <div class="home">
    <div v-if="isLoadingComponent">Local View</div>
<!--    <HelloWorld msg="Welcome to Your Vue.js + TypeScript App" />-->
  </div>
</template>

<script lang="ts" setup>
import { defineAsyncComponent, ref } from 'vue';

const isLoadingComponent = ref(true);

// const HelloWorld = defineAsyncComponent(() =>
//   import('app_exposes/HelloWorld.vue').finally(() => {
//     isLoadingComponent.value = false;
//   }),
// );
</script>
