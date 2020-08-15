export function deepCopy<T>(src: T): T {
  return JSON.parse(JSON.stringify(src));
}
