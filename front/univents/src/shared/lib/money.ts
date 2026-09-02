export function formatMoney(cents: number, currency = "BRL") {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency,
  }).format(cents / 100);
}

export function formatPrice(cents: number) {
  return cents === 0 ? "Gratuito" : formatMoney(cents);
}
