<template>
  <v-card>
    <v-card-title class="grey darken-2 white--text">
      Edit Book
    </v-card-title>

    <AlertMessage/>

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
      >Close</v-btn>
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
        setOverlay: VuexMutation.SET_OVERLAY,
        closeDialog: VuexMutation.CLOSE_DIALOG,
        setFiles: VuexMutation.SET_FILES,
        setEditingBook: VuexMutation.SET_EDITING_BOOK,
        unsetEditingBook: VuexMutation.UNSET_EDITING_BOOK,
      }),
      cancel () {
        this.closeDialog();
        this.setFiles([]);
        this.unsetEditingBook();
      },
      async update() {
        this.setOverlay(true);

        await this.updateBook().then(() => {
          this.setFiles([]);
        }).catch(() => {
          console.error('failed to update book');
        });

        this.setOverlay(false);
      },
      async deleteBook() {
        this.deleteDialog = false;
        this.setOverlay(true);

        await this.deleteEditingBook().then(() => {
          this.closeDialog();
          this.setFiles([]);
          this.unsetEditingBook();
        }).catch(error => {
          console.error('failed to delete book:', error)
        });

        this.setOverlay(false);
      },
    },
  })
</script>
