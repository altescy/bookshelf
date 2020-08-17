<template>
  <div>
    <v-card
      v-for="book in filteredBooks()"
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
        <FileIcon
         v-for="file in book.Files"
         v-bind:key="file.ID"
         :book="book"
         :file="file"
         mode="download"
        />
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
  import Vue from 'vue';
  import {mapActions, mapMutations, mapState} from 'vuex';
  import * as Model from '@/model';
  import * as VuexAction from '@/vuex/action_types';
  import * as VuexMutation from '@/vuex/mutation_types';
  import FileIcon from '@/components/FileIcon.vue';

  export default Vue.extend({
    name: 'BookList',

    components: {
      FileIcon,
    },

    computed: {
      ...mapState(['books', 'search']),
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
      openBookEditDialog(book: Model.Book) {
        this.setEditingBook(book);
        this.setDialogType('edit');
        this.openDialog();
      },
      filteredBooks() {
        const query = this.search.toLowerCase();
        return this.books.filter(function(e: Model.Book): boolean  {
          const text = [e.Title, e.Author, e.Publisher].join('  ').toLowerCase();
          return text.indexOf(query) > -1;
        });
      },
    },

    mounted() {
      this.fetchAllBooks();
    },
  })
</script>
