<script lang="ts" >
import {defineComponent, inject, onMounted, ref, watch} from "vue";
import {State} from "./store/store";

export default defineComponent({
  setup() {
    const state = inject<State>('state')

    return {
      configuration: state.configuration,
    }
  },
  beforeMount() {
    console.log('beforeMount')
  }
})
</script>

<template>
  <nav class="navbar">
    <ul v-if="configuration" class="navbar-nav">
      <li class="nav-item">
        <router-link to="/">Home</router-link>
      </li>
      <li v-for="cat of configuration.categories" :key="cat.category" class="nav-item has-dropdown">
        <a href="#">{{ cat.title }}</a>
        <ul v-if="cat && cat.entries" class="dropdown">
          <li v-for="ent of cat.entries" :key="ent.title" class="dropdown-item">
            <router-link :to="ent.module.path">{{ent.title}}</router-link>
          </li>
        </ul>
      </li>
    </ul>
  </nav>

  <div class="content">
    <router-view />
  </div>
</template>

<style scoped lang="scss">
.content {
  display: flex;
  min-height: calc(100vh - 70px);
  line-height: 1.1;
  text-align: center;
  flex-direction: column;
  justify-content: center;
}

.navbar {
  height: 70px;
  width: 100%;
  background: black;
  color: white;
}

.navbar-nav {
  list-style-type: none;
  margin: 0;
  padding: 0;

  display: flex;
  align-items: center;
  justify-content: space-evenly;
  height: 100%;
}

.dropdown {
  opacity: 0;
  position: absolute;
  width: 500px;
  z-index: 2;
  background: black;

  display: flex;
  align-items: center;
  justify-content: space-around;
  height: 3rem;
  margin-top: 2rem;
  padding: 0.5rem;

  transform: translateX(-40%);
  transition: opacity .15s ease-out;
}

.dropdown-item a {
  width: 100%;
  height: 100%;
}

.has-dropdown:focus-within .dropdown   {
  opacity: 1;
  pointer-events: auto;
}
</style>
