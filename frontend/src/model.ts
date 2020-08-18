export type AlertMessageType = 'success' | 'warning' | 'error';

export interface AlertMessage {
  id: number;
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
  Files: BookFile[];
}

export interface BookFile {
  ID: number;
  BookID: number;
  MimeType: string;
  Path: string;
}

export type DialogType = 'register' | 'edit';

export interface State {
  alertMessages: AlertMessage[];
  books: Book[];
  dialog: boolean;
  dialogType: DialogType;
  mimes: Map<string, string>;
  search: string;
  editingBook: Book;
  files: File[];
  overlay: boolean;
}

export interface OnixTextContent {
  ContentAudience: string;
  Text: string;
  TextType: string;
}
