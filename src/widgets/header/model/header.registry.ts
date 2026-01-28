import { HeaderConfigI } from "./header.types";

export const headerRegistry: Record<string, HeaderConfigI> = {
  landing: {
    variant: 'landing',
    title: 'TrieAuth',
    centerActions: [ // center is the only who visibleOn always or empty transfer to menu
      { type: 'link', label: 'Features', to: '/' },
      { type: 'link', label: 'Pricing', to: '/' },
      { type: 'link', label: 'Docs', to: '/' },
    ],
    rightActions: [
      { type: 'authButton', visibleOn: 'desktop' }, 
      // { type: 'authButton', visibleOn: 'mobile' }
    ],
  },

  auth: {
    variant: 'auth',
    title: 'TrieAuth',
    centerActions: [
      { type: 'link', label: 'Features', to: '/' },
      { type: 'link', label: 'Pricing', to: '/' },
      { type: 'link', label: 'Docs', to: '/' },
    ],
    rightActions: [
      { type: 'authButton', visibleOn: 'desktop' }, 
      { type: 'authButton', visibleOn: 'mobile' }
    ],
  },

  projects: {
    variant: 'projects',
    title: 'Projects',
    leftActions: [{ type: 'back' }],
    // rightActions: [{ type: 'createProject' }],
  },
} satisfies Record<string, HeaderConfigI>
