import UserIcon from "../Icons/UserIcon";

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
    <div className="trieoh trieoh-avacard">
      <div className="trieoh-avacard__container">
        <UserIcon className="trieoh-avacard__content"/>
      </div>
      <h3 className="trieoh-avacard__title">{mainText}</h3>
      <span className="trieoh-avacard__sub-title">{subText}</span>
    </div>
  )
}