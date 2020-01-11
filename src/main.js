import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import Vuex from 'vuex'
import VueClipboard from 'vue-clipboard2'
import VueShortkey from 'vue-shortkey'

import 'bootswatch/dist/flatly/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import {
  library,
} from '@fortawesome/fontawesome-svg-core'
import {
  fas,
} from '@fortawesome/free-solid-svg-icons'
import {
  FontAwesomeIcon,
} from '@fortawesome/vue-fontawesome'

import axios from 'axios'

import App from './App.vue'
import './style.scss'
import store from './store'

Vue.config.productionTip = false
Vue.use(BootstrapVue)
Vue.use(VueShortkey, {
  prevent: ['input', 'textarea'],
})
Vue.use(VueClipboard)

library.add(fas)
Vue.component('fa-icon', FontAwesomeIcon)

axios.defaults.baseURL = 'v2'

const go = new Go()
WebAssembly.instantiateStreaming(fetch('cryptocore.wasm'), go.importObject)
  .then(async obj => await go.run(obj.instance))

const instance = new Vue({
  store,
  render: h => h(App),
  mounted: () => store.dispatch('reload_users'),
}).$mount('#app')

// Wait for the cryptocore to be loaded (which makes encryption available)
new Promise(resolve => {
  (function waitForCryptocore() {
    if (window.opensslEncrypt) {
      return resolve()
    }
    setTimeout(waitForCryptocore, 100)
  }())
}).then(() => {
  store.commit('cryptocore_loaded')
})
