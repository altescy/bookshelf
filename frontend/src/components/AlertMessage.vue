<template>
  <div>
    <v-alert
      dense
      v-for="msg in alertMessages"
      v-bind:key="msg.id"
      :type="msg.type"
      @click="clearAlert(msg)"
      style="margin: 0;"
    >
      {{ msg.message }}
    </v-alert>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import {mapMutations, mapState} from 'vuex';
import * as Model from '@/model';
import * as VuexMutataion from '@/vuex/mutation_types';

export default Vue.extend({
  name: 'AlertMessage',

  computed: {
    ...mapState(['alertMessages']),
  },

  methods: {
    ...mapMutations({
      deleteAlertMessage: VuexMutataion.DELETE_ALERT_MESSAGE,
    }),
    clearAlert (msg: Model.AlertMessage) {
      this.deleteAlertMessage(msg);
    },

  },

});
</script>
