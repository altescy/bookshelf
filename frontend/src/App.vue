<template>
  <v-app>
    <Header/>
    <v-main>
      <AlertMessage v-if="!dialog"/>
      <router-view/>
    </v-main>
    <v-btn
      bottom
      color="pink"
      dark
      fab
      fixed
      right
      @click="openBookRegistrationDialog()"
    >
      <v-icon>mdi-plus</v-icon>
    </v-btn>
    <v-dialog
      v-model="dialog"
      persistent
      width="800px"
    >
      <BookRegistrationDialog v-if="dialogType === 'register'"/>
      <BookEditDialog v-if="dialogType === 'edit'"/>
    </v-dialog>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import {mapMutations, mapState} from 'vuex';
import Header from './components/Header.vue';
import AlertMessage from './components/AlertMessage.vue';
import BookRegistrationDialog from './components/BookRegistrationDialog.vue';
import BookEditDialog from './components/BookEditDialog.vue';
import * as VuexMutation from '@/vuex/mutation_types';


export default Vue.extend({
  name: 'App',

  components: {
    Header,
    AlertMessage,
    BookRegistrationDialog,
    BookEditDialog,
  },

  computed: {
    ...mapState(['dialog', 'dialogType', 'globalAlertMessage']),
  },

  methods: {
    ...mapMutations({
      openDialog: VuexMutation.OPEN_DIALOG,
      setDialogType: VuexMutation.SET_DIALOG_TYPE,
      unsetEditingBook: VuexMutation.UNSET_EDITING_BOOK,
    }),
    openBookRegistrationDialog () {
      this.unsetEditingBook();
      this.setDialogType('register');
      this.openDialog();
    },
  },
});
</script>
