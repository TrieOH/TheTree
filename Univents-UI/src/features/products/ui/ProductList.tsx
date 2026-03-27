import {
  Search,
  ChevronLeft,
  ChevronRight,
  PackageX,
  Ticket,
  Coins,
  Package,
  Layers,
  RotateCcw
} from "lucide-react";
import { Fragment, useEffect, useMemo, useState } from "react";
import { ProductCard } from "./ProductCard";
import { ProductDetails } from "./ProductDetails";
import type { ProductI, ProductType } from "../model";
import { Input } from "@/shared/ui/shadcn/input";
import { Button } from "@/shared/ui/shadcn/button";
import WaveSpinnerLoading from "@/shared/ui/loader/WaveSpinnerLoading";

interface ProductListProps {
  products: ProductI[];
  isLoading?: boolean;
  itemsPerPage?: number;
  inventory: Record<string, number>;
}

const typeConfig: Record<ProductType | "all", { label: string; icon: React.ElementType; color: string }> = {
  all: { label: "Todos", icon: Package, color: "bg-slate-500" },
  merchandise: { label: "Produtos", icon: Package, color: "bg-blue-500" },
  ticket: { label: "Ingressos", icon: Ticket, color: "bg-green-500" },
  token: { label: "Tokens", icon: Coins, color: "bg-amber-500" },
  bundle: { label: "Combos", icon: Layers, color: "bg-purple-500" },
};

export function ProductList({
  products,
  isLoading = false,
  itemsPerPage = 12,
  inventory
}: ProductListProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState<ProductType | "all">("all");
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedProduct, setSelectedProduct] = useState<ProductI | null>(null);

  const handleProductSelect = (product: ProductI) => {
    setSelectedProduct(product);
  };

  const handleDetailsClose = () => {
    setSelectedProduct(null);
  }

  const filteredProducts = useMemo(() => {
    return products.filter(product => {
      const matchesSearch = !searchQuery ||
        product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        product.description?.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesType = selectedType === "all" || product.type === selectedType;

      return matchesSearch && matchesType;
    });
  }, [products, searchQuery, selectedType]);

  const totalPages = Math.ceil(filteredProducts.length / itemsPerPage) || 1;
  const paginatedProducts = useMemo(() => {
    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    return filteredProducts.slice(start, end);
  }, [filteredProducts, currentPage, itemsPerPage]);

  useEffect(() => {
    setCurrentPage(1);
  }, [searchQuery, selectedType]);

  const handleTypeSelect = (type: ProductType | "all") => {
    setSelectedType(type);
  };

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  };

  const clearFilters = () => {
    setSearchQuery("");
    setSelectedType("all");
    setCurrentPage(1);
  };

  const hasActiveFilters = searchQuery || selectedType !== "all";

  return (
    <>
      <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Header */}
        <div className="flex flex-col gap-3 mb-4">
          <div className="flex items-center gap-2">
            <form onSubmit={(e) => { e.preventDefault(); }} className="relative flex-1">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Buscar por nome..."
                value={searchQuery}
                onChange={(e) => { setSearchQuery(e.target.value); }}
                className="pl-9 h-10"
              />
            </form>

            <Button
              variant="outline"
              onClick={clearFilters}
              disabled={!hasActiveFilters}
              className="h-10 px-3"
              title="Limpar filtros"
            >
              <RotateCcw className="h-4 w-4" />
              <span className="hidden sm:inline">Limpar</span>
            </Button>
          </div>

          {/* Type Filter Pills */}
          <div className="-mx-4 px-4 sm:mx-0 sm:px-0 overflow-x-auto scrollbar-hide">
            <div className="flex items-center gap-2 pb-1 sm:flex-wrap sm:pb-0 w-max sm:w-auto min-w-full sm:min-w-0">
              {(Object.keys(typeConfig) as Array<keyof typeof typeConfig>).map((type) => {
                const config = typeConfig[type];
                const Icon = config.icon;
                const isSelected = selectedType === type;
                const count = type === "all" ? products.length : products.filter(p => p.type === type).length;

                return (
                  <button
                    key={type}
                    onClick={() => { handleTypeSelect(type); }}
                    disabled={count === 0}
                    className={`
                    shrink-0 inline-flex items-center gap-1.5 px-3 py-2 rounded-full text-sm font-medium
                    transition-all duration-200 border
                    ${isSelected
                        ? `${config.color} text-white border-transparent`
                        : "bg-background border-border hover:border-primary/50 hover:bg-muted text-foreground"
                      }
                    ${count === 0 ? "opacity-40 cursor-not-allowed" : ""}
                  `}
                  >
                    <Icon className="h-3.5 w-3.5" />
                    <span>{config.label}</span>
                    <span className={`text-[10px] opacity-70 ${isSelected ? 'text-white/70' : 'text-muted-foreground'}`}>
                      {count}
                    </span>
                  </button>
                );
              })}
            </div>
          </div>
        </div>

        {/* Loading State */}
        {isLoading ? (
          <div className="mt-12">
            <WaveSpinnerLoading text="Carregando produtos..." />
          </div>
        ) : (
          <div className="space-y-4">
            {filteredProducts.length > 0 ? (
              <>
                <div className="flex items-center justify-between text-sm text-muted-foreground">
                  <span>{filteredProducts.length} resultado{filteredProducts.length !== 1 ? 's' : ''}</span>
                </div>

                <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
                  {paginatedProducts.map((product) => (
                    <ProductCard
                      key={product.id}
                      product={product}
                      inventoryRemaining={inventory[product.id] ?? product.inventory_remaining}
                      onProductSelect={() => { handleProductSelect(product); }}
                    />
                  ))}
                </div>

                {totalPages > 1 && (
                  <div className="flex items-center justify-center gap-2 pt-6 border-t mt-6">
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => { handlePageChange(currentPage - 1); }}
                      disabled={currentPage <= 1}
                      className="h-8 w-8"
                    >
                      <ChevronLeft className="h-4 w-4" />
                    </Button>

                    <div className="flex items-center gap-1">
                      {Array.from({ length: totalPages }, (_, i) => i + 1)
                        .filter(page =>
                          page === 1 ||
                          page === totalPages ||
                          (page >= currentPage - 1 && page <= currentPage + 1)
                        )
                        .map((page, idx, arr) => (
                          <Fragment key={page}>
                            {idx > 0 && arr[idx - 1] !== page - 1 && (
                              <span className="text-muted-foreground px-1">...</span>
                            )}
                            <Button
                              variant={currentPage === page ? "default" : "outline"}
                              size="icon"
                              onClick={() => { handlePageChange(page); }}
                              className="h-8 w-8 text-xs"
                            >
                              {page}
                            </Button>
                          </Fragment>
                        ))}
                    </div>

                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => { handlePageChange(currentPage + 1); }}
                      disabled={currentPage >= totalPages}
                      className="h-8 w-8"
                    >
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                )}
              </>
            ) : (
              <div className="flex flex-col items-center justify-center text-center py-16 px-4">
                <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
                  <PackageX className="w-8 h-8 text-muted-foreground" />
                </div>
                <h2 className="text-lg font-semibold text-foreground mb-2">
                  Nenhum produto encontrado
                </h2>
                <p className="text-muted-foreground text-sm max-w-xs mb-6">
                  Não conseguimos encontrar nenhum produto que corresponda aos filtros atuais.
                </p>
                <Button
                  onClick={clearFilters}
                  className="gap-2"
                  size="sm"
                >
                  <RotateCcw className="h-4 w-4" />
                  Limpar todos os filtros
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
      <ProductDetails
        isOpen={!!selectedProduct}
        onOpenChange={(isOpen) => {
          if (!isOpen) handleDetailsClose()
        }}
        product={selectedProduct}
        inventoryRemaining={selectedProduct ? inventory[selectedProduct.id] ?? selectedProduct.inventory_remaining : 0}
      />
    </>
  );
}