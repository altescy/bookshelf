import Vue from 'vue';
import Vuex from 'vuex';
import axios, {AxiosError, AxiosResponse} from 'axios';
import * as Model from '@/model';
import * as VuexMutation from '@/vuex/mutation_types';
import * as VuexAction from '@/vuex/action_types';
import {deepCopy} from '@/utils';

Vue.use(Vuex)

const API_ENDPOINT = 'http://localhost/api';
const OPENBD_ENDPOINT = 'https://api.openbd.jp/v1';

const emptyBook: Model.Book = {
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

const initialState: Model.State = {
  alertMessages: [],
  books: [],
  dialog: false,
  dialogType: 'register',
  mimes: new Map(),
  search: '',
  editingBook: deepCopy(emptyBook),
  files: [],
}

function buildBookParams(book: Model.Book): URLSearchParams {
  const params = new URLSearchParams();
  params.append('ISBN', book.ISBN);
  params.append('Title', book.Title);
  params.append('Author', book.Author);
  params.append('Publisher', book.Publisher);
  params.append('PubDate', book.PubDate);
  params.append('CoverURL', book.CoverURL);
  params.append('Description', book.Description);
  return params;
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
  const getDescription = (): string => {
    const contents = data.onix.CollateralDetail.TextContent;
    if (!contents) return '';
    const description = contents.find((c: Model.OnixTextContent) => c.TextType === '03');
    return description? description.Text : '';
  }
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
    Description: getDescription(),
    Files: [],
  };
  return book;
}

function handleAPIError(commit: Function, error: AxiosError) {
  const alertMessage: Model.AlertMessage = {
    id: 0,
    type: 'error',
    message: String(error),
  };
  commit(VuexMutation.SET_ALERT_MESSAGE, alertMessage);
  console.error(error)
  throw error
}

export default new Vuex.Store({
  state: initialState,
  mutations: { [VuexMutation.SET_ALERT_MESSAGE](state: Model.State, alertMessage: Model.AlertMessage) {
      alertMessage.id = state.alertMessages.length;
      state.alertMessages = state.alertMessages.concat(alertMessage);
    },
    [VuexMutation.DELETE_ALERT_MESSAGE](state: Model.State, alertMessage: Model.AlertMessage) {
      state.alertMessages = state.alertMessages.filter((m: Model.AlertMessage) => m.id != alertMessage.id);
    },
    [VuexMutation.OPEN_DIALOG](state: Model.State) {
      state.dialog = true;
    },
    [VuexMutation.CLOSE_DIALOG](state: Model.State) {
      state.dialog = false;
    },
    [VuexMutation.SET_DIALOG_TYPE](state: Model.State, type: Model.DialogType) {
      state.dialogType = type;
    },
    [VuexMutation.SET_EDITING_BOOK](state: Model.State, book) {
      state.editingBook = deepCopy(book);
    },
    [VuexMutation.UNSET_EDITING_BOOK](state: Model.State) {
      state.editingBook = deepCopy(emptyBook);
    },
    [VuexMutation.SET_BOOKS](state: Model.State, books) {
      state.books = books;
    },
    [VuexMutation.ADD_BOOK](state: Model.State, book) {
      state.books = state.books.concat(book);
    },
    [VuexMutation.UPDATE_BOOK](state: Model.State, book: Model.Book) {
      const books  = deepCopy(state.books)
      for(const i in books) {
        if(books[i].ID == book.ID) {
          books[i] = book;
          break;
        }
      }
      state.books = books;
    },
    [VuexMutation.DELETE_BOOK_BY_ID](state: Model.State, bookID: number) {
      const books = deepCopy(state.books);
      state.books = books.filter((b: Model.Book) => b.ID != bookID);
    },
    [VuexMutation.SET_FILES](state: Model.State, files: File[]) {
      state.files = files;
    },
    [VuexMutation.SET_MIMES](state: Model.State, mimes: Map<string, string>) {
      state.mimes = mimes;
    },
    [VuexMutation.SET_SEARCH_QUERY](state: Model.State, query: string) {
      state.search = query;
    },
  },
  actions: {
    [VuexAction.OPEN_DIALOG]({ commit }, type: Model.DialogType) {
      commit(VuexMutation.OPEN_DIALOG, type);
    },
    async [VuexAction.AUTOCOMPLETE_EDITING_BOOK_BY_ISBN]({ commit }) {
      const isbn = this.state.editingBook.ISBN;
      try {
        const response = await axios.get(OPENBD_ENDPOINT + '/get?isbn=' + isbn)
        if (response.status === 200) {
          const completedBook = extractBookFromOpenBDResponse(response);
          const book = deepCopy(this.state.editingBook);
          book.ISBN = completedBook.ISBN;
          book.Title = completedBook.Title;
          book.Author = completedBook.Author;
          book.Publisher = completedBook.Publisher;
          book.PubDate = completedBook.PubDate;
          book.CoverURL = completedBook.CoverURL;
          book.Description = completedBook.Description;
          commit(VuexMutation.SET_EDITING_BOOK, book);
        } else {
          throw 'failed to fetch book information'
        }
      } catch (error) {
        handleAPIError(commit, error);
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
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.REGISTER_EDITING_BOOK]({ commit }) {
      try {
        // validate files
        const files = this.state.files;
        const types: string[] = [];
        for (const file of files) {
          if(types.includes(file.type)) throw 'file type conflict';
          types.push(file.type);
        }

        // register book entry
        if (!this.state.editingBook.Title) {
          throw 'Title is empty.';
        }
        const params = buildBookParams(this.state.editingBook);
        const bookResponse = await axios.post(API_ENDPOINT + '/book', params);
        if (bookResponse.status === 200) {
          commit(VuexMutation.ADD_BOOK, bookResponse.data);
        } else {
          throw bookResponse.data.err || 'unexpected error';
        }

        // upload files
        console.log(this.state.files)
        const bookID = String(bookResponse.data.ID);
        const formData = new FormData();
        for(const file of this.state.files) {
          formData.append("files", file);
        }
        const filesResponse = await axios.post(API_ENDPOINT + '/book/' + bookID + '/files', formData);
        if (filesResponse.status === 200) {
          for(const result of filesResponse.data) {
            if (result.status !== 'ok') {
              const alertMessage: Model.AlertMessage = {
                id: 0,
                type: "error",
                message: result.content,
              }
              commit(VuexMutation.SET_ALERT_MESSAGE, alertMessage);
            }
          }
        } else {
          throw filesResponse.data.err || 'unexpected error';
        }
      } catch (error) {
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.UPDATE_BOOK]({ commit }) {
      try {
        if (!this.state.editingBook.Title) {
          throw 'Title is empty.';
        }
        const bookID = String(this.state.editingBook.ID);
        const params = buildBookParams(this.state.editingBook);
        const response = await axios.put(API_ENDPOINT + '/book/' + bookID, params);
        if (response.status === 200) {
          commit(VuexMutation.UPDATE_BOOK, response.data);
        } else {
          throw response.data.err || 'unexpected error';
        }

      } catch (error) {
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.DELETE_EDITING_BOOK]({ commit }) {
      try {
        const bookID = String(this.state.editingBook.ID);
        const response = await axios.delete(API_ENDPOINT + '/book/' + bookID);
        if (response.status === 200) {
          commit(VuexMutation.DELETE_BOOK_BY_ID, this.state.editingBook.ID);
        } else {
          throw response.data.err || 'unexpected error';
        }

      } catch (error) {
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.FETCH_MIMES]({ commit }) {
      try {
        const response = await axios.get(API_ENDPOINT + '/mimes');
        if (response.status === 200) {
          commit(VuexMutation.SET_MIMES, response.data);
        } else {
          throw response.data.err || 'unexpected error';
        }
      } catch (error) {
        handleAPIError(commit, error);
      }
    }
  },
  modules: {
  }
})
