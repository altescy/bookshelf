<template>
  <v-card>
    <v-card-title class="grey darken-2 white--text">
      Register Book
    </v-card-title>

    <AlertMessage/>

    <BookEditor/>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        text
        color="primary"
        @click.prevent="register()"
      >Register</v-btn>
      <v-btn
        text
        @click="cancel()"
      >Cancel</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
  import Vue from 'vue'
  import {mapActions, mapMutations, mapState} from 'vuex';
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
      ...mapState(['editingBook']),
    },

    methods: {
      ...mapActions({
        autocomplete: VuexAction.AUTOCOMPLETE_EDITING_BOOK_BY_ISBN,
        registerBook: VuexAction.REGISTER_EDITING_BOOK,
      }),
      ...mapMutations({
        setOverlay: VuexMutation.SET_OVERLAY,
        closeDialog: VuexMutation.CLOSE_DIALOG,
        unsetEditingBook: VuexMutation.UNSET_EDITING_BOOK,
        setFiles: VuexMutation.SET_FILES,
      }),

      cancel () {
        this.closeDialog();
        this.setFiles([]);
        this.unsetEditingBook();
      },
      async register() {
        this.setOverlay(true);

        await this.registerBook().then(() => {
          this.closeDialog();
          this.setFiles([]);
          this.unsetEditingBook();
        }).catch(error => {
          console.error("failed to register book:", error)
        });

        this.setOverlay(false);
      },
    },
  })
</script>
