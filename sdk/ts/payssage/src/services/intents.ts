import { client } from "./utils/_client";
import type { IntentResponseI } from "./utils/types";

export const intentService = {
  /**
   * Get all Workspace Intents
   * @param name Workspace name
   * @returns The ApiResponse containing all workspace Intents
   */
  getWorkspaceIntents: (
    name: string
  ) =>
    client.get<IntentResponseI>(`/workspaces/${name}/intents`),
};
