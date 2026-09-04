import "@/features/upload-queue/associations";
import { UploadQueueProvider } from "./upload-queue-provider";

export default function UploadQueueRuntime() {
  return <UploadQueueProvider>{null}</UploadQueueProvider>;
}
