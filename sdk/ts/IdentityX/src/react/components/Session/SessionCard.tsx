import { FaTrashAlt } from "react-icons/fa";
import { deviceIconMap } from "../../../utils/icons/device-icon-map";
import type { SessionI } from "../../../types/sessions-types";
import { getDeviceInfo } from "../../../utils/ua/get-device-info";

interface SessionProps extends SessionI {
  is_current: boolean
}

export default function SessionCard({
  is_current,
  user_agent,
  user_ip
}: SessionProps) {
  const deviceI = getDeviceInfo(user_agent)
  const DeviceIcon = deviceIconMap[deviceI.device];
  return (
    <div className="trieoh-session">
      <div className="trieoh-session__main">
        <div className="trieoh-session__content">
          <DeviceIcon size={40} />
          <div className="trieoh-session__info">
            <h3>{deviceI.browser} - {deviceI.os}</h3>
            <span className="trieoh-session__meta">
              {is_current && <strong>• Sessão Atual •</strong>}
              <span>{ user_ip }</span>
            </span>
          </div>
        </div>
        <FaTrashAlt size={20} color="red" />
      </div>
    </div>
  )
}