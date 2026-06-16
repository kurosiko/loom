export enum ThreadType {
  Media = "media",
  Text = "text",
}

export interface Texture {
  id: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Thread {
  id: string;
  textureId: string;
  title: string;
  type: ThreadType;
  createdAt: string;
  updatedAt: string;
}

export interface Member {
  id: string;
  username: string;
  avatarUrl?: string;
}

export interface Message {
  id: string;
  threadId: string;
  memberId: string;
  content: string;
  createdAt: string;
}
