import type { HeaderConfigI, HeaderVariant } from "./header.types";

export const headerRegistry: Record<HeaderVariant, HeaderConfigI> = {
  landing: {
    variant: 'landing',
    title: 'TrieAuth',
    titlePosition: 'left',
    centerActions: [
      { id: "landing_link_000", type: 'link', label: 'Features', to: '/' },
      { id: "landing_link_001", type: 'link', label: 'Pricing', to: '/' },
      { id: "landing_link_002", type: 'link', label: 'Docs', to: '/' },
    ],
    rightActions: [
      { id: "landing_auth_button_000", type: 'authButton', visibleOn: 'fixed' },
    ],
  },

  projects: {
    variant: 'projects',
    title: 'Projects',
    titlePosition: 'none',
    disableMobileMenu: true,
    leftActions: [{ id: "projects_back_000", type: 'back', visibleOn: "fixed", to: "/" }],
    rightActions: [{ id: "projects_create_project_000", type: 'createProject', visibleOn: "fixed" }],
    showUserMenu: true,
  },
  "projects/config": {
    variant: 'projects/config',
    title: 'Project Config',
    titlePosition: 'none',
    disableMobileMenu: true,
    leftActions: [{ id: "projects_config_back_000", type: 'back', visibleOn: "fixed", to: "/projects" }],
    showUserMenu: true,
  },
  "schemas/editor": {
    variant: "schemas/editor",
    title: 'Schema Editor',
    titlePosition: 'none',
    disableMobileMenu: true,
    leftActions: [{ id: "schemas_editor_back_000", type: 'back', visibleOn: "fixed", to: "/projects/config" }],
    showUserMenu: true,
  },
  auth: { variant: 'auth' },
  none: { variant: 'none' }
} satisfies Record<HeaderVariant, HeaderConfigI>

headerRegistry.auth = {...headerRegistry.landing, variant: 'auth' }