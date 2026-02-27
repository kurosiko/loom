import MemberList from "./_com_texture/MemberList";
import ThreadList from "./_com_texture/ThreadList";

export default function TextureLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex *:outline-1">
      <ThreadList />
      <div className="flex-1 h-full">
        {children}
      </div>
      <MemberList />
    </div>   
  )
}