<template>
  <v-container>
    <v-row class="mx-2">
      <v-col cols="8">
        <v-row>
          <v-col cols="10">
            <v-text-field
              v-model="editingBook.ISBN"
              prepend-icon="mdi-numeric"
              placeholder="ISBN"
            ></v-text-field>
          </v-col>
          <v-col cols="2">
            <v-btn
              @click="autocomplete()"
            >auto complete</v-btn>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="12">
        <v-text-field
          v-model="editingBook.Title"
          prepend-icon="mdi-book-open-page-variant"
          placeholder="Title"
        ></v-text-field>
      </v-col>
      <v-col cols="12">
        <v-text-field
          v-model="editingBook.Author"
          prepend-icon="mdi-account"
          placeholder="Author"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="editingBook.Publisher"
          prepend-icon="mdi-domain"
          placeholder="Publisher"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="editingBook.PubDate"
          type="date"
          prepend-icon="mdi-calendar"
          placeholder="Publication Date"
        ></v-text-field>
      </v-col>
      <v-col cols="12">
        <v-text-field
          v-model="editingBook.CoverURL"
          prepend-icon="mdi-camera"
          placeholder="Cover URL"
        ></v-text-field>
      </v-col>
      <v-col cols="12">
        <v-textarea
          v-model="editingBook.Description"
          prepend-icon="mdi-text"
          placeholder="Description"
          rows="1"
        ></v-textarea>
      </v-col>
      <v-col cols="12" v-if="editingBook.Files.length > 0">
        <v-icon
         style="margin-right: 8px"
        >mdi-file-document-multiple-outline</v-icon>
        <FileIcon
         v-for="file in editingBook.Files"
         v-bind:key="file.ID"
         :book="editingBook"
         :file="file"
         mode="delete"
         style="display: inline"
        />
      </v-col>
      <v-col cols="12">
        <v-file-input
         v-model="files"
         small-chips
         multiple
         label="Files"
         :accept="acceptMimes()"
        ></v-file-input>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
  import Vue from 'vue'
  import {mapActions, mapMutations, mapState} from 'vuex';
  import * as VuexAction from '@/vuex/action_types';
  import * as VuexMutation from '@/vuex/mutation_types';
  import FileIcon from '@/components/FileIcon.vue';

  export default Vue.extend({
    name: 'BookRegistrationDialog',

    components: {
      FileIcon,
    },

    computed: {
      ...mapState(['editingBook', 'mimes']),
      files: {
        get(): File[] {
          return this.$store.state.files;
        },
        set(files: File[]) {
          this.setFiles(files);
        },
      },
    },

    methods: {
      ...mapActions({
        autocomplete: VuexAction.AUTOCOMPLETE_EDITING_BOOK_BY_ISBN,
      }),
      ...mapMutations({
        setFiles: VuexMutation.SET_FILES,
      }),
      acceptMimes(): string {
        const ms: string[] = [];
        for(const ext in this.mimes) ms.push(ext);
        return ms.join(',');
      },
    },
  })
</script>
