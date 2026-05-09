import { FaTrashAlt } from "react-icons/fa";
import { deviceIconMap } from "../../../utils/icons/device-icon-map";
import type { SessionI } from "../../../types/sessions-types";
import { getDeviceInfo } from "../../../utils/ua/get-device-info";
import type { MouseEvent } from "react";
import { timeAgo } from "../../../utils/date-utils";

interface SessionProps extends SessionI {
  is_current: boolean;
  onClick: (e: MouseEvent<SVGElement>, id: string) => void;
}

export default function SessionCard({
  is_current,
  session_id,
  user_agent,
  issued_at,
  user_ip,
  onClick
}: SessionProps) {
  const deviceI = getDeviceInfo(user_agent)
  const DeviceIcon = deviceIconMap[deviceI.device];
  return (
    <div className="border-t border-[rgba(0,0,0,0.3)] p-[1.25rem_0.75rem] text-trieoh-neutral2">
      <div className="relative flex flex-col items-center gap-2 flex-1 text-center @[640px]:flex-row @[640px]:items-center @[640px]:gap-4 @[640px]:text-left">
        <DeviceIcon size={40} />
        <div className="flex flex-col min-w-0">
          <h3 className="text-trieoh-base font-semibold m-0">{deviceI.browser} - {deviceI.os}</h3>
          <span className="inline-flex flex-col text-trieoh-sm font-extralight @[640px]:flex-row @[640px]:gap-2 @[640px]:items-center">
            {is_current && <strong className="text-trieoh-primary font-normal">• Sessão Atual •</strong>}
            <span>{ `${user_ip} - ${timeAgo(issued_at)}` }</span>
          </span>
        </div>
        {!is_current && (
          <FaTrashAlt 
            size={20} 
            color="red" 
            onClick={(e) => onClick(e, session_id)}
            className="absolute top-0 right-0 cursor-pointer opacity-60 transition-[opacity,transform] duration-200 hover:opacity-100 hover:scale-[1.05] active:scale-[0.95]"
          />
        )}
      </div>
    </div>
  )
}