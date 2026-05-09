import { GoPerson } from "react-icons/go";


interface CardAvatarProps {
  /** The main text */
  mainText: string;
  /** Sub Text that will appear below the main text */
  subText: string;
}

export default function CardAvatar({
  mainText,
  subText,
}: CardAvatarProps) {
  return (
    <div className="font-inter flex flex-col items-center">
      <div className="flex justify-center items-center p-2 bg-[oklch(0.8853_0_0/30%)] rounded-trieoh-full mb-[0.625rem]">
        <GoPerson className="w-[64px] h-[64px] p-[0.625rem] bg-trieoh-neutral1 rounded-trieoh-full shadow-[0_0.25rem_1rem_rgba(0,0,0,0.25)]" size={70} />
      </div>
      <h3 className="text-center text-trieoh-xl font-medium m-0">{mainText}</h3>
      <span className="text-center text-[0.875rem] font-light opacity-60">{subText}</span>
    </div>
  )
}