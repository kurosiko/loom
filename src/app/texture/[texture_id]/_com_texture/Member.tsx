export default function Member({
  id
}:Readonly<{
  id: string;
}>) {
  return (
    <div className="flex items-center rounded-lg px-2.5 py-1.5 gap-2.5">
      <img
        src={'/camp_icon.png'}
        alt="Icon"
        className="size-10"
      />
      <p>{ id }</p>
    </div>
  )
}