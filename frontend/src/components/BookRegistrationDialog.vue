<template>
  <v-card>
    <v-card-title class="grey darken-2 white--text">
      Register Book
    </v-card-title>

    <AlertMessage />

    <BookEditor/>

    <div style="margin: 0 30px">
      <v-file-input
       small-chips
       multiple
       label="Files"
       accept="azw,.azw3,.cbr,.cbz,.cbt,.cb7,.epub,.epub3,.fb2,.fb2.zip,.mobi,.pdf,.txt"
      ></v-file-input>
    </div>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        text
        color="primary"
        @click="cancel()"
      >Cancel</v-btn>
      <v-btn
        text
        @click.prevent="register()"
      >Register</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
  import Vue from 'vue'
  import {mapActions, mapMutations, mapState} from 'vuex';
  import * as Model from '@/model';
  import * as VuexAction from '@/vuex/action_types';
  import * as VuexMutation from '@/vuex/mutation_types';
  import AlertMessage from '@/components/AlertMessage.vue';
  import BookEditor from '@/components/BookEditor.vue';

  export default Vue.extend({
    name: 'BookRegistrationDialog',

    components: {
      AlertMessage,
      BookEditor,
    },

    computed: {
      ...mapState(['editingBook'])
    },

    methods: {
      ...mapActions({
        autocomplete: VuexAction.AUTOCOMPLETE_EDITING_BOOK_BY_ISBN,
        registerBook: VuexAction.REGISTER_EDITING_BOOK,
      }),
      ...mapMutations({
        closeDialog: VuexMutation.CLOSE_DIALOG,
        unsetEditingBook: VuexMutation.UNSET_EDITING_BOOK,
        setAlertMessage: VuexMutation.SET_ALERT_MESSAGE,
      }),
      clearAlert() {
        const msg: Model.AlertMessage = {
          type: 'success',
          message: '',
        };
        this.setAlertMessage(msg);
      },
      cancel () {
        this.closeDialog();
        this.clearAlert();
        this.unsetEditingBook();
      },
      async register() {
        this.registerBook().then(() => {
          this.closeDialog();
          this.clearAlert();
          this.unsetEditingBook();
        });
      },
    },
  })
</script>
