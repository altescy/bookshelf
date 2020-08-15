<template>
  <div>
    <v-card
      v-for="book in filteredBooks"
      v-bind:key="book.id"
      class="mx-auto"
      max-width="1000px"
      outlined
      style="margin: 10px 0"
    >
      <v-list-item three-line>
        <v-list-item-content>
          <div class="overline mb-4">
            <span>{{ book.Author }}</span>
            <v-btn
             text style="float: right"
             @click="openBookEditDialog(book)"
            >
              <v-icon small>
                mdi-pencil
              </v-icon>
            </v-btn>
          </div>

          <v-list-item-title class="headline mb-1">{{ book.Title }}</v-list-item-title>
          <v-list-item-subtitle>{{ book.Description }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-card-actions>
        <v-btn text>Button</v-btn>
        <v-btn text>Button</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
  import Vue from 'vue'
  import {mapActions, mapMutations, mapState} from 'vuex';
  import * as VuexAction from '@/vuex/action_types';
  import * as VuexMutation from '@/vuex/mutation_types';

  export default Vue.extend({
    name: 'BookList',

    computed: {
      ...mapState(['books', 'search']),
      filteredBooks: function () {
        return this.books.filter(e => {
          const text = [e.Title, e.Author, e.Publisher].join('  ').toLowerCase();
          return text.indexOf(this.search.toLowerCase()) > -1;
        });
      },
    },

    methods: {
      ...mapActions({
        fetchAllBooks: VuexAction.FETCH_ALL_BOOKS,
      }),
      ...mapMutations({
        openDialog: VuexMutation.OPEN_DIALOG,
        setDialogType: VuexMutation.SET_DIALOG_TYPE,
        setEditingBook: VuexMutation.SET_EDITING_BOOK,
      }),
      openBookEditDialog(book) {
        this.setEditingBook(book)
        this.setDialogType('edit');
        this.openDialog();
      },
    },

    mounted() {
      this.fetchAllBooks();
    }
  })
</script>
