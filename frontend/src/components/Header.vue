<template>
  <v-app-bar
     :clipped-left="$vuetify.breakpoint.lgAndUp"
     app
     color="blue darken-3"
     dark
  >
    <v-avatar
      size="32px"
      item
    >
      <v-icon large>mdi-bookshelf</v-icon>
    </v-avatar>

    <v-toolbar-title
      style="width: 300px"
      class="ml-0 pl-4"
    >
      <span class="hidden-sm-and-down">Bookshelf</span>
    </v-toolbar-title>

    <v-text-field
      v-model="query"
      flat
      solo-inverted
      hide-details
      prepend-inner-icon="mdi-magnify"
      label="Search"
      class="hidden-sm-and-down"
    ></v-text-field>

    <v-spacer></v-spacer>

    <a href="/opds">
      <v-btn icon>
        <v-icon>mdi-rss</v-icon>
      </v-btn>
    </a>

  </v-app-bar>
</template>

<script lang="ts">
  import Vue from 'vue';
  import {mapState, mapMutations} from 'vuex';
  import * as VuexMutation from '@/vuex/mutation_types';

  export default Vue.extend({
    name: 'Header',

    computed: {
      ...mapState(['search']),
      query: {
        get(): string {
          return this.search;
        },
        set(value: string) {
          this.setSearchQuery(value)
        },
      }
    },

    methods: {
      ...mapMutations({
        setSearchQuery: VuexMutation.SET_SEARCH_QUERY,
      }),
    },

  });
</script>
