<template>
  <v-card>
    <v-card-title class="grey darken-2 white--text">
      Edit Book
    </v-card-title>

    <AlertMessage />

    <BookEditor/>

    <v-card-actions>
      <v-btn
        text
        color="red"
        @click="deleteDialog = true"
      >Delete</v-btn>
      <v-spacer></v-spacer>
      <v-btn
        text
        color="primary"
        @click.prevent="update()"
      >Update</v-btn>
      <v-btn
        text
        @click="cancel()"
      >Cancel</v-btn>
    </v-card-actions>
    <v-dialog
     v-model="deleteDialog"
     max-width="500px"
    >
      <v-card>
        <v-card-title class="headline">
          Delete this book?
        </v-card-title>
        <v-card-text>
          Are you sure you want to delete this book?
        </v-card-text>
        <v-card-actions>
          <v-btn
            text
            color="primary"
            @click="deleteBook()"
          >OK</v-btn>
          <v-btn
           text
           @click="deleteDialog = false"
          >Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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

    props: ['bookID'],

    data: function() {
      return {
        deleteDialog: false,
      }
    },

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
        updateBook: VuexAction.UPDATE_BOOK,
        deleteEditingBook: VuexAction.DELETE_EDITING_BOOK,
      }),
      ...mapMutations({
        closeDialog: VuexMutation.CLOSE_DIALOG,
        setEditingBook: VuexMutation.SET_EDITING_BOOK,
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
      async update() {
        this.updateBook().then(() => {
          this.closeDialog();
          this.clearAlert();
          this.unsetEditingBook();
        });
      },
      async deleteBook() {
        this.deleteEditingBook().then(() => {
          this.deleteDialog = false;
          this.closeDialog();
          this.unsetEditingBook();
        }, () => {
          console.error('failed to delete book')
          this.deleteDialog = false;
        })
      },
    },
  })
</script>
