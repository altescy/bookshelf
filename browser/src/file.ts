export function getBaseMime(mime: string): string {
  return mime.split(";")[0];
}

export const ExtToMime: Map<string, string> = new Map([
  [".azw3", "application/x-mobi8-ebook"],
  [".epub", "application/epub+zip"],
  [".fb2", "application/fb2+zip"],
  [".mobi", "application/x-mobipocket-ebook"],
  [".pdf", "application/pdf"],
  [".txt", "text/plain"],
]);

export const MimeToAlias: Map<string, string> = new Map([
  ["application/x-mobi8-ebook", "azw3"],
  ["application/epub+zip", "epub"],
  ["application/fb2+zip", "fb2"],
  ["application/x-mobipocket-ebook", "mobi"],
  ["application/pdf", "pdf"],
  ["text/plain", "txt"],
]);

export const MimeToColor: Map<string, string> = new Map([
  ["application/x-mobi8-ebook", "orange"],
  ["application/epub+zip", "lime"],
  ["application/fb2+zip", "light-blue"],
  ["application/x-mobipocket-ebook", "amber"],
  ["application/pdf", "cyan"],
  ["text/plain", "blue-grey"],
]);
