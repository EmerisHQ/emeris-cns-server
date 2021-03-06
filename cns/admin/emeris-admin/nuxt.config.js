export default {
  // Target: https://go.nuxtjs.dev/config-target
  target: 'static',

  generate: {
    exclude: [
      /^\/chain/ // path starts with /chain
    ]
  },

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    title: 'demeris-admin',
    htmlAttrs: {
      lang: 'en'
    },
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/admin/favicon.ico' },
      { rel: 'dns-prefetch', href: 'https://fonts.gstatic.com' },
      {
        rel: 'stylesheet',
        type: 'text/css',
        href: 'https://fonts.googleapis.com/css?family=Nunito',
      },
      {
        rel: 'stylesheet',
        type: 'text/css',
        href:
          'https://cdn.materialdesignicons.com/4.9.95/css/materialdesignicons.min.css',
      },
    ]
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: ['./assets/scss/main.scss'],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [{ src: '~/plugins/after-each.js', mode: 'client' }],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: false,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    // Doc: https://buefy.github.io/#/documentation
    ['nuxt-buefy', { materialDesignIcons: false }],
    'bootstrap-vue/nuxt',
    '@nuxtjs/axios',
    [
      '@nuxtjs/firebase',
      {
        config: {
          apiKey: "AIzaSyA09i5TB-Jb0-llyRaHPjytr9iFNY8V1TI",
          authDomain: "emeris-admin-ui.firebaseapp.com",
          projectId: "emeris-admin-ui",
          storageBucket: "emeris-admin-ui.appspot.com",
          messagingSenderId: "456830583626",
          appId: "1:456830583626:web:cc57b5b475143771f177b3"
        
        },
        services: {
          auth: {
            persistence: 'local',
            initialize: {
              onAuthStateChangedAction: 'onAuthStateChangedAction',
              subscribeManually: false
            },
            ssr: false,
          }
        }
      }
    ],
  ],

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {
  },

  axios: {
    baseUrl: process.env.CNS_URL || "/v1/cns",
    apiUrl: process.env.API_URL || "/v1"
  },

  router: {
    base: process.env.BASE_URL || "/admin",
    middleware: 'auth'
  }
};
