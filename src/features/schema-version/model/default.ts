import type { VersionFieldList, VersionFieldResult } from "./types";

export const defaultVersionFieldList: VersionFieldList = [
  {
    key: "name",
    title: "Name",
    type: "string",
    owner: "user",
    mutable: true,
    required: true,
    position: 0,
    default_value: null,
    options: [],
    required_rules: [],
    visibility_rules: []
  },
  {
    key: "age",
    title: "Age",
    type: "int",
    owner: "user",
    mutable: true,
    required: false,
    position: 1,
    default_value: null,
    options: [],
    required_rules: [],
    visibility_rules: []
  },
];

export const defaultEmailVersionField: VersionFieldResult = {
  key: "email",
  title: "Email",
  type: "email",
  owner: "system",
  mutable: false,
  required: true,
  position: -1,
  default_value: null,
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
  default_value: null,
  options: [],
  required_rules: [],
  visibility_rules: [],
  id: 'passwordTemplate',
  object_id: 'passwordTemplate'
}