<template>

  <b-modal
    id="register"
    lazy
    ok-title="Register"
    size="lg"
    title="Register a new account"
    @ok="onOK">
    <b-form
      @submit.prevent
      id="registerForm">

      <b-form-group label="Username">
        <b-form-input
          type="text"
          name="username"
          placeholder="Username"/>
      </b-form-group>

      <b-form-group label="Password">
        <b-form-input
          type="password"
          name="password"
          placeholder="Password"/>
      </b-form-group>

      <b-form-group label="Repeat Password">
        <b-form-input
          type="password"
          name="check_password"
          placeholder="Repeat Password"/>
      </b-form-group>

    </b-form>
  </b-modal>

</template>

<script>
export default {
  methods: {
    onOK() {
      const form = document.getElementById('registerForm')
      const entry = {}

      for (var pair of new FormData(form).entries()) {
        entry[pair[0]] = pair[1]
      }

      this.$store.dispatch('register', entry)
      form.reset()
    },
  },

  name: 'Register',
}
</script>
