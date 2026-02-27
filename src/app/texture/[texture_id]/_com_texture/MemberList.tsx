"use client";
import { useEffect, useState } from "react";
import Member from "./Member";

export default function MemberList() {
  const [member, setMember] = useState<string[]>([]);
  useEffect(() => {
    setMember(["user1", "user2", "user3"]);
  }, []);
  return (
    <div className="p-2.5 w-50">
      {member?.map((id) => (
        <Member key={id} id={id} />
      ))}
      sidebar
    </div>
  );
}
