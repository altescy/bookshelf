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
      @click="openDialog()"
    >
      <v-icon>mdi-plus</v-icon>
    </v-btn>
    <v-dialog
      v-model="dialog"
      persistent
      width="800px"
    >
      <BookRegistrationDialog/>
    </v-dialog>

  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import {mapMutations, mapState} from 'vuex';
import Header from './components/Header.vue';
import AlertMessage from './components/AlertMessage.vue';
import BookRegistrationDialog from './components/BookRegistrationDialog.vue';
import * as VuexMutation from '@/vuex/mutation_types';


export default Vue.extend({
  name: 'App',

  components: {
    Header,
    AlertMessage,
    BookRegistrationDialog,
  },

  computed: {
    ...mapState(['dialog', 'globalAlertMessage']),
  },

  methods: {
    ...mapMutations({
      openDialog: VuexMutation.OPEN_DIALOG,
    }),

  },
});
</script>
