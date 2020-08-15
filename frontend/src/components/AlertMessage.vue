<template>
  <v-alert
    dense
    v-if="alertMessage.message"
    :type="alertMessage.type"
    @click="clearAlert()"
  >
    {{ alertMessage.message }}
  </v-alert>
</template>

<script lang="ts">
import Vue from 'vue';
import {mapMutations, mapState} from 'vuex';
import * as Model from '@/model';
import * as VuexMutataion from '@/vuex/mutation_types';

export default Vue.extend({
  name: 'AlertMessage',

  computed: {
    ...mapState(['alertMessage']),
  },

  methods: {
    ...mapMutations({
      setAlertMessage: VuexMutataion.SET_ALERT_MESSAGE,
    }),
    clearAlert () {
      const msg: Model.AlertMessage = {
        type: 'success',
        message: '',
      };
      this.setAlertMessage(msg);
    },

  },

});
</script>
