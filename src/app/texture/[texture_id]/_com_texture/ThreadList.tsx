"use client";

import { useEffect, useState } from "react";
import Thread from "./Thread";
import { ThreadType, type ThreadData } from "./type";
export default function ThreadList() {
  const [threadList, setThreadList] = useState<ThreadData[]>([]);
  useEffect(() => {
    //fetch thread list here and refresh threadList
    const sampleThreadData: ThreadData[] = [
      {
        title: "general",
        id: "123456",
        type: ThreadType.Text,
      },
      {
        title: "ボンサイダー super sugoi threads v2",
        id: "1234567",
        type: ThreadType.Text,
      },
    ];
    setThreadList(sampleThreadData);
  }, []);
  return (
    <div className="p-2.5 w-50">
      {threadList?.map((thread: ThreadData) => {
        return <Thread title={thread.title} id={thread.id} type={thread.type} key={thread.id}/>;
      })}
    </div>
  );
}
