<template>
  <div id="app">
    <aside-menu :menu="menu" />
    <nav-bar />
    <nuxt />
    <footer-bar />
  </div>
</template>

<script>
import AsideMenu from '@/components/AsideMenu'
import FooterBar from '@/components/FooterBar'
import NavBar from '@/components/NavBar'

export default {
  name: 'Home',
  components: {
    FooterBar,
    NavBar,
    AsideMenu,
  },
  computed: {
    menu() {
      return [
        'General',
        [
          {
            to: '/',
            icon: 'desktop-mac',
            label: 'Dashboard',
          },
          {
            to: '/chains',
            label: 'Chains',
            icon: 'table',
          },
        ],
      ]
    },
  },
  mounted() {
    document.documentElement.classList.add('has-aside-left')
    document.documentElement.classList.add('has-aside-mobile-transition')
    document.documentElement.classList.add('has-aside-expanded')
  },
}
</script>
