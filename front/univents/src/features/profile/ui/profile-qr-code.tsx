import QRCode from "qrcode";
import { useEffect, useState } from "react";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";

export interface ProfileQrCodeProps {
  value: string;
  label?: string;
  size?: number;
  className?: string;
}

export function ProfileQrCode({
  value,
  label = "QR Code do perfil",
  size = 184,
  className,
}: ProfileQrCodeProps) {
  const [source, setSource] = useState<string>();

  useEffect(() => {
    let active = true;
    QRCode.toDataURL(value, {
      width: size,
      margin: 1,
      errorCorrectionLevel: "M",
      color: { dark: "#111827", light: "#ffffff" },
    }).then((url) => {
      if (active) setSource(url);
    });
    return () => {
      active = false;
    };
  }, [size, value]);

  return (
    <div className={className}>
      {source ? (
        <img
          src={source}
          alt={label}
          width={size}
          height={size}
          className="rounded-md"
        />
      ) : (
        <Skeleton
          style={{ width: size, height: size }}
          className="rounded-md"
        />
      )}
    </div>
  );
}
