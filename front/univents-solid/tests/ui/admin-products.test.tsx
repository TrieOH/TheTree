import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { ProductI, VariantI } from "@/features/products/model";
import { AdminCreateProductCard } from "@/features/products/ui/AdminCreateProductCard";
import { AdminCreateVariantCard } from "@/features/products/ui/AdminCreateVariantCard";
import { AdminProductCard } from "@/features/products/ui/AdminProductCard";
import { AdminVariantCard } from "@/features/products/ui/AdminVariantCard";

afterEach(cleanup);

const mockProduct: ProductI = {
  id: "prod-1",
  edition_id: "ed-1",
  vendor_code: "CAMISETA-DEV",
  requires_registration: true,
  created_at: "2026-03-01T12:00:00Z",
  updated_at: null,
  deleted_at: null,
};

const mockVariant: VariantI = {
  id: "var-1",
  edition_id: "ed-1",
  product_id: "prod-1",
  vendor_code: "CAMISETA-DEV-G",
  name: "Camiseta Dev Tamanho G",
  description: "100% Algodão",
  price: 8900,
  stock: 25,
  gallery_urls: [],
  created_at: "2026-03-01T12:00:00Z",
  updated_at: null,
  deleted_at: null,
};

describe("AdminCreateProductCard", () => {
  it("renders with title, description and triggers onCreate callback", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateProductCard onCreate={onCreate} animate={false} />);

    const btn = screen.getByRole("button", { name: /Novo produto/i });
    expect(btn).toBeInTheDocument();
    expect(
      screen.getByText("Cadastre produtos físicos ou digitais com variações e controle de estoque."),
    ).toBeInTheDocument();

    fireEvent.click(btn);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminProductCard", () => {
  it("renders product info, registration badge and actions", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    const onManageVariants = vi.fn();

    render(() => (
      <AdminProductCard
        product={mockProduct}
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
        onManageVariants={onManageVariants}
      />
    ));

    expect(screen.getByText("CAMISETA-DEV")).toBeInTheDocument();
    expect(screen.getByText("Exige cadastro")).toBeInTheDocument();

    const editBtn = screen.getByLabelText("Editar CAMISETA-DEV");
    fireEvent.click(editBtn);
    expect(onEdit).toHaveBeenCalledWith(mockProduct);

    const deleteBtn = screen.getByLabelText("Excluir CAMISETA-DEV");
    fireEvent.click(deleteBtn);
    expect(onDelete).toHaveBeenCalledWith(mockProduct);

    const variantsBtn = screen.getByRole("button", { name: /Gerenciar variações/i });
    fireEvent.click(variantsBtn);
    expect(onManageVariants).toHaveBeenCalledWith(mockProduct);
  });
});

describe("AdminCreateVariantCard", () => {
  it("renders with title, description and triggers onCreate callback", () => {
    const onCreate = vi.fn();
    render(() => <AdminCreateVariantCard onCreate={onCreate} animate={false} />);

    const btn = screen.getByRole("button", { name: /Nova variação/i });
    expect(btn).toBeInTheDocument();
    expect(
      screen.getByText("Adicione novos tamanhos, cores ou modelos com preços e estoques independentes."),
    ).toBeInTheDocument();

    fireEvent.click(btn);
    expect(onCreate).toHaveBeenCalledTimes(1);
  });
});

describe("AdminVariantCard", () => {
  it("renders variant name, code, price, stock and triggers actions", () => {
    const onEdit = vi.fn();
    const onDelete = vi.fn();

    render(() => (
      <AdminVariantCard
        variant={mockVariant}
        animate={false}
        onEdit={onEdit}
        onDelete={onDelete}
      />
    ));

    expect(screen.getByText("Camiseta Dev Tamanho G")).toBeInTheDocument();
    expect(screen.getByText("CAMISETA-DEV-G")).toBeInTheDocument();
    expect(screen.getByText("100% Algodão")).toBeInTheDocument();
    expect(screen.getByText(/89,00/)).toBeInTheDocument();
    expect(screen.getByText(/25 un/)).toBeInTheDocument();

    const editBtn = screen.getByLabelText("Editar Camiseta Dev Tamanho G");
    fireEvent.click(editBtn);
    expect(onEdit).toHaveBeenCalledWith(mockVariant);

    const deleteBtn = screen.getByLabelText("Excluir Camiseta Dev Tamanho G");
    fireEvent.click(deleteBtn);
    expect(onDelete).toHaveBeenCalledWith(mockVariant);
  });
});
