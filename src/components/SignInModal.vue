<template>

  <b-modal
    id="signIn"
    lazy
    ok-title="Sign-In"
    size="lg"
    title="Sign in using existing account"
    @ok="issueLogin">
    <b-form
      @submit.prevent="issueLogin"
      id="signInForm">

      <b-form-group label="Username">
        <b-form-input
          type="text"
          name="username"
          placeholder="Username"
          v-model="login.username"/>
      </b-form-group>

      <b-form-group label="Password">
        <b-form-input
          type="password"
          name="password"
          placeholder="Password"
          v-model="login.password"/>
      </b-form-group>

    </b-form>
  </b-modal>

</template>

<script>
export default {
  data() {
    return {
      login: {
        password: '',
        username: '',
      },
    }
  },

  methods: {
    issueLogin() {
      this.$store.dispatch('sign_in', this.login)
      this.login = { password: '', username: '' }
    },
  },

  name: 'Signin',
}
</script>
