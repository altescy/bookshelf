<template>
  <div style="margin: 0 3px">
    <v-btn
     rounded
     small
     :color="getColor(file.MimeType)"
     @click="onClick()"
    >
      <v-icon
       v-if="mode === 'delete'"
       small
      >mdi-close</v-icon>
      {{ getAlias(file.MimeType) }}
    </v-btn>

    <v-dialog
     v-model="deleteDialog"
     max-width="500px"
    >
      <v-card>
        <v-card-title class="headline">
          Delete this file?
        </v-card-title>
        <v-card-text>
          Are you sure you want to delete {{ getAlias(file.MimeType) }} file?
        </v-card-text>
        <v-card-actions>
          <v-btn
            text
            color="primary"
            @click="deleteFile()"
          >OK</v-btn>
          <v-btn
           text
           @click="deleteDialog = false"
          >Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
  import Vue from 'vue';
  import {mapActions, mapState} from 'vuex';
  import * as VuexAction from '@/vuex/action_types';
  import {getBaseMime, MimeToAlias, MimeToColor} from '@/file'

  export default Vue.extend({
    name: 'FileIcon',

    props: ['book', 'file', 'mode'],

    data: function() {
      return {
        deleteDialog: false,
      }
    },

    computed: {
      ...mapState(['mimes']),
    },

    methods: {
      ...mapActions({
        deleteFileAction: VuexAction.DELETE_FILE,
      }),
      getAlias(mime: string): string {
        mime = getBaseMime(mime);
        const alias = MimeToAlias.get(mime);
        if(alias === undefined) throw 'unknown mime';
        return String(alias);
      },
      getColor(mime: string): string {
        mime = getBaseMime(mime);
        const color = MimeToColor.get(mime);
        if(color === undefined) return "grey";
        return color
      },
      onClick: function() {
        if (this.mode === 'download') {
          const alias = this.getAlias(this.file.MimeType);
          const url = '/api/book/' + this.file.BookID + '/file/' + alias;

          const fileLink = document.createElement('a');
          fileLink.href = url;
          fileLink.download = this.book.Title + '.' + alias;
          fileLink.click();
          fileLink.remove();
        } else if (this.mode === 'delete') {
          this.deleteDialog = true;
        }
      },
      async deleteFile() {
        this.deleteDialog = false;
        await this.deleteFileAction(this.file).catch(error => {
          console.error("failed to delete file:", error);
        });
      },
    },
  });
</script>
