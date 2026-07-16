import { Extension } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    certificateLineHeight: {
      setCertificateLineHeight: (lineHeight: number) => ReturnType
    }
  }
}

export const CertificateLineHeight = Extension.create({
  name: 'certificateLineHeight',

  addGlobalAttributes() {
    return [
      {
        types: ['paragraph'],
        attributes: {
          lineHeight: {
            default: 1.25,
            parseHTML: (element) =>
              Number.parseFloat(element.style.lineHeight) || 1.25,
            renderHTML: ({ lineHeight }) => ({
              style: `line-height: ${String(lineHeight)}`,
            }),
          },
        },
      },
    ]
  },

  addCommands() {
    return {
      setCertificateLineHeight:
        (lineHeight) =>
        ({ commands }) =>
          commands.updateAttributes('paragraph', { lineHeight }),
    }
  },
})
