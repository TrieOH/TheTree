import type { HeaderConfigI, HeaderVariant } from "./header.types";

export const headerRegistry: Record<HeaderVariant, HeaderConfigI> = {
  landing: {
    variant: 'landing',
    title: 'TrieAuth',
    titlePosition: 'left',
    centerActions: [
      { type: 'link', label: 'Features', to: '/' },
      { type: 'link', label: 'Pricing', to: '/' },
      { type: 'link', label: 'Docs', to: '/' },
    ],
    rightActions: [
      { type: 'authButton', visibleOn: 'fixed' },
    ],
  },

  projects: {
    variant: 'projects',
    title: 'Projects',
    titlePosition: 'none',
    disableMobileMenu: true,
    leftActions: [{ type: 'back', visibleOn: "fixed", to: "/" }],
    rightActions: [{ type: 'createProject', visibleOn: "fixed" }],
  },
  auth: { variant: 'auth' },
  none: { variant: 'none' }
} satisfies Record<HeaderVariant, HeaderConfigI>

headerRegistry.auth = {...headerRegistry.landing, variant: 'auth' }