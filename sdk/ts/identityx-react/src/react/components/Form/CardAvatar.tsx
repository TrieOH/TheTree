import { User } from "lucide-react";


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
    <div className="font-sans flex flex-col items-center">
      <div className="flex justify-center items-center p-2 bg-foreground/10 rounded-full mb-2.5">
        <User className="w-16 h-16 p-2.5 bg-background rounded-full shadow-[0_0.25rem_1rem_rgba(0,0,0,0.25)]" size={70} />
      </div>
      <h3 className="text-center text-xl font-medium m-0">{mainText}</h3>
      <span className="text-center text-sm font-light opacity-60">{subText}</span>
    </div>
  )
}