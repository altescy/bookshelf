import Vue from 'vue';
import Vuex from 'vuex';
import axios, {AxiosResponse} from 'axios';
import * as Model from '@/model';
import * as VuexMutation from '@/vuex/mutation_types';
import * as VuexAction from '@/vuex/action_types';

Vue.use(Vuex)

const API_ENDPOINT = 'http://localhost/api';
const OPENBD_ENDPOINT = 'https://api.openbd.jp/v1';

const initialState: Model.State = {
  alertMessage: {
    type: 'success',
    message: '',
  },
  books: [],
  dialog: false,
  search: '',
  editingBook: {
    ID: 0,
    CreatedAt: '',
    UpdatedAt: '',
    ISBN: '',
    Title: '',
    Author: '',
    Description: '',
    CoverURL: '',
    Publisher: '',
    PubDate: '',
    Files: [],
  }
}

function extractBookFromOpenBDResponse(response: AxiosResponse): Model.Book {
  const data = response.data[0];
  if (!data) {
    throw 'invalid ISBN';
  }
  const convertPubdate = (pubdate: string): string => {
    const year = pubdate.slice(0, 4);
    const month = pubdate.slice(4, 6);
    const date = pubdate.slice(6, 8);
    return year + '-' + month + '-' + date;
  };
  const book: Model.Book = {
    ID: 0,
    CreatedAt: '',
    UpdatedAt: '',
    ISBN: data.summary.isbn,
    Title: data.summary.title,
    Author: data.summary.author,
    Publisher: data.summary.publisher,
    PubDate: convertPubdate(data.summary.pubdate),
    CoverURL: data.summary.cover,
    Description: data.onix.CollateralDetail.TextContent.slice(-1)[0].Text,
    Files: [],
  };
  return book;
}

export default new Vuex.Store({
  state: initialState,
  mutations: {
    [VuexMutation.SET_ALERT_MESSAGE](state: Model.State, alertMessage: Model.AlertMessage) {
      state.alertMessage= alertMessage;
    },
    [VuexMutation.OPEN_DIALOG](state: Model.State) {
      state.dialog = true;
    },
    [VuexMutation.CLOSE_DIALOG](state: Model.State) {
      state.dialog = false;
    },
    [VuexMutation.SET_EDITING_BOOK](state: Model.State, book) {
      state.editingBook = book;
    },
    [VuexMutation.UNSET_EDITING_BOOK](state: Model.State) {
      state.editingBook = {
        ID: 0,
        CreatedAt: '',
        UpdatedAt: '',
        ISBN: '',
        Title: '',
        Author: '',
        Description: '',
        CoverURL: '',
        Publisher: '',
        PubDate: '',
        Files: [],
      };
    },
    [VuexMutation.SET_BOOKS](state: Model.State, books) {
      state.books = books;
    },
    [VuexMutation.ADD_BOOK](state: Model.State, book) {
      state.books = state.books.concat(book);
    },
    [VuexMutation.SET_SEARCH_QUERY](state: Model.State, query: string) {
      state.search = query;
    },
  },
  actions: {
    async [VuexAction.AUTOCOMPLETE_EDITING_BOOK_BY_ISBN]({ commit }) {
      const isbn = this.state.editingBook.ISBN;
      try {
        const response = await axios.get(OPENBD_ENDPOINT + '/get?isbn=' + isbn)
        if (response.status === 200) {
          const book = extractBookFromOpenBDResponse(response);
          commit(VuexMutation.SET_EDITING_BOOK, book);
        } else {
          throw 'failed to fetch book information'
        }
      } catch (error) {
        const alertMessage: Model.AlertMessage = {
          type: 'error',
          message: String(error),
        }
        commit(VuexMutation.SET_ALERT_MESSAGE, alertMessage);
        console.error(error)
      }
    },
    async [VuexAction.FETCH_ALL_BOOKS]({ commit }) {
      try {
        const response = await axios.get(API_ENDPOINT + '/books');
        if (response.status === 200) {
          commit(VuexMutation.SET_BOOKS, response.data);
        } else {
          throw response.data.err || 'unexpected error';
        }
      } catch (error) {
        const alertMessage: Model.AlertMessage = {
          type: 'error',
          message: String(error),
        }
        commit(VuexMutation.SET_ALERT_MESSAGE, alertMessage);
        throw error
      }
    },
    async [VuexAction.REGISTER_EDITING_BOOK]({ commit }) {
      try {
        if (!this.state.editingBook.Title) {
          throw 'Title is empty.';
        }

        const params = new URLSearchParams();
        params.append('ISBN', this.state.editingBook.ISBN);
        params.append('Title', this.state.editingBook.Title);
        params.append('Author', this.state.editingBook.Author);
        params.append('Publisher', this.state.editingBook.Publisher);
        params.append('PubDate', this.state.editingBook.PubDate);
        params.append('CoverURL', this.state.editingBook.CoverURL);
        params.append('Description', this.state.editingBook.Description);

        const response = await axios.post(API_ENDPOINT + '/book', params);
        if (response.status === 200) {
          commit(VuexMutation.ADD_BOOK, response.data);
        } else {
          throw response.data.err || 'unexpected error';
        }
      } catch (error) {
        const alertMessage: Model.AlertMessage = {
          type: 'error',
          message: String(error),
        }
        console.error(error)
        commit(VuexMutation.SET_ALERT_MESSAGE, alertMessage);
        throw error
      }
    },
  },
  modules: {
  }
})
