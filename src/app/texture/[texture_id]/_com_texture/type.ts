export enum ThreadType{
  Media,
  Text
}
export type ThreadData = {
  title: string,
  id: string,
  type:ThreadType
};