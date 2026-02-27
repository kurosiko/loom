export default async function ThreadPage({
  params
}: {
  params: Promise<{ texture_id: string; thread_id: string }>
}) {
  const { texture_id, thread_id } = await params
  return (
    <>
      <h1>page of {texture_id}/ { thread_id }</h1>
    </>
  )
}