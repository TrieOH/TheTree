import type { FormFieldI } from '@/shared/model/field'
import type { ProductCreateI } from '.'

export const getProductFields = (): FormFieldI<ProductCreateI>[] => [
  {
    name: 'name',
    label: 'Nome do Produto',
    type: 'text',
    placeholder: 'O nome de exibição do produto.',
    required: true,
  },
  {
    name: 'description',
    label: 'Descrição',
    type: 'textarea',
    placeholder: 'Uma breve descrição do produto.',
  },
  {
    name: 'type',
    label: 'Tipo',
    type: 'select',
    placeholder: 'O tipo de produto.',
    options: [
      { value: 'merchandise', label: 'Mercadoria' },
      { value: 'ticket', label: 'Ingresso' },
      { value: 'token', label: 'Token' },
      { value: 'bundle', label: 'Pacote' },
    ],
    required: true,
  },
  {
    name: 'price_cents',
    label: 'Preço (em centavos)',
    type: 'number',
    placeholder: 'O preço do produto em centavos (ex: 10000 para R$100,00).',
    required: true,
  },
  {
    name: 'has_inventory',
    label: 'Gerenciar Estoque',
    type: 'checkbox',
    placeholder: 'Marque para gerenciar o estoque deste produto.',
  },
  {
    name: 'inventory_quantity',
    label: 'Quantidade em Estoque',
    type: 'number',
    placeholder: 'A quantidade disponível em estoque (se o gerenciamento de estoque estiver ativado).',
    rules: {
      isVisible: ({ has_inventory }) => has_inventory,
    },
  },
  {
    name: 'available_from',
    label: 'Disponível a partir de',
    type: 'datetime',
    placeholder: 'Data e hora em que o produto estará disponível para compra.',
  },
  {
    name: 'available_until',
    label: 'Disponível até',
    type: 'datetime',
    placeholder: 'Data e hora em que o produto deixará de estar disponível para compra.',
  },
  {
    name: 'ticket_id',
    label: 'ID do Ingresso (se tipo for Ingresso)',
    type: 'text',
    placeholder: 'O ID do ingresso associado a este produto (apenas se o tipo for "Ingresso").',
    rules: {
      isVisible: ({ type }) => type === "ticket"
    }
  },
]
