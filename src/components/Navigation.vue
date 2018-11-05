<template>
  <b-navbar toggleable="md" type="dark" variant="primary" class="mb-4">
    <b-navbar-brand class="mr-3">
      <fa-icon icon="key" /> cloud<strong>keys</strong>
    </b-navbar-brand>
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav class="mr-auto">
        <b-nav-item href="#"
                    class="mr-2"
                    title="Add new password (Ctrl+N)"
                    v-shortkey.once="['ctrl', 'n']"
                    @shortkey=""
                    v-if="this.$store.state.account_info.master_password"
                    v-b-modal.addPassword>
          <fa-icon icon="plus-circle" /> Add Password

        </b-nav-item>

        <b-nav-item href="#"
                    title="Import passwords from KeePass"
                    v-if="this.$store.state.account_info.master_password"
                    class="d-none d-sm-none d-md-none d-lg-block d-xl-block">
          <fa-icon icon="file-import" /> Import from KeePass
        </b-nav-item>
      </b-navbar-nav>

      <b-navbar-nav class="ml-auto">

        <b-nav-item-dropdown text="Select Account" right v-if="this.$store.state.accounts.length > 0">
          <b-dropdown-item href="#" 
                           v-for="user in this.$store.state.accounts"
                           :key="user"
                           @click="$store.dispatch('load_user_data', user)">
            <fa-icon icon="user" />
            <strong v-if="$store.state.account_info.loaded == user">{{ user }}</strong>
            <span v-else>{{ user }}</span>
          </b-dropdown-item>

          <b-dropdown-divider />
          <b-dropdown-item href="#"
                           v-b-modal.signIn>
            <fa-icon icon="sign-in-alt" /> Sign in
          </b-dropdown-item>
        </b-nav-item-dropdown>

        <b-nav-item href="#"
                    class="mr-3"
                    title="Sign in to an existing account"
                    v-if="this.$store.state.accounts.length == 0"
                    v-b-modal.signIn>
          <fa-icon icon="sign-in-alt" /> Login
        </b-nav-item>

        <b-nav-item href="#"
                    title="Lock current account (Ctrl+L)"
                    v-shortkey.once="['ctrl', 'l']"
                    @shortkey="$store.commit('lock_database')"
                    @click="$store.commit('lock_database')"
                    v-if="this.$store.state.account_info.master_password"
                    class="mr-2 d-none d-sm-none d-md-none d-lg-block d-xl-block">
          <fa-icon icon="lock" />
        </b-nav-item>

        <b-nav-item href="#"
                    class="mr-2"
                    title="Sign out current account"
                    v-if="this.$store.state.account_info.loaded"
                    @click="$store.dispatch('sign_out')">
          <fa-icon icon="sign-out-alt" />
          <span class="d-inline d-sm-inline d-md-inline d-lg-none d-xl-none">Sign out</span>
        </b-nav-item>

        <b-nav-item href="#"
                    title="Register new account"
                    class="mr-2 d-none d-sm-none d-md-none d-lg-block d-xl-block"
                    v-b-modal.register>
          <fa-icon icon="user-plus" />
        </b-nav-item>

        <b-nav-item href="#"
                    title="Info"
                    class="ml-2 d-none d-sm-none d-md-none d-lg-block d-xl-block">
          <fa-icon icon="info-circle" />
        </b-nav-item>

      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
export default {
  name: 'navigation',
}
</script>
