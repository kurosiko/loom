import TextureIcon from "./_com_thread/TextureIcon";

export default function TextureLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  
  return (
    <div className="*:py-2.5 min-h-full">
      
      <div className="divide-x-2 divide-gray-600 flex *:px-2.5">
        <TextureIcon id="example_id_1"/>
        <TextureIcon id="example_id_2"/>
        <TextureIcon id="example_id_3"/>
      </div>
      {children}
    </div>
  )
    
}