<script lang="ts">
import {computed, defineComponent, ref} from 'vue';

export default defineComponent({
  setup() {
    const number = ref(0);

    function increment() {
      number.value++;
    }

    return { number, increment };
  },
})
</script>

<template>
  <div class="content">
    <button @click="increment">
      Remote Counter {{ number }}
    </button>
  </div>
</template>

<style lang="scss" scoped>
.content {
  display: flex;
  line-height: 1.1;
  text-align: center;
  flex-direction: column;
  justify-content: center;
}
</style>

