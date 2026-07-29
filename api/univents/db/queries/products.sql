-- name: CreateProduct :one
INSERT INTO products (edition_id, vendor_code, requires_registration)
VALUES (@edition_id, @vendor_code, @requires_registration)
RETURNING *;

-- name: CreateProductVariant :one
INSERT INTO product_variants (edition_id, product_id, vendor_code, name, description, price, stock)
VALUES (@edition_id, @product_id, @vendor_code, @name, @description, @price, @stock)
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetProductByVendorCode :one
SELECT *
FROM products
WHERE edition_id = @edition_id
  AND vendor_code = @vendor_code
  AND deleted_at IS NULL;

-- name: ListProductsByEdition :many
SELECT *
FROM products
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: GetProductVariantByID :one
SELECT *
FROM product_variants
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetProductVariantByVendorCode :one
SELECT *
FROM product_variants
WHERE edition_id = @edition_id
  AND vendor_code = @vendor_code
  AND deleted_at IS NULL;

-- name: ListProductVariantsByProduct :many
SELECT *
FROM product_variants
WHERE product_id = @product_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: PatchProduct :one
UPDATE products
SET
    vendor_code           = @vendor_code,
    requires_registration = @requires_registration,
    updated_at            = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: PatchProductVariant :one
UPDATE product_variants
SET
    vendor_code  = @vendor_code,
    name         = @name,
    description  = @description,
    price        = @price,
    stock        = @stock,
    gallery_urls = @gallery_urls,
    updated_at   = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProduct :exec
WITH deleted_variants AS (
    UPDATE product_variants
    SET deleted_at = now(), updated_at = now()
    WHERE product_id = @id AND deleted_at IS NULL
)
UPDATE products
SET deleted_at = now(), updated_at = now()
WHERE products.id = @id AND products.deleted_at IS NULL;

-- name: DeleteProductVariant :exec
UPDATE product_variants
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL;
