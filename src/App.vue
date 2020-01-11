<template>
  <div>
    <Navigation />

    <b-container>

      <Landing v-if="this.$store.state.accounts.length == 0 || !this.$store.state.account_info.loaded" />
      <PasswordView v-if="this.$store.state.account_info.loaded" />

    </b-container>

    <AddPasswordModal />
    <RegisterModal />
    <SignInModal />
  </div>
</template>

<script>
import AddPasswordModal from './components/AddPasswordModal.vue'
import Landing from './components/Landing.vue'
import Navigation from './components/Navigation.vue'
import PasswordView from './components/PasswordView.vue'
import RegisterModal from './components/RegisterModal.vue'
import SignInModal from './components/SignInModal.vue'

export default {
  components: {
    AddPasswordModal,
    Landing,
    Navigation,
    PasswordView,
    RegisterModal,
    SignInModal,
  },

  name: 'App',
}
</script>

<style lang="scss">
</style>
