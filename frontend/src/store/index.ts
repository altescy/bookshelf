import Vue from 'vue';
import Vuex from 'vuex';
import axios, {AxiosResponse} from 'axios';
import * as Model from '@/model';
import * as VuexMutation from '@/vuex/mutation_types';
import * as VuexAction from '@/vuex/action_types';
import {MimeToAlias} from '@/file';
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

function validateBook(book: Model.Book) {
  if (!book.Title) {
    throw 'Title is empty.';
  }
}

function validateFiles(files: File[]) {
  const types: string[] = [];
  for (const file of files) {
    if(types.includes(file.type)) throw 'file type conflict';
    types.push(file.type);
  }
}

async function uploadFiles(commit: Function, bookID: number, files: File[]) {
  const formData = new FormData();
  for(const file of files) {
    formData.append("files", file);
  }
  const filesResponse = await axios.post(API_ENDPOINT + '/book/' + bookID + '/files', formData);
  if (filesResponse.status === 200) {
    for(const result of filesResponse.data) {
      if (result.status === 'ok'){
        commit(VuexMutation.UPDATE_FILE, result.content)
      } else {
        const msg: Model.AlertMessage = {
          id: 0,
          type: "error",
          message: result.file + " : " + result.content,
        }
        commit(VuexMutation.SET_ALERT_MESSAGE, msg);
      }
    }
  } else {
    throw filesResponse.data.err || 'unexpected error';
  }
}

function handleAPIError(commit: Function, error: Error) {
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
      state.books = [book].concat(state.books);
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
    [VuexMutation.DELEFTE_FILE_BY_ID](state: Model.State, fileID: number) {
      // delete file from book list
      state.books = deepCopy(state.books).map((b: Model.Book) => {
        b.Files = b.Files.filter((f: Model.BookFile) => f.ID !== fileID);
        return b;
      });
      // delete file from editing book
      const editingBook = deepCopy(state.editingBook);
      editingBook.Files = editingBook.Files.filter((f: Model.BookFile) => f.ID !== fileID);
      state.editingBook = editingBook;
    },
    [VuexMutation.UPDATE_FILE](state: Model.State, file: Model.BookFile) {
      // update file in book list
      state.books = deepCopy(state.books).map((b: Model.Book): Model.Book => {
        if (b.ID !== file.BookID) return b;
        if (!b.Files) b.Files = []
        b.Files = b.Files.filter((f: Model.BookFile): boolean => f.MimeType !== file.MimeType)
        b.Files.push(file)
        return b;
      });
      // update file in editing book
      if (state.editingBook.ID === file.BookID) {
        const book = state.editingBook;
        book.Files = book.Files.filter((f: Model.BookFile): boolean => f.MimeType !== file.MimeType)
        book.Files.push(file);
        state.editingBook = book;
      }
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
        const book = this.state.editingBook;
        validateBook(book);

        const files = this.state.files;
        validateFiles(files);

        // register book entry
        const params = buildBookParams(book);
        const bookResponse = await axios.post(API_ENDPOINT + '/book', params);
        if (bookResponse.status === 200) {
          commit(VuexMutation.ADD_BOOK, bookResponse.data);
        } else {
          throw bookResponse.data.err || 'unexpected error';
        }

        // upload files
        const bookID = bookResponse.data.ID;
        await uploadFiles(commit, bookID, files);
      } catch (error) {
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.UPDATE_BOOK]({ commit }) {
      try {
        const book = this.state.editingBook;
        validateBook(book);

        const files = this.state.files;
        validateFiles(files);

        await uploadFiles(commit, book.ID, files);

        const params = buildBookParams(this.state.editingBook);
        const response = await axios.put(API_ENDPOINT + '/book/' + book.ID, params);
        if (response.status === 200) {
          commit(VuexMutation.UPDATE_BOOK, response.data);
          commit(VuexMutation.SET_EDITING_BOOK, response.data);
          const msg: Model.AlertMessage = {
            id: 0,
            type: 'success',
            message: 'successfully updated',
          };
          commit(VuexMutation.SET_ALERT_MESSAGE, msg);
        } else {
          throw response.data.err || 'unexpected error';
        }
      } catch (error) {
        handleAPIError(commit, error);
      }
    },
    async [VuexAction.DELETE_EDITING_BOOK]({ commit }) {
      try {
        const book = this.state.editingBook;
        const response = await axios.delete(API_ENDPOINT + '/book/' + book.ID);
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
    },
    async [VuexAction.DELETE_FILE]({ commit }, file: Model.BookFile) {
      try {
        const response = await axios.delete(API_ENDPOINT + '/book/' + file.BookID + '/file/' + MimeToAlias.get(file.MimeType));
        if (response.status === 200) {
          commit(VuexMutation.DELEFTE_FILE_BY_ID, file.ID);
          const msg: Model.AlertMessage = {
            id: 0,
            type: 'success',
            message: 'successfully deleted : ' + MimeToAlias.get(file.MimeType),
          };
          commit(VuexMutation.SET_ALERT_MESSAGE, msg);
        } else {
          throw response.data.err || 'unexpected error';
        }
      } catch (error) {
        handleAPIError(commit, error);
      }
    },
  },
  modules: {
  }
})
