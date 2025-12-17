import { FaTrashAlt } from "react-icons/fa";
import { deviceIconMap } from "../../../utils/icons/device-icon-map";
import type { SessionI } from "../../../types/sessions-types";
import { getDeviceInfo } from "../../../utils/ua/get-device-info";
import type { MouseEvent } from "react";

interface SessionProps extends SessionI {
  is_current: boolean;
  onClick: (e: MouseEvent<SVGElement>, id: string) => void;
}

export default function SessionCard({
  is_current,
  session_id,
  user_agent,
  user_ip,
  onClick
}: SessionProps) {
  const deviceI = getDeviceInfo(user_agent)
  const DeviceIcon = deviceIconMap[deviceI.device];
  return (
    <div className="trieoh-session">
      <div className="trieoh-session__content">
        <DeviceIcon size={40} />
        <div className="trieoh-session__info">
          <h3>{deviceI.browser} - {deviceI.os}</h3>
          <span className="trieoh-session__meta">
            {is_current && <strong>• Sessão Atual •</strong>}
            <span>{ user_ip }</span>
          </span>
        </div>
        {!is_current && <FaTrashAlt size={20} color="red" onClick={(e) => onClick(e, session_id)}/>}
      </div>
    </div>
  )
}