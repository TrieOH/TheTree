// Central registry for upload associations. Feature-specific handlers are
// imported here so queued uploads can resume regardless of the current route.
import "@/features/events/api/upload-association";
import "@/features/editions/api/upload-association";
import "@/features/programs/api/upload-association";
