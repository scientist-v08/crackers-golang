-- ============================================================
-- EXPENSES
-- ============================================================

-- name: CreateExpense :one
INSERT INTO expenses (
    reason_for_expense,
    amount
)
VALUES ($1, $2)
RETURNING
    id,
    reason_for_expense,
    amount;


-- name: GetAllExpenses :many
SELECT
    id,
    reason_for_expense,
    amount
FROM expenses
ORDER BY id;


-- name: GetAllExpensesTotalAmt :one
SELECT COALESCE(SUM(amount), 0)::BIGINT AS total_amount
FROM expenses;


-- ============================================================
-- INVENTORY
-- ============================================================

-- name: AddInventoryItem :one
INSERT INTO inventories (
    brand_or_company,
    state,
    item,
    num_of_boxes,
    num_of_cartons,
    price_per_carton,
    sub_total
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING
    id,
    brand_or_company,
    state,
    item,
    num_of_boxes,
    num_of_cartons,
    price_per_carton,
    sub_total;


-- name: FindInventoryByID :one
SELECT
    id,
    brand_or_company,
    state,
    item,
    num_of_boxes,
    num_of_cartons,
    price_per_carton,
    sub_total
FROM inventories
WHERE id = $1;


-- name: FindPaginatedInventory :many
SELECT
    id,
    brand_or_company,
    state,
    item,
    num_of_boxes,
    num_of_cartons,
    price_per_carton,
    sub_total
FROM inventories
WHERE
    ($1 = '' OR state = $1)
    AND
    ($2 = '' OR item ILIKE '%' || $2 || '%')
ORDER BY id DESC
LIMIT $3
OFFSET $4;


-- name: CountInventory :one
SELECT COUNT(*)::BIGINT AS total
FROM inventories
WHERE
    ($1 = '' OR state = $1)
    AND
    ($2 = '' OR item ILIKE '%' || $2 || '%');


-- name: SumInventorySubTotal :one
SELECT COALESCE(SUM(sub_total), 0)::BIGINT AS total
FROM inventories
WHERE
    ($1 = '' OR state = $1)
    AND
    ($2 = '' OR item ILIKE '%' || $2 || '%');


-- name: UpdateInventoryState :exec
UPDATE inventories
SET
    state = $2,
    sub_total = $3
WHERE id = $1;


-- name: UpdateInventoryDiff :exec
UPDATE inventories
SET
    num_of_cartons = $2,
    sub_total = $3
WHERE id = $1;


-- ============================================================
-- ROUTES
-- ============================================================

-- name: FindRoutesByRole :many
SELECT
    id,
    route,
    heading,
    role
FROM routes
WHERE role = $1;


-- ============================================================
-- USERS
-- ============================================================

-- name: FindUserByEmail :one
SELECT
    id,
    created_at,
    updated_at,
    deleted_at,
    email,
    password,
    roles
FROM users
WHERE email = $1
  AND deleted_at IS NULL;


-- name: CreateUser :one
INSERT INTO users (
    email,
    password,
    roles
)
VALUES ($1, $2, $3)
RETURNING
    id,
    created_at,
    updated_at,
    deleted_at,
    email,
    password,
    roles;

-- ============================================================
-- BILLINGS
-- ============================================================

-- name: CreateBilling :one
INSERT INTO billings (
    "user",
    mobile,
    grand_total,
    finalized_amt
)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    "user",
    mobile,
    grand_total,
    finalized_amt;


-- name: CreatePurchase :one
INSERT INTO purchases (
    billing_id,
    mobile,
    mrp_or_net,
    item,
    quantity,
    discount,
    sub_total
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING
    id,
    billing_id,
    mobile,
    mrp_or_net,
    item,
    quantity,
    discount,
    sub_total;

-- name: CreatePurchases :copyfrom
INSERT INTO purchases (
    billing_id,
    mobile,
    mrp_or_net,
    item,
    quantity,
    discount,
    sub_total
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- ============================================================
-- PRICE LIST
-- ============================================================

-- name: GetAllPriceList :many
SELECT
    id,
    price,
    item,
    brand
FROM price_lists
ORDER BY brand, item;