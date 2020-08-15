<template>
  <div class="home">
    <BookList/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import BookList from '@/components/BookList.vue';

@Component({
  components: {
    BookList,
  },
})
export default class Home extends Vue {}
</script>
