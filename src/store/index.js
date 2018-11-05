import Vue from 'vue'
import Vuex from 'vuex'

import axios from 'axios'
import uuidv4 from 'uuid/v4'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    //account_info: {
    //    data: [{
    //            "title": "Test entry",
    //            "username": "testuser",
    //            "password": "quitesecretpass",
    //            "url": "https://example.com",
    //            "comment": "",
    //            "tags": "",
    //            "id": "f106cdd5-7b52-4c51-a386-6f2016ee70c8",
    //        },
    //        {
    //            "title": "Test entry 2",
    //            "username": "testuser",
    //            "password": "quitesecretpass",
    //            "url": "https://example.com",
    //            "comment": "",
    //            "tags": "foobar",
    //            "id": "e956814f-c7dd-4730-b383-624f8bbc6923",
    //        },
    //    ],
    //    data_raw: "",
    //    loaded: "Luzifer",
    //    master_password: "foobar",
    //    selected: null,
    //},
    //accounts: ["Luzifer"],
    account_info: {},
    accounts: [],
    cryptocore_available: false,
    filter: null,
  },
  actions: {

    add_entry(context, entry) {
      entry.id = uuidv4()
      context.commit('add_entry', entry)
      context.dispatch('encrypt_data')
    },

    decrypt_data(context) {
      // Do not try to decrypt empty data
      if (context.state.account_info.data_raw == "") return

      new Promise((resolve, reject) => {
        opensslDecrypt(
          context.state.account_info.data_raw,
          context.state.account_info.master_password,
          (data, error) => {
            if (error) reject(error)
            resolve(data)
          },
        )
      }).then((plain_data) => {
        context.commit('decrypted_data', plain_data)
      }).catch((error) => console.log(error))
    },

    encrypt_data(context) {
      new Promise((resolve, reject) => {
        opensslEncrypt(
          JSON.stringify(context.state.account_info.data),
          context.state.account_info.master_password,
          (data, error) => {
            if (error) reject(error)
            resolve(data)
          },
        )
      }).then((enc_data) => {
        context.commit('encrypted_data', enc_data)
        context.dispatch('save_user_data')
      }).catch((error) => console.log(error))
    },

    enter_master_password(context, password) {
      context.commit('enter_master_password', password)
      context.dispatch('decrypt_data')
    },

    load_user_data(context, username) {
      axios.get(`/user/${username}/data`)
        .then((response) => {
          context.commit('account_loaded', {
            checksum: response.data.checksum,
            data: [],
            data_raw: response.data.data,
            loaded: username,
            master_password: null,
            selected: null,
          })
        })
        .catch((error) => console.log(error))
    },

    register(context, data) {
      axios.post("/register", data)
        .then((response) => {
          context.dispatch('reload_users', data.username)
        })
        .catch((error) => console.log(error))
    },

    reload_users(context, autoload = null) {
      axios.get("/users")
        .then((response) => {
          let users = []
          for (let user in response.data) {
            let login_state = response.data[user]
            if (login_state == 'logged-in') {
              users.push(user)
            }
          }

          context.commit('active_users', users)
          if (autoload != null) {
            context.dispatch('load_user_data', autoload)
          } else {
            context.commit('clear_active_user')
          }
        })
        .catch((error) => console.log(error))
    },

    save_user_data(context) {
      new Promise((resolve, reject) => {
        sha256sum(context.state.account_info.data_raw, (data, error) => {
          if (error) return reject(error)
          resolve(data)
        })
      }).then((checksum) => {
        axios.put(`/user/${context.state.account_info.loaded}/data`, {
          'checksum': checksum,
          'old_checksum': context.state.account_info.checksum,
          'data': context.state.account_info.data_raw,
        }).then((response) => {
          context.commit('update_checksum', checksum)
        }).catch((error) => console.log(error))
      }).catch((error) => console.log(error))

    },

    sign_in(context, auth) {
      console.log(auth)
      axios.post("/login", auth)
        .then((response) => {
          context.dispatch('reload_users', auth.username)
        })
        .catch((error) => {
          if (error.response) {
            console.log(error.response)
          } else {
            console.log(error)
          }
        })
    },

    sign_out(context) {
      axios.post(`/user/${context.state.account_info.loaded}/logout`)
        .then((response) => {
          context.dispatch('reload_users')
        })
        .catch((error) => console.log(error))
    },
  },
  getters: {
    filtered_entries: (state) => {
      if (state.filter === "" || state.filter === null) {
        return state.account_info.data
      }

      let entries = []

      for (let item of state.account_info.data) {
        if (item.title.indexOf(state.filter) > -1) {
          entries.push(item)
          continue
        }
        if (item.username.indexOf(state.filter) > -1) {
          entries.push(item)
          continue
        }
        if (item.tags.indexOf(state.filter) > -1) {
          entries.push(item)
          continue
        }
      }

      return entries
    },
    selected_entry: (state) => {
      if (!state.account_info.selected) {
        return null
      }
      for (let item of state.account_info.data) {
        if (item.id == state.account_info.selected) {
          return item
        }
      }
      return null
    }
  },
  mutations: {
    account_loaded(state, account_info) {
      state.account_info = account_info
    },

    active_users(state, users) {
      state.accounts = users
    },

    add_entry(state, entry) {
      state.account_info.data.push(entry)
    },

    clear_active_user(state) {
      state.account_info = {}
    },

    cryptocore_loaded(state) {
      state.cryptocore_available = true
    },

    decrypted_data(state, plain_data) {
      let elms = JSON.parse(plain_data)
      for (let e of elms) {
        // Migrate old entries without UUID
        if (!e.id) e.id = uuidv4()
        state.account_info.data.push(e)
      }
    },

    encrypted_data(state, enc_data) {
      state.account_info.data_raw = enc_data
    },

    enter_master_password(state, password) {
      state.account_info.master_password = password
    },

    lock_database(state) {
      state.account_info.data = []
      state.account_info.master_password = null
      state.account_info.selected = null
    },

    select_entry(state, idx) {
      state.account_info.selected = idx
    },

    update_checksum(state, checksum) {
      state.account_info.checksum = checksum
    },

    update_filter(state, value) {
      state.filter = value
    },
  }
})
