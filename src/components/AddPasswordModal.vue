<template>

  <b-modal
    id="addPassword"
    lazy
    ok-title="Add"
    size="lg"
    title="Add Password"
    @ok="onOK">
    <b-form
      @submit.prevent
      id="addPasswordForm">

      <b-form-group label="Title">
        <b-form-input
          type="text"
          name="title"
          required
          placeholder="Title"/>
      </b-form-group>

      <b-form-group label="Username">
        <b-form-input
          type="text"
          name="username"
          placeholder="Username"/>
      </b-form-group>

      <b-form-group label="Password">
        <b-input-group>
          <b-form-input
            type="text"
            name="password"
            placeholder="Password"/>
          <b-input-group-append>
            <b-btn
              variant="primary"
              title="Generate password"><fa-icon icon="sync" /></b-btn>
          </b-input-group-append>
        </b-input-group>
      </b-form-group>

      <b-form-group label="URL">
        <b-form-input
          type="text"
          name="url"
          placeholder="URL"/>
      </b-form-group>

      <b-form-group label="Comments">
        <b-form-textarea
          type="text"
          name="comments"
          rows="4"
          placeholder="Comments / notes for the entry"/>
      </b-form-group>

      <b-form-group label="Tags">
        <b-form-input
          type="text"
          name="tags"
          placeholder="Tags"/>
      </b-form-group>

    </b-form>
  </b-modal>

</template>

<script>
export default {
  methods: {
    onOK() {
      const form = document.getElementById('addPasswordForm')
      const entry = {}

      for (var pair of new FormData(form).entries()) {
        entry[pair[0]] = pair[1]
      }

      this.$store.dispatch('add_entry', entry)
      form.reset()
    },
  },

  name: 'Addpassword',
}
</script>
