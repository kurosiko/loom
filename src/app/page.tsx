export default function Home() {
  return (
    <div className="bg-black text-white w-full text-center items-center flex flex-col gap-5 *:p-5 h-screen">
      <h1 className="text-8xl font-bold bg-linear-to-r from-violet-500 to-fuchsia-500 inline-block text-transparent bg-clip-text">
        Create own texture
      </h1>
      <button className="text-4xl bg-linear-to-r from-violet-500 to-fuchsia-500 inline-block text-transparent bg-clip-text border-2 border-violet-500 rounded-lg px-5 py-3 hover:bg-clip-border hover:text-white transition-all ease-linear duration-300">
        Join now
      </button>
    </div>
  );
}
