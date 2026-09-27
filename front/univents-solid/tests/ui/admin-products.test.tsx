import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { ProductI, VariantI } from "@/features/products/model";
import { AdminCreateProductCard } from "@/features/products/ui/AdminCreateProductCard";
import { AdminCreateVariantCard } from "@/features/products/ui/AdminCreateVariantCard";
import { AdminProductCard } from "@/features/products/ui/AdminProductCard";
import { AdminVariantCard } from "@/features/products/ui/AdminVariantCard";
import { ManageProductDialog } from "@/features/products/ui/ManageProductDialog";
import { ManageVariantDialog } from "@/features/products/ui/ManageVariantDialog";

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

describe("ManageProductDialog", () => {
  it("guides through product and initial variant creation with validation and focus retention", async () => {
    const onOpenChange = vi.fn();
    const onCreate = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageProductDialog
        open={true}
        onOpenChange={onOpenChange}
        product={null}
        onCreate={onCreate}
      />
    ));

    expect(screen.getByText("Novo produto")).toBeInTheDocument();

    // Step 1: Empty vendor_code validation
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(
        screen.getByText("O código do produto deve ter pelo menos 2 caracteres."),
      ).toBeInTheDocument();
    });

    const vendorCodeInput = screen.getByLabelText(/Código do produto/i) as HTMLInputElement;
    vendorCodeInput.focus();
    expect(document.activeElement).toBe(vendorCodeInput);

    fireEvent.input(vendorCodeInput, { target: { value: "CANECA-OFICIAL" } });
    expect(document.activeElement).toBe(vendorCodeInput);

    // Toggle requires registration
    const regCheckbox = screen.getByLabelText(/Exigir cadastro do comprador/i);
    fireEvent.click(regCheckbox);

    // Advance to Step 2: Primeira Variação
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/Código da variação/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Nome da variação/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Preço/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Estoque inicial/i)).toBeInTheDocument();
    });

    // Fill Step 2 fields
    const varCodeInput = screen.getByLabelText(/Código da variação/i);
    fireEvent.input(varCodeInput, { target: { value: "CAN-BRANCA-350" } });

    const varNameInput = screen.getByLabelText(/Nome da variação/i);
    fireEvent.input(varNameInput, { target: { value: "Caneca Branca 350ml" } });

    const priceInput = screen.getByLabelText(/Preço/i);
    fireEvent.input(priceInput, { target: { value: "4500" } }); // R$ 45,00

    const stockInput = screen.getByLabelText(/Estoque inicial/i);
    fireEvent.input(stockInput, { target: { value: "50" } });

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados do produto e variação")).toBeInTheDocument();
      expect(screen.getByText("Pronto para cadastrar")).toBeInTheDocument();
      expect(screen.getAllByText("CANECA-OFICIAL")[0]).toBeInTheDocument();
      expect(screen.getAllByText("Caneca Branca 350ml")[0]).toBeInTheDocument();
      expect(screen.getAllByText("R$ 45,00")[0]).toBeInTheDocument();
      expect(screen.getByText("50 unidades")).toBeInTheDocument();
    });

    // Submit form
    const submitBtn = screen.getByRole("button", { name: /Criar produto/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onCreate).toHaveBeenCalledWith({
        vendor_code: "CANECA-OFICIAL",
        requires_registration: true,
        variant_vendor_code: "CAN-BRANCA-350",
        name: "Caneca Branca 350ml",
        description: null,
        price: 4500,
        stock: 50,
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("handles product editing mode with SKU and registration fields only", async () => {
    const onOpenChange = vi.fn();
    const onUpdate = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageProductDialog
        open={true}
        onOpenChange={onOpenChange}
        product={mockProduct}
        onUpdate={onUpdate}
      />
    ));

    expect(screen.getByText("Editar produto")).toBeInTheDocument();

    const vendorCodeInput = screen.getByLabelText(/Código do produto/i) as HTMLInputElement;
    expect(vendorCodeInput.value).toBe("CAMISETA-DEV");

    fireEvent.input(vendorCodeInput, { target: { value: "CAMISETA-DEV-2026" } });

    // Advance to Step 2: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Pronto para salvar")).toBeInTheDocument();
      expect(screen.getAllByText("CAMISETA-DEV-2026")[0]).toBeInTheDocument();
    });

    const submitBtn = screen.getByRole("button", { name: /Salvar alterações/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onUpdate).toHaveBeenCalledWith({
        vendor_code: "CAMISETA-DEV-2026",
        requires_registration: true,
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });
});

describe("ManageVariantDialog", () => {
  it("guides through variant creation with price, stock and focus retention", async () => {
    const onOpenChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageVariantDialog
        open={true}
        onOpenChange={onOpenChange}
        variant={null}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Nova variação")).toBeInTheDocument();

    // Step 1: Validation failure on empty fields
    const continueBtn = screen.getByRole("button", { name: /Continuar/i });
    fireEvent.click(continueBtn);

    await waitFor(() => {
      expect(
        screen.getByText("O código deve ter pelo menos 2 caracteres."),
      ).toBeInTheDocument();
    });

    const codeInput = screen.getByLabelText(/Código da variação/i) as HTMLInputElement;
    codeInput.focus();
    expect(document.activeElement).toBe(codeInput);

    fireEvent.input(codeInput, { target: { value: "CAN-PRETA-350" } });
    expect(document.activeElement).toBe(codeInput);

    const nameInput = screen.getByLabelText(/Nome da variação/i);
    fireEvent.input(nameInput, { target: { value: "Caneca Preta 350ml" } });

    const descInput = screen.getByLabelText(/Descrição \(opcional\)/i);
    fireEvent.input(descInput, { target: { value: "Cerâmica fosca de alta qualidade." } });

    // Advance to Step 2: Preço e Estoque
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/Preço unitário/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/Estoque disponível/i)).toBeInTheDocument();
    });

    const priceInput = screen.getByLabelText(/Preço unitário/i);
    fireEvent.input(priceInput, { target: { value: "5990" } }); // R$ 59,90

    // Leave stock empty for unlimited stock

    // Advance to Step 3: Resumo
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Dados da variação")).toBeInTheDocument();
      expect(screen.getByText("Pronto para cadastrar")).toBeInTheDocument();
      expect(screen.getAllByText("Caneca Preta 350ml")[0]).toBeInTheDocument();
      expect(screen.getAllByText("CAN-PRETA-350")[0]).toBeInTheDocument();
      expect(screen.getAllByText("R$ 59,90")[0]).toBeInTheDocument();
      expect(screen.getByText("Ilimitado")).toBeInTheDocument();
    });

    // Submit form
    const submitBtn = screen.getByRole("button", { name: /Criar variação/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        vendor_code: "CAN-PRETA-350",
        name: "Caneca Preta 350ml",
        description: "Cerâmica fosca de alta qualidade.",
        price: 5990,
        stock: null,
        gallery_urls: [],
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });

  it("populates existing variant data on edit and updates values", async () => {
    const onOpenChange = vi.fn();
    const onSubmit = vi.fn().mockResolvedValue(true);

    render(() => (
      <ManageVariantDialog
        open={true}
        onOpenChange={onOpenChange}
        variant={mockVariant}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByText("Editar variação")).toBeInTheDocument();

    const nameInput = screen.getByLabelText(/Nome da variação/i) as HTMLInputElement;
    expect(nameInput.value).toBe("Camiseta Dev Tamanho G");

    fireEvent.input(nameInput, { target: { value: "Camiseta Dev GG" } });

    // Advance to Step 2
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByLabelText(/Preço unitário/i)).toBeInTheDocument();
    });

    // Advance to Step 3
    fireEvent.click(screen.getByRole("button", { name: /Continuar/i }));

    await waitFor(() => {
      expect(screen.getByText("Pronto para salvar alterações")).toBeInTheDocument();
      expect(screen.getAllByText("Camiseta Dev GG")[0]).toBeInTheDocument();
    });

    const submitBtn = screen.getByRole("button", { name: /Salvar alterações/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        vendor_code: "CAMISETA-DEV-G",
        name: "Camiseta Dev GG",
        description: "100% Algodão",
        price: 8900,
        stock: 25,
        gallery_urls: [],
      });
      expect(onOpenChange).toHaveBeenCalledWith(false);
    });
  });
});
