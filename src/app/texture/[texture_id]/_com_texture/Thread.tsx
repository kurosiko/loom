import {
  AudioLinesIcon,
  LibraryIcon,
  Sidebar,
  Text,
  TextIcon,
} from "lucide-react";
import { ThreadType } from "./type";
import { useId } from "react";

export default function Thread({
  title,
  id,
  type,
}: {
  title: Readonly<string>;
  id: Readonly<string>;
  type: Readonly<ThreadType>;
}) {
  const mask_id = useId();
  return (
    <div className="flex outline hover:bg-neutral-600 gap-2.5 w-full h-7.5 *:my-auto">
      {type === ThreadType.Text ? (
        <LibraryIcon className="shrink-0" absoluteStrokeWidth enableBackground="currentColor"/>
      ) : (
        <AudioLinesIcon className="shrink-0" absoluteStrokeWidth />
      )}
      <p
        className={`flex-1 overflow-x-clip text-nowrap text-ellipsis overflow-hidden`}
        id=""
      >
        {title}
      </p>
    </div>
  );
}
