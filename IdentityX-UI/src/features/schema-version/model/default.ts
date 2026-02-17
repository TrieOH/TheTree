import type { VersionFieldResult } from "./types";

export const defaultEmailVersionField: VersionFieldResult = {
  key: "email",
  title: "Email",
  type: "email",
  owner: "system",
  mutable: false,
  required: true,
  position: -1,
  default_value: '',
  options: [],
  required_rules: [],
  visibility_rules: [],
  id: 'emailTemplate',
  object_id: 'emailTemplate'
}

export const defaultPasswordVersionField: VersionFieldResult = {
  key: "password",
  title: "Password",
  type: "string",
  owner: "system",
  mutable: false,
  required: true,
  position: -1,
  default_value: '',
  options: [],
  required_rules: [],
  visibility_rules: [],
  id: 'passwordTemplate',
  object_id: 'passwordTemplate'
}