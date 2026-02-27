export default function TextureIcon({
  id
}: Readonly<{
  id: string;
}>) {
  return (
    <a href={ `/texture/${ id }` }>
      <div
        className="size-15"
      >
        <img
          src={'/camp_icon.png'}
          alt="Camp Icon"
          className="size-full"
        />
      </div>
    </a>
  )
}