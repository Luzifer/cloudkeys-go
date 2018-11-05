<template>
  <div>
    <b-row align-h="center">

      <b-col cols="12" v-if="this.$store.state.account_info.data_raw.length == 0 && !this.$store.state.account_info.master_password">
        <b-alert show variant="primary">
          It looks like you've not stored any passwords yet&hellip;<br>Therefore you now need to choose your <strong>master password</strong> which is used to encrypt your passwords you will store in the future. Choose wise, do not re-use your login password and don't tell anyone but also <strong>do never loose this password</strong>. Nobody will be able to recover your passwords if you do.
        </b-alert>
      </b-col>

      <b-col cols="8" v-if="this.$store.state.account_info.loaded && !this.$store.state.account_info.master_password">
        <b-form @submit="onMasterPassword" id="masterPassword">
          <b-form-group :label="masterPasswordHint">
            <b-form-input type="password"
                          size="lg"
                          data-lpignore="true"
                          name="master"
                          required
                          :disabled="!this.$store.state.cryptocore_available"
                          placeholder="Master Password" />
          </b-form-group>
        </b-form>
      </b-col>

      <b-col cols="12" v-if="!this.$store.state.cryptocore_available">
        <b-alert show variant="warning">
          The cryptographic library is not yet ready. Until then it is not possible to decrypt the password store.
        </b-alert>
      </b-col>

      <b-col cols="12" v-if="this.$store.state.account_info.data.length == 0 && this.$store.state.account_info.master_password">
        <b-alert show variant="primary">
          It looks like you've not stored any passwords yet&hellip;<br>After you've now set your <strong>master password</strong> you can add a new entry using the buttons on the upper left: Create a new password by hand or import them from an existing KeePass file&hellip;
        </b-alert>
      </b-col>
    </b-row>

    <b-row v-if="this.$store.state.account_info.data.length > 0" class="w-100">
      <b-col cols="4">
        <b-form>
          <b-form-group>
            <b-form-input type="text"
                          placeholder="Filter / Search"
                          @change="updateFilter"
                          @input="updateFilter"
                          v-shortkey.focus="{'f': ['ctrl', 'f'], 'slash': ['/']}" />
          </b-form-group>
        </b-form>

        <b-list-group class="mt-3">
          <b-list-group-item v-for="entry in $store.getters.filtered_entries"
                             button
                             :key="entry.id"
                             @click="loadEntry(entry.id)">
            {{ entry.title }}<br>
            <small>{{ entry.username }}</small>
          </b-list-group-item>
        </b-list-group>
      </b-col>

      <b-col cols="8">
        <div v-if="this.selectedEntry">
          <h1>{{selectedEntry.title}}</h1>

          <b-form-group label="Username">
            <b-input-group>
              <b-form-input type="text" readonly :value="selectedEntry.username" />
              <b-input-group-append>
                <b-btn :variant="copyClass('username')"
                       v-shortkey="['ctrl', 'b']"
                       @shortkey="copy(selectedEntry.username, 'username')"
                       @click="copy(selectedEntry.username, 'username')"
                       title="Copy to Clipboard (Ctrl+B)">
                  <fa-icon icon="clipboard" />
                </b-btn>
              </b-input-group-append>
            </b-input-group>
          </b-form-group>

          <b-form-group label="Password">
            <b-input-group>
              <b-form-input :type="this.passwordVisible ? 'text' : 'password'" readonly :value="selectedEntry.password" />
              <b-input-group-append>
                <b-btn variant="warning"
                       @click="showPassword"
                       title="Show password in plain txt">
                  <fa-icon icon="eye" />
                </b-btn>
                <b-btn :variant="copyClass('password')"
                       v-shortkey="['ctrl', 'c']"
                       @shortkey="copy(selectedEntry.password, 'password')"
                       @click="copy(selectedEntry.password, 'password')"
                       title="Copy to Clipboard (Ctrl+C)">
                  <fa-icon icon="clipboard" />
                </b-btn>
              </b-input-group-append>
            </b-input-group>
          </b-form-group>

          <b-form-group label="URL">
            <b-input-group>
              <b-form-input type="text" readonly :value="selectedEntry.url" />
              <b-input-group-append>
                <b-btn :variant="copyClass('url')"
                       v-shortkey="['ctrl', 'u']"
                       @shortkey="copy(selectedEntry.url, 'url')"
                       @click="copy(selectedEntry.url, 'url')"
                       title="Copy to Clipboard (Ctrl+U)">
                  <fa-icon icon="clipboard" />
                </b-btn>
              </b-input-group-append>
            </b-input-group>
          </b-form-group>

          <b-form-group label="Comments">
            <b-textarea readonly>{{ selectedEntry.comment }}</b-textarea>
          </b-form-group>

          <b-form-group label="Tags">
            <b-form-input type="text" readonly :value="selectedEntry.tags" />
          </b-form-group>

        </div>

        <b-alert variant="info" show v-else>
          <fa-icon icon="arrow-left" />Select a database entry on the left to view it&hellip;
        </b-alert>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  name: 'passwordview',
  data: function() {
    return {
      copySuccess: null,
      copyError: null,
      passwordVisible: false,
      timerHandle: null,
    }
  },
  computed: {
    masterPasswordHint: function() {
      if (this.$store.state.account_info.data_raw) {
        return "Please enter your master password to decrypt your passwords:"
      } else {
        return "Please set your master password:"
      }
    },
    selectedEntry: function() {
      return this.$store.getters.selected_entry
    },
  },
  methods: {
    copy: function(payload, component) {
      this.$copyText(payload).then(()=>{
        this.copySuccess = component
        this.copyError = ""
      }, ()=>{
        this.copySuccess = ""
        this.copyError = component
      })

      // Remove old timer if exists
      if (this.timerHandle !== null) {
        window.cancelTimeout(this.timerHandle)
      }

      // Create a new timer to restore buttons to primary color
      let that = this
      this.timerHandle = window.setTimeout(function(){
        that.copySuccess = "" 
        that.copyError = ""
        that.timerHandle = null
      }, 3000)
    },
    copyClass: function(component) {
      if (this.copySuccess == component) return "success"
      if (this.copyError == component) return "danger"
      return "primary"
    },
    loadEntry: function(idx) {
      this.passwordVisible = false
      this.$store.commit('select_entry', idx)
    },
    onMasterPassword(e) {
e.preventDefault()

      let form = document.getElementById('masterPassword')
      let entry = {}

      for(var pair of new FormData(form).entries()) {
        entry[pair[0]] = pair[1]
      }

      this.$store.dispatch('enter_master_password', entry.master)
      form.reset()

      return false
    },
    showPassword: function() {
      this.passwordVisible = !this.passwordVisible
      let that = this
      window.setTimeout(function(){ that.passwordVisible = false }, 5000)
    },
    updateFilter: function(value) {
      this.$store.commit('update_filter', value)
    },
  },
}
</script>

<style scoped lang="scss">
button {
  cursor: pointer;
}
</style>
