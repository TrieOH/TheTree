import { Trash2 } from "lucide-react";
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
    <div className="border-t border-border p-[1.25rem_0.75rem] text-foreground">
      <div className="relative flex flex-col items-center gap-2 flex-1 text-center @[640px]:flex-row @[640px]:items-center @[640px]:gap-4 @[640px]:text-left">
        <DeviceIcon size={40} />
        <div className="flex flex-col min-w-0">
          <h3 className="text-base font-semibold m-0">{deviceI.browser} - {deviceI.os}</h3>
          <span className="inline-flex flex-col text-sm font-extralight @[640px]:flex-row @[640px]:gap-2 @[640px]:items-center">
            {is_current && <strong className="text-primary font-normal">• Sessão Atual •</strong>}
            <span>{ `${user_ip} - ${timeAgo(issued_at)}` }</span>
          </span>
        </div>
        {!is_current && (
          <Trash2 
            size={20} 
            className="absolute top-0 right-0 cursor-pointer text-destructive opacity-60 transition-[opacity,transform] duration-200 hover:opacity-100 hover:scale-[1.05] active:scale-[0.95]"
            onClick={(e) => onClick(e, session_id)}
          />
        )}
      </div>
    </div>
  )
}