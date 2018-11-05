<template>

  <b-modal id="signIn"
           lazy
           ok-title="Sign-In"
           size="lg"
           title="Sign in using existing account"
           @ok="onOK">
    <b-form @submit.prevent id="signInForm">

      <b-form-group label="Username">
        <b-form-input type="text"
                      name="username"
                      placeholder="Username"/>
      </b-form-group>

      <b-form-group label="Password">
        <b-form-input type="password"
                      name="password"
                      placeholder="Password"/>
      </b-form-group>

    </b-form>
  </b-modal>

</template>

<script>
export default {
  name: 'signin',
  methods: {
    onOK: function() {
      let form = document.getElementById('signInForm')
      let entry = {}

      for(var pair of new FormData(form).entries()) {
        entry[pair[0]] = pair[1]
      }

      this.$store.dispatch('sign_in', entry)
      form.reset()
    },
  }
}
</script>
