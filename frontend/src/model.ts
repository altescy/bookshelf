export type AlertMessageType = 'success' | 'warning' | 'error';

export interface AlertMessage {
  type: AlertMessageType;
  message: string;
}

export interface Book {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  ISBN: string;
  Title: string;
  Author: string;
  Description: string;
  CoverURL: string;
  Publisher: string;
  PubDate: string;  // format: 2020-01-02
  Files: File[];
}

export interface File {
  ID: number;
  BookID: number;
  MimeType: string;
  Path: string;
}

export type DialogType = 'register' | 'edit';

export interface State {
  alertMessage: AlertMessage;
  books: Book[];
  dialog: boolean;
  dialogType: DialogType;
  search: string;
  editingBook: Book;
}

export interface OnixTextContent {
  ContentAudience: string;
  Text: string;
  TextType: string;
}
